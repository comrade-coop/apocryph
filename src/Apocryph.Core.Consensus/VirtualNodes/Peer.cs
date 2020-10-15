using System;
using System.Linq;

namespace Apocryph.Core.Consensus.VirtualNodes
{
    public struct Peer : IEquatable<Peer>
    {
        public byte[] Value { get; set; }

        public Peer(byte[] value)
        {
            Value = value;
        }

        public bool Equals(Peer other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return Value.SequenceEqual(other.Value);
        }

        public override bool Equals(object? obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((Peer)obj);
        }

        public override int GetHashCode()
        {
            var hash = new HashCode();
            Array.ForEach(Value, hash.Add);
            return hash.ToHashCode();
        }

        // Via https://stackoverflow.com/a/311179
        public override string ToString()
        {
            return BitConverter.ToString(Value).Replace("-", "");
        }
    }
}