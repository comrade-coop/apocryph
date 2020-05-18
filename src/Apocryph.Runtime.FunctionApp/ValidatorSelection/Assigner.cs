// NOTE: File is ignored by .csproj file

using System.Collections.Generic;
using System.Threading.Tasks;
using Apocryph.Runtime.FunctionApp.Consensus.Core;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.ValidatorSelection
{
    public class Assigner
    {
        private class NodeData
        {
            public byte[] Salt { get; set; } = new byte[0];
            public PublicKey? Key { get; set; }
        }

        private Node[]? _nodes;
        private byte[]? _agentId;
        private Dictionary<Node, NodeData> _nodeData { get; set; } = new Dictionary<Node, NodeData>();
        private IAsyncCollector<object>? _output;

        [FunctionName(nameof(Assigner))]
        public async Task Run<T>([PerperStreamTrigger] PerperStreamContext context,
            [Perper("nodes")] Node[] nodes,
            [Perper("agentId")] byte[] agentId,
            [PerperStream("gossips")] IAsyncEnumerable<SlotClaim> gossip,
            [PerperStream("salts")] IAsyncEnumerable<(Node, byte[])> salts,
            [PerperStream("keys")] IAsyncEnumerable<PrivateKey> keys,
            [PerperStream("output")] IAsyncCollector<object> output)
        {
            _nodes = nodes;
            _agentId = agentId;
            _output = output;

            foreach (var node in nodes)
            {
                _nodeData[node] = new NodeData();
            }

            await Task.WhenAll(
                ProcessSalts(salts),
                ProcessClaims(gossip),
                ProcessGeneratedKeys(keys));
        }

        private async Task ProcessSalts(IAsyncEnumerable<(Node, byte[])> salts)
        {
            await foreach(var (node, salt) in salts)
            {
                _nodeData[node].Salt = salt;
            }
        }

        private async Task ProcessGeneratedKeys(IAsyncEnumerable<PrivateKey> keys)
        {
            await foreach(var privateKey in keys)
            {
                if (AddKey(privateKey.PublicKey))
                {
                    await _output!.AddAsync(new SlotClaim { Key = privateKey.PublicKey, AgentId = _agentId! }); // Gossip
                    await _output!.AddAsync((true, GetNodeForKey(privateKey.PublicKey))); // Local Node
                }
            }
        }

        private async Task ProcessClaims(IAsyncEnumerable<SlotClaim> claims)
        {
            await foreach(var claim in claims)
            {
                if (claim.AgentId == _agentId)
                {
                    if (AddKey(claim.Key))
                    {
                        await _output!.AddAsync(claim); // Forward gossip
                        await _output!.AddAsync((false, GetNodeForKey(claim.Key))); // Remote Node
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
                && slotKey.GetDifficulty(_agentId!, nodeData.Salt) > key.GetDifficulty(_agentId!, nodeData.Salt))
            {
                return false;
            }

            nodeData.Key = key;

            return true;
        }
    }
}