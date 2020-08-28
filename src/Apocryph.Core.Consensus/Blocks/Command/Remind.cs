using System;
using System.Linq;

namespace Apocryph.Core.Consensus.Blocks.Command
{
    public class Remind : IEquatable<Remind>, ICommand
    {
        public DateTime DueDateTime { get; }
        public (string, byte[]) Message { get; }

        public Remind(DateTime dueDateTime, (string, byte[]) message)
        {
            DueDateTime = dueDateTime;
            Message = message;
        }

        public bool Equals(Remind? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return DueDateTime.Equals(other.DueDateTime) && Message.Item1.Equals(other.Message.Item1) && Message.Item2.SequenceEqual(other.Message.Item2);
        }

        public override bool Equals(object? obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((Remind)obj);
        }

        public override int GetHashCode()
        {
            return HashCode.Combine(DueDateTime, Message);
        }
    }
}