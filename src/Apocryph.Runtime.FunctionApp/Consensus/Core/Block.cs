using System;
using Apocryph.Agent;
using Ipfs;

namespace Apocryph.Runtime.FunctionApp.Consensus.Core
{
    public class Block : IEquatable<Block>
    {
        public object State { get; set; }
        public string Sender { get; set; }
        public object Message { get; set; }
        public AgentCommand[] Commands { get; set; }

        public int ProposerStake { get; set; }

        public Cid Previous { get; set; }

        public bool Equals(Block? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return State.Equals(other.State) && Sender == other.Sender && Message.Equals(other.Message) && Commands.Equals(other.Commands) && Previous.Equals(other.Previous);
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
            return HashCode.Combine(State, Sender, Message, Commands, Previous);
        }
    }
}