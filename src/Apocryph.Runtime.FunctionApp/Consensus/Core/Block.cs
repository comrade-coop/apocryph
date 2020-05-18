using System;
using System.Collections.Generic;

namespace Apocryph.Runtime.FunctionApp.Consensus.Core
{
    public class Block : IEquatable<Block>
    {
        public Node? Proposer { get; set; }
        public byte[]? State { get; }
        public object[] Commands { get; }
        public IDictionary<Guid, (string, string[])> Capabilities { get; }

        public Block(byte[]? state, object[] commands, IDictionary<Guid, (string, string[])> capabilities)
        {
            State = state;
            Commands = commands;
            Capabilities = capabilities;
        }

        public bool Equals(Block? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return Equals(State, other.State) && Commands.Equals(other.Commands) && Capabilities.Equals(other.Capabilities);
        }

        public override bool Equals(object? obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((Block)obj);
        }

        public override int GetHashCode()
        {
            return HashCode.Combine(State, Commands, Capabilities);
        }
    }
}