using System.Collections.Generic;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class Assigner
    {
        private class NodeData
        {
            public byte[] Salt { get; set; } = new byte[0];
            public PublicKey? Key { get; set; }
        }

        private Node[]? _nodes;
        private byte[]? _chainId;
        private Dictionary<Node, NodeData> _nodeData { get; set; } = new Dictionary<Node, NodeData>();
        private IAsyncCollector<object>? _output;
        private string _chain;

        [FunctionName(nameof(Assigner))]
        public async Task Run<T>([PerperStreamTrigger] PerperStreamContext context,
            [Perper("chain")] string chain,
            [Perper("privateKey")] PrivateKey privateKey,
            [PerperStream("gossips")] IAsyncEnumerable<SlotClaim> gossip,
            [PerperStream("salts")] IAsyncEnumerable<(Node, byte[])> salts,
            [PerperStream("output")] IAsyncCollector<object> output)
        {
            _nodes = new Node[] { new Node { Id = 1 } }; //Slots?
            _chainId = GetPrefix(privateKey);
            _chain = chain;
            _output = output;

            foreach (var node in _nodes)
            {
                _nodeData[node] = new NodeData();
            }

            await Task.WhenAll(
                ProcessSalts(salts),
                ProcessClaims(gossip),
                ProcessGeneratedKeys(privateKey));
        }

        private async Task ProcessSalts(IAsyncEnumerable<(Node, byte[])> salts)
        {
            await foreach (var (node, salt) in salts)
            {
                _nodeData[node].Salt = salt;
            }
        }

        private async Task ProcessGeneratedKeys(PrivateKey privateKey)
        {
            if (AddKey(privateKey.PublicKey))
            {
                await _output!.AddAsync(new SlotClaim { Key = privateKey.PublicKey, ChainId = _chainId! }); // Gossip
                await _output!.AddAsync((new Node {Chain = _chain}, _nodes));
            }
        }

        private async Task ProcessClaims(IAsyncEnumerable<SlotClaim> claims)
        {
            await foreach (var claim in claims)
            {
                if (claim.ChainId == _chainId)
                {
                    if (AddKey(claim.Key))
                    {
                        await _output!.AddAsync(claim); // Forward gossip
                        await _output!.AddAsync((new Node {Chain = _chain}, _nodes));
                    }
                }
            }
        }

        private Node GetNodeForKey(PublicKey key)
        {
            return _nodes![(int)(key.GetPosition() % _nodes!.Length)];
        }

        private bool AddKey(PublicKey key)
        {
            var node = GetNodeForKey(key);
            var nodeData = _nodeData[node];

            if (nodeData.Key is PublicKey slotKey
                && slotKey.GetDifficulty(_chainId!, nodeData.Salt) > key.GetDifficulty(_chainId!, nodeData.Salt))
            {
                return false;
            }

            nodeData.Key = key;

            return true;
        }

        private static byte[] GetPrefix(PrivateKey privateKey)
        {
            return default!; //prefix of the private key
        }
    }
}