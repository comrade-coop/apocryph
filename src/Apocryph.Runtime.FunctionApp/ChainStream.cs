using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.VirtualNodes;
using Apocryph.Core.Consensus.Communication;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class ChainStream
    {
        private PerperStreamContext? _context;
        private Dictionary<byte[], Chain>? _chains;
        private IPerperStream? _gossips;
        private IPerperStream? _queries;
        private Assigner assigner;
        private IAsyncCollector<object>? _output;

        private Dictionary<int, IEnumerable<IPerperStream>> _streams = new Dictionary<int, IEnumerable<IPerperStream>>();

        public ChainStream()
        {
            assigner = new Assigner(CreateNode);
        }

        [FunctionName(nameof(ChainStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("chains")] Dictionary<byte[], Chain> chains,
            [Perper("gossips")] IPerperStream gossips,
            [Perper("queries")] IPerperStream queries,
            [Perper("slotGossips")] IAsyncEnumerable<SlotClaim> slotGossips,
            [Perper("salts")] IAsyncEnumerable<(byte[], int, byte[])> salts,
            [Perper("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            _context = context;
            _chains = chains;
            _gossips = gossips;
            _queries = queries;
            _output = output;

            foreach (var (chainId, chain) in chains)
            {
                assigner.AddChain(chainId, chain.SlotCount);
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
                    var chains = _chains!;
                    var chainData = _chains![chainId];
                    var queries = _queries!;
                    var gossips = _gossips!;
                    var chain = _context!.GetStream();

                    var filter = _context!.DeclareStream(typeof(FilterStream));

                    var consensus = await _context!.StreamFunctionAsync(typeof(ConsensusStream), new
                    {
                        chain = chain.Subscribe(),
                        filter = filter.Subscribe(),
                        queries = queries.Subscribe(),
                        chainData,
                        node,
                        nodes = assigner.GetNodes(chainId)
                    });

                    var validator = await _context!.StreamFunctionAsync(typeof(ValidatorStream), new
                    {
                        consensus = consensus.Subscribe(),
                        filter = filter.Subscribe(),
                        queries = queries.Subscribe(),
                        chainData,
                        node
                    });

                    var ibc = await _context!.StreamFunctionAsync(typeof(IBCStream), new
                    {
                        chain = chain.Subscribe(),
                        validator = validator.Subscribe(),
                        gossips = gossips.Subscribe(),
                        node,
                        nodes = assigner.GetNodes()
                    });

                    await _context!.StreamFunctionAsync(filter, new
                    {
                        ibc = ibc.Subscribe(),
                        gossips = gossips.Subscribe(),
                        chains,
                        node
                    });

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