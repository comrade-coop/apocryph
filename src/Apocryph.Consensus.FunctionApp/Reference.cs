using System;
using System.Linq;
using Apocryph.HashRegistry;

namespace Apocryph.Consensus
{
    public class Reference : IEquatable<Reference>
    {
        public Hash<Chain> Chain { get; }
        public int AgentNonce { get; }
        public string[] AllowedMessageTypes { get; }
        // public MerkleTreeProof<Message> Source { get; }

        public Reference(Hash<Chain> chain, int agentNonce, string[] allowedMessageTypes)
        {
            Chain = chain;
            AgentNonce = agentNonce;
            AllowedMessageTypes = allowedMessageTypes;
        }

        public override bool Equals(object? other)
        {
            return other is Reference otherReference && Equals(otherReference);
        }

        public bool Equals(Reference? other)
        {
            return other != null && Chain.Equals(other.Chain) && AgentNonce.Equals(other.AgentNonce) && AllowedMessageTypes.SequenceEqual(other.AllowedMessageTypes);
        }

        public override int GetHashCode()
        {
            return HashCode.Combine(Chain, AgentNonce, AllowedMessageTypes);
        }
    }
}