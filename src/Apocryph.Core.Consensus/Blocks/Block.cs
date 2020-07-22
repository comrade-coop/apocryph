using System;
using System.Linq;
using System.Collections;
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

        public Block(Guid chainId, Node? proposer, Guid proposerAccount, Dictionary<string, byte[]> states, object[] inputCommands, object[] commands, Dictionary<Guid, (string, string[])> capabilities)
        {
            ChainId = chainId;
            Proposer = proposer;
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
            return ChainId == other.ChainId
                && ProposerAccount == other.ProposerAccount
                && (Proposer != null ? Proposer.Equals(other.Proposer) : other.Proposer == null)
                && States.Count() == other.States.Count() && States.All(kv => other.States.ContainsKey(kv.Key) && kv.Value.SequenceEqual(other.States[kv.Key]))
                && (InputCommands as IStructuralEquatable).Equals(other.InputCommands, EqualityComparer<object>.Default)
                && (Commands as IStructuralEquatable).Equals(other.Commands, EqualityComparer<object>.Default)
                && Capabilities.Count() == other.Capabilities.Count() && Capabilities.All(kv => other.Capabilities.ContainsKey(kv.Key) && (kv.Value as IStructuralEquatable).Equals(other.Capabilities[kv.Key], StructuralComparisons.StructuralEqualityComparer));
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
            var hash = new HashCode();
            hash.Add(ChainId);
            hash.Add(ProposerAccount);
            hash.Add(Proposer);
            foreach(var (key, state) in States)
            {
                hash.Add(key);
                Array.ForEach(state, hash.Add);
            }
            Array.ForEach(InputCommands, hash.Add);
            Array.ForEach(Commands, hash.Add);
            foreach(var (key, capability) in Capabilities)
            {
                hash.Add(key);
                hash.Add(capability.Item1);
                Array.ForEach(capability.Item2, hash.Add);
            }
            return hash.ToHashCode();
        }

        public override string ToString()
        {
            return $"Block({ChainId}, {Proposer}, {ProposerAccount}, InputCommands = [{string.Join(", ", InputCommands)}], Commands = [{string.Join(", ", Commands)}], States = [{string.Join(", ", States.Keys)}], Capabilities = [{string.Join(", ", Capabilities.Keys)}])";
        }
    }
}