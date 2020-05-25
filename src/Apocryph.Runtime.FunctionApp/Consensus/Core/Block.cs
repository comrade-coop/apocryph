using System;
using System.Collections.Generic;

namespace Apocryph.Runtime.FunctionApp.Consensus.Core
{
    public class Block : IEquatable<Block>
    {
        public Node? Proposer { get; set; }
        public object[] InputCommands { get; }
        public object[] Commands { get; }
        public IDictionary<string, byte[]> States { get; }
        public IDictionary<Guid, (string, string[])> Capabilities { get; }

        public Block(IDictionary<string, byte[]> states, object[] inputCommands, object[] commands, IDictionary<Guid, (string, string[])> capabilities)
        {
            States = states;
            InputCommands = inputCommands;
            Commands = commands;
            Capabilities = capabilities;
        }

        public bool Equals(Block? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return States.Equals(other.States) && InputCommands.Equals(other.InputCommands) && Commands.Equals(other.Commands) && Capabilities.Equals(other.Capabilities);
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
            return HashCode.Combine(States, InputCommands, Commands, Capabilities);
        }
    }
}