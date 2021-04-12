using System;
using System.Linq;

namespace Apocryph.Ipfs
{
    public abstract class ByteSequence<TSelf> : IEquatable<TSelf>
        where TSelf : ByteSequence<TSelf>
    {
        public byte[] Bytes { get; }

        public ByteSequence(byte[] bytes)
        {
            Bytes = bytes;
        }

        public bool Equals(TSelf? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return Bytes.SequenceEqual(other.Bytes);
        }

        public override bool Equals(object? obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((TSelf)obj);
        }

        public override int GetHashCode()
        {
            var hash = new HashCode();
            Array.ForEach(Bytes, hash.Add);
            return hash.ToHashCode();
        }

        // Via https://stackoverflow.com/a/311179
        public override string ToString()
        {
            return BitConverter.ToString(Bytes).Replace("-", "");
        }

        public static bool operator ==(ByteSequence<TSelf>? left, TSelf? right)
        {
            if (left is null) return right is null;
            if (right is null) return false;
            return left.Bytes.SequenceEqual(right.Bytes);
        }
        public static bool operator !=(ByteSequence<TSelf>? left, TSelf? right) => !(left == right);
    }
}