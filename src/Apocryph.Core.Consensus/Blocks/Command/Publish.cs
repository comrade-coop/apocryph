using System;
using System.Linq;

namespace Apocryph.Core.Consensus.Blocks.Command
{
    public class Publish : IEquatable<Publish>, ICommand
    {
        public (string, byte[]) Message { get; }

        public Publish((string, byte[]) message)
        {
            Message = message;
        }

        public bool Equals(Publish? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return Message.Item1.Equals(other.Message.Item1) && Message.Item2.SequenceEqual(other.Message.Item2);
        }

        public override bool Equals(object? obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((Publish)obj);
        }

        public override int GetHashCode()
        {
            return Message.GetHashCode();
        }
    }
}