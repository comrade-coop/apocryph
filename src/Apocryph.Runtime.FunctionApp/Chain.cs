using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.VirtualNodes;
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
        private IAsyncCollector<object>? _output;
        private byte[]? _chainId;
        private Node[]? _nodes;

        private Dictionary<PrivateKey, IEnumerable<IAsyncDisposable>> _streams = new Dictionary<PrivateKey, IEnumerable<IAsyncDisposable>>();

        [FunctionName(nameof(Chain))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("chainId")] byte[] chainId,
            [Perper("slotCount")] int slotCount,
            [Perper("gossips")] IAsyncDisposable gossips,
            [Perper("queries")] IAsyncDisposable queries,
            [PerperStream("slotGossips")] IAsyncEnumerable<SlotClaim> slotGossips,
            [PerperStream("salts")] IAsyncEnumerable<(int, byte[])> salts,
            [PerperStream("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _chainId = chainId;
            _context = context;
            _gossips = gossips;
            _queries = queries;
            _output = output;
            _nodes = Enumerable.Range(0, slotCount).Select(slot => new Node(_chainId!, slot)).ToArray();

            var assigner = new Assigner(slotCount, _chainId, AddPrivateKey, RemovePrivateKey);

            await Task.WhenAll(
                ProcessGossip(assigner, slotGossips),
                ProcessSalts(assigner, salts),
                Miner.RunAsync(assigner, cancellationToken));
        }

        private async Task ProcessGossip(Assigner assigner, IAsyncEnumerable<SlotClaim> slotGossips)
        {
            await foreach (var gossip in slotGossips)
            {
                if (gossip.ChainId.SequenceEqual(_chainId))
                {
                    assigner.AddKey(gossip.Key, null);
                    // Forward gossip
                }
            }
        }

        private async Task ProcessSalts(Assigner assigner, IAsyncEnumerable<(int, byte[])> salts)
        {
            await foreach (var (slot, salt) in salts)
            {
                assigner.SetSalt(slot, salt);
            }
        }

        private async void AddPrivateKey(int slot, PrivateKey privateKey)
        {
            var queries = _queries!;
            var gossips = _gossips!;
            var nodes = _nodes!;
            var node = new Node(_chainId!, slot);

            var filter = _context!.DeclareStream(typeof(Filter));
            var consensus = await _context!.StreamFunctionAsync(typeof(Consensus), new { filter, queries, nodes, node });
            var validator = await _context!.StreamFunctionAsync(typeof(Validator), new { consensus, queries, nodes, node });
            var ibc = await _context!.StreamFunctionAsync(typeof(IBC), new { validator, gossips, node });
            await _context!.StreamFunctionAsync(filter, new { ibc, gossips, nodes, node });

            _streams[privateKey] = new[] { filter, consensus, validator, ibc };

            await Task.WhenAll(new[] { filter, consensus, validator, ibc }.Select(
                stream => _output!.AddAsync(stream)));

            await _output!.AddAsync(new SlotClaim { Key = privateKey.PublicKey, ChainId = _chainId! });
        }

        private async void RemovePrivateKey(int slot, PrivateKey privateKey)
        {
            foreach (var stream in _streams[privateKey])
            {
                // TODO: Remove from peering instead?
                await stream.DisposeAsync();
            }

            _streams.Remove(privateKey);
        }
    }
}