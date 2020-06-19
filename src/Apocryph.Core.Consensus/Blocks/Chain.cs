using System;

namespace Apocryph.Core.Consensus.Blocks
{
    public class Chain : IEquatable<Chain>
    {
        public int SlotCount { get; set; }
        public Block GenesisBlock { get; set; }

        public Chain(int slotCount, Block genesisBlock)
        {
            SlotCount = slotCount;
            GenesisBlock = genesisBlock;
        }

        public bool Equals(Chain? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return SlotCount.Equals(other.SlotCount) && GenesisBlock.Equals(other.GenesisBlock);
        }

        public override bool Equals(object? obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((Chain)obj);
        }

        public override int GetHashCode()
        {
            return HashCode.Combine(SlotCount, GenesisBlock);
        }
    }
}