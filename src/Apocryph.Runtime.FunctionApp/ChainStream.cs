using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class ChainStream
    {
        private PerperStreamContext? _context;
        private Dictionary<Guid, Chain>? _chains;
        private IPerperStream? _gossips;
        private IPerperStream? _queries;
        private string? _hashRegistryWorker;
        private Peer? _self;
        private Assigner assigner = new Assigner();
        private IAsyncCollector<object>? _output;

        private Dictionary<Node, IEnumerable<IPerperStream>> _streams = new Dictionary<Node, IEnumerable<IPerperStream>>();

        public ChainStream()
        {
            assigner.SlotOccupantChanged += SlotOccupantChanged;
        }

        [FunctionName(nameof(ChainStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("self")] Peer self,
            [Perper("chains")] Dictionary<Guid, Chain> chains,
            [Perper("gossips")] IPerperStream gossips,
            [Perper("queries")] IPerperStream queries,
            [Perper("hashRegistryWorker")] string hashRegistryWorker,
            [Perper("slotGossips")] IAsyncEnumerable<SlotClaim> slotGossips,
            [Perper("salts")] IAsyncEnumerable<(Guid, int, byte[])> salts,
            [Perper("output")] IAsyncCollector<object> output,
            CancellationToken cancellationToken)
        {
            var proofLength = 64; // TODO: Move constant to parameter
            await Task.Delay(1000);
            _context = context;
            _chains = chains;
            _gossips = gossips;
            _queries = queries;
            _hashRegistryWorker = hashRegistryWorker;
            _output = output;
            _self = self;

            foreach (var (chainId, chain) in chains)
            {
                assigner.AddChain(chainId, chain.SlotCount);
                await _output!.AddAsync(new Message<(Guid, Node?[])>((chainId, assigner.GetNodes(chainId)), MessageType.Valid));
            }

            await TaskHelper.WhenAllOrFail(
                ProcessGossip(slotGossips),
                ProcessSalts(salts),
                Miner.RunAsync(self, proofLength, assigner, cancellationToken));
        }

        private async Task ProcessGossip(IAsyncEnumerable<SlotClaim> slotGossips)
        {
            await foreach (var gossip in slotGossips)
            {
                assigner.ProcessClaim(gossip);
                // Handle forwarding
            }
        }

        private async Task ProcessSalts(IAsyncEnumerable<(Guid, int, byte[])> salts)
        {
            await foreach (var (chainId, slot, salt) in salts)
            {
                assigner.SetSalt(chainId, slot, salt);
            }
        }

        private async void SlotOccupantChanged(Node node, Peer peer, byte[] proof)
        {
            try
            {
                var peerIsSelf = _self!.Equals(peer);

                if (_streams.ContainsKey(node) && !peerIsSelf)
                {
                    foreach (var stream in _streams[node])
                    {
                        // TODO: Remove from peering as well
                        await stream.DisposeAsync();
                    }

                    _streams.Remove(node);
                }

                if (peerIsSelf)
                {

                    var chains = _chains!;
                    var chainData = _chains![node.ChainId];
                    var queries = _queries!.Filter("Receiver", node);
                    var gossips = _gossips!;
                    var hashRegistryWorker = _hashRegistryWorker!;
                    var chain = _context!.GetStream();

                    var filter = _context!.DeclareStream($"Filter-{node}", typeof(FilterStream));
                    var validator = _context!.DeclareStream($"Validator-{node}", typeof(ValidatorStream));

                    var consensus = await _context!.StreamFunctionAsync($"Consensus-{node}", typeof(ConsensusStream), new
                    {
                        chain = chain.Subscribe(),
                        validator = validator.Subscribe(),
                        filter = filter.Subscribe(),
                        queries = queries.Subscribe(),
                        hashRegistryWorker,
                        chainData,
                        node,
                        proposerAccount = Guid.NewGuid(),
                        nodes = assigner.GetNodes(node.ChainId).ToList()
                    });

                    await _context!.StreamFunctionAsync(validator, new
                    {
                        consensus = consensus.Subscribe(),
                        filter = filter.Subscribe(),
                        queries = queries.Subscribe(),
                        hashRegistryWorker,
                        chainData,
                        node
                    });

                    var ibc = await _context!.StreamFunctionAsync($"IBC-{node}", typeof(IBCStream), new
                    {
                        chain = chain.Subscribe(),
                        validator = validator.Subscribe(),
                        gossips = gossips.Subscribe(),
                        nodes = assigner.GetNodes(),
                        node
                    });

                    await _context!.StreamFunctionAsync(filter, new
                    {
                        ibc = ibc.Subscribe(),
                        gossips = gossips.Subscribe(),
                        hashRegistryWorker,
                        chains,
                        node
                    });

                    _streams[node] = new[] { filter, consensus, validator, ibc };

                    await Task.WhenAll(new[] { filter, consensus, validator, ibc }.Select(
                        stream => _output!.AddAsync(stream)));
                    await _output!.AddAsync(new SlotClaim(node.ChainId, peer, proof));
                }

                await _output!.AddAsync((node, peer));
            }
            catch (Exception e)
            {
                Console.WriteLine(e);
                throw;
            }
        }
    }
}