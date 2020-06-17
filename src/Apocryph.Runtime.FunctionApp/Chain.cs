using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.VirtualNodes;
using Apocryph.Core.Consensus.Communication;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class Chain
    {
        private PerperStreamContext? _context;
        private IAsyncDisposable? _gossips;
        private IAsyncDisposable? _queries;
        private Assigner assigner;
        private IAsyncCollector<object>? _output;

        private Dictionary<int, IEnumerable<IAsyncDisposable>> _streams = new Dictionary<int, IEnumerable<IAsyncDisposable>>();

        public Chain()
        {
            assigner = new Assigner(CreateNode);
        }

        [FunctionName(nameof(Chain))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("chains")] Dictionary<byte[], int> chains,
            [Perper("gossips")] IAsyncDisposable gossips,
            [Perper("queries")] IAsyncDisposable queries,
            [PerperStream("slotGossips")] IAsyncEnumerable<SlotClaim> slotGossips,
            [PerperStream("salts")] IAsyncEnumerable<(byte[], int, byte[])> salts,
            [PerperStream("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _context = context;
            _gossips = gossips;
            _queries = queries;
            _output = output;

            foreach (var (chainId, slotCount) in chains)
            {
                assigner.AddChain(chainId, slotCount);
            }

            await Task.WhenAll(
                ProcessGossip(slotGossips),
                ProcessSalts(salts),
                Miner.RunAsync(assigner, cancellationToken));
        }

        private async Task ProcessGossip(IAsyncEnumerable<SlotClaim> slotGossips)
        {
            await foreach (var gossip in slotGossips)
            {
                assigner.AddKey(gossip.ChainId, gossip.Key, null);
                // Forward gossip
            }
        }

        private async Task ProcessSalts(IAsyncEnumerable<(byte[], int, byte[])> salts)
        {
            await foreach (var (chainId, slot, salt) in salts)
            {
                assigner.SetSalt(chainId, slot, salt);
            }
        }

        private Node CreateNode(byte[] chainId, int slot, PublicKey publicKey, PrivateKey? privateKey)
        {
            var node = new Node(chainId, slot);
            Task.Run(async () =>
            {
                if (_streams.ContainsKey(slot))
                {
                    foreach (var stream in _streams[slot])
                    {
                        // TODO: Remove from peering instead?
                        await stream.DisposeAsync();
                    }

                    _streams.Remove(slot);
                }

                if (privateKey != null)
                {
                    var queries = _queries!;
                    var gossips = _gossips!;
                    var chain = _context!.GetStream();

                    var filter = _context!.DeclareStream(typeof(Filter));
                    var consensus = await _context!.StreamFunctionAsync(typeof(Consensus), new { chain, filter, queries, node, nodes = assigner.GetNodes(chainId) });
                    var validator = await _context!.StreamFunctionAsync(typeof(ValidatorStream), new { consensus, queries, node });
                    var ibc = await _context!.StreamFunctionAsync(typeof(IBC), new { chain, validator, gossips, node, nodes = assigner.GetNodes() });
                    await _context!.StreamFunctionAsync(filter, new { ibc, gossips, node });

                    _streams[slot] = new[] { filter, consensus, validator, ibc };

                    await Task.WhenAll(new[] { filter, consensus, validator, ibc }.Select(
                        stream => _output!.AddAsync(stream)));
                    await _output!.AddAsync(new SlotClaim { Key = privateKey.Value.PublicKey, ChainId = chainId });
                }

                await _output!.AddAsync(new Message<(byte[], Node?[])>((chainId, assigner.GetNodes(chainId)), MessageType.Valid));
            });

            return node;
        }
    }
}