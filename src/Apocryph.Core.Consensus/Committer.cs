using System.Collections.Generic;
using System.Linq;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.VirtualNodes;

namespace Apocryph.Core.Consensus
{
    public class Committer
    {
        private readonly Dictionary<(Block, GossipVerb), HashSet<Node>> _gossips;

        public Committer()
        {
            _gossips = new Dictionary<(Block, GossipVerb), HashSet<Node>>();
        }

        public void AddGossip(Gossip<Block> gossip)
        {
            var fact = (gossip.Value, gossip.Verb);
            if (!_gossips.ContainsKey(fact))
            {
                _gossips[fact] = new HashSet<Node>();
            }
            _gossips[fact].Add(gossip.Sender);
        }

        public bool IsGossipConfirmed(Block block, GossipVerb verb, Node?[] nodes)
        {
            var fact = (block, verb);
            if (!_gossips.ContainsKey(fact))
            {
                return false;
            }
            var confirmations = _gossips[(block, verb)].Intersect(nodes).Count();
            var total = nodes.Length;
            return 3 * confirmations > 2 * total;
        }

        public bool IsGossipConfirmed(Gossip<Block> gossip, Node?[] nodes)
        {
            return IsGossipConfirmed(gossip.Value, gossip.Verb, nodes);
        }
    }
}