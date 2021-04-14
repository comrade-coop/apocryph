using System;
using System.Linq;
using Apocryph.Ipfs;

namespace Apocryph.Consensus
{
    public class Reference : IEquatable<Reference>
    {
        public Hash<Chain> Chain { get; private set; }
        public int AgentNonce { get; private set; }
        public string[] AllowedMessageTypes { get; private set; }
        // public MerkleTreeProof<Message> Source { get; private set; }

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
            var hashCode = new HashCode();
            hashCode.Add(Chain);
            hashCode.Add(AgentNonce);
            Array.ForEach(AllowedMessageTypes, hashCode.Add);
            return hashCode.ToHashCode();
        }
    }
}