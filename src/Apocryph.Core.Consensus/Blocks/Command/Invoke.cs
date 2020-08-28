using System;
using System.Linq;
using System.Text;

namespace Apocryph.Core.Consensus.Blocks.Command
{
    public class Invoke : IEquatable<Invoke>, ICommand
    {
        public Guid Reference { get; }
        public (string, byte[]) Message { get; }

        public Invoke(Guid reference, (string, byte[]) message)
        {
            Reference = reference;
            Message = message;
        }

        public bool Equals(Invoke? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return Reference.Equals(other.Reference) && Message.Item1.Equals(other.Message.Item1) && Message.Item2.SequenceEqual(other.Message.Item2);
        }

        public override bool Equals(object? obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((Invoke)obj);
        }

        public override int GetHashCode()
        {
            return HashCode.Combine(Reference, Message.Item1);
        }

        public override string ToString()
        {
            return $"Invoke({Reference}, {Message.Item1}, {Encoding.UTF8.GetString(Message.Item2)})";
        }
    }
}