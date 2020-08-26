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
        public IDictionary<Block, int> BlockCounts { get; }

        public SnowballReport(Node source, IDictionary<Block, int> blockCounts) : base(source)
        {
            BlockCounts = blockCounts;
        }
    }

    public class ConsensusReport : Report
    {
        public int Round { get; }
        public Block LastBlock { get; }

        public ConsensusReport(Node source, int round, Block lastBlock) : base(source)
        {
            Round = round;
            LastBlock = lastBlock;
        }
    }
}