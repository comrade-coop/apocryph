using System;
using System.Collections.Generic;
using Apocryph.Core.Consensus.VirtualNodes;

namespace Apocryph.Core.Consensus.Blocks
{
    public class Block : IEquatable<Block>
    {
        public Guid ChainId { get; set; }
        public Guid ProposerAccount { get; set; }
        public Node? Proposer { get; set; }
        public object[] InputCommands { get; }
        public object[] Commands { get; }
        public Dictionary<string, byte[]> States { get; }
        public Dictionary<Guid, (string, string[])> Capabilities { get; }

        public Block(Guid chainId, Guid proposerAccount, Dictionary<string, byte[]> states, object[] inputCommands, object[] commands, Dictionary<Guid, (string, string[])> capabilities)
        {
            ChainId = chainId;
            ProposerAccount = proposerAccount;
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