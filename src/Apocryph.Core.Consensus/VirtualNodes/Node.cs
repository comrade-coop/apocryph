using System;

namespace Apocryph.Core.Consensus.VirtualNodes
{
    public class Node : IEquatable<Node>
    {
        public int Id { get; set; }
        public string Chain { get; set; }

        public bool Equals(Node? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return Id == other.Id;
        }

        public override bool Equals(object? obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((Node)obj);
        }

        public override int GetHashCode()
        {
            return HashCode.Combine(Id);
        }
    }
}