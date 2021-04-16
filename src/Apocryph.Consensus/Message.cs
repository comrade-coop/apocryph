using System;

namespace Apocryph.Consensus
{
    public class Message : IEquatable<Message>
    {
        // public Reference Source { get; private set; }
        public Reference Target { get; private set; }  // NOTE: Currently publications are encoded as negative Target.AgentNonce values
        public ReferenceData Data { get; private set; }
        // public DateTime SendTime { get; private set; }

        public Message(Reference target, ReferenceData data)
        {
            Target = target;
            Data = data;
        }

        public override bool Equals(object? other)
        {
            return other is Message otherMessage && Equals(otherMessage);
        }

        public bool Equals(Message? other)
        {
            return other != null && Target.Equals(other.Target) && Data.Equals(other.Data);
        }

        public override int GetHashCode()
        {
            return HashCode.Combine(Target, Data);
        }
    }
}