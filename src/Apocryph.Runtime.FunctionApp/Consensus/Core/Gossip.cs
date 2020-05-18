namespace Apocryph.Runtime.FunctionApp.Consensus.Core
{
    public class Gossip<T>
    {
        public T Value { get; }
        public Node Sender { get; }
        public GossipVerb Verb { get; }

        public Gossip(T value, Node sender, GossipVerb verb)
        {
            Value = value;
            Sender = sender;
            Verb = verb;
        }
    }

    public enum GossipVerb
    {
        // IdentityChanged,
        Confirm,
        Reject
    }
}