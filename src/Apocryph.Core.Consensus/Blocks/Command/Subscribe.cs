using System;

namespace Apocryph.Core.Consensus.Blocks.Command
{
    public class Subscribe : IEquatable<Subscribe>, ICommand
    {
        public string Target { get; }

        public Subscribe(string target)
        {
            Target = target;
        }

        public bool Equals(Subscribe? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return Target.Equals(other.Target);
        }

        public override bool Equals(object? obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((Subscribe)obj);
        }

        public override int GetHashCode()
        {
            return Target.GetHashCode();
        }
    }
}