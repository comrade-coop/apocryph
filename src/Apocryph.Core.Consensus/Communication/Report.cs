using System.Collections.Generic;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.VirtualNodes;

namespace Apocryph.Core.Consensus.Communication
{
    public class Report
    {
        public Node Source { get; }

        public Report(Node source)
        {
            Source = source;
        }
    }

    public class SnowballReport : Report
    {
        public Dictionary<Block, int> BlockCounts { get; }

        public SnowballReport(Node source, Dictionary<Block, int> blockCounts) : base(source)
        {
            BlockCounts = blockCounts;
        }
    }

    public class ConsensusReport : Report
    {
        public int Round { get; }
        public Node?[] Proposers { get; }
        public Block LastBlock { get; }

        public ConsensusReport(Node source, Node?[] proposers, int round, Block lastBlock) : base(source)
        {
            Round = round;
            Proposers = proposers;
            LastBlock = lastBlock;
        }
    }
}