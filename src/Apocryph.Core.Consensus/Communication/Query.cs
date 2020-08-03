using Apocryph.Core.Consensus.VirtualNodes;

namespace Apocryph.Core.Consensus.Communication
{
    public class Query<T>
    {
        public T Value { get; }
        public Node Sender { get; }
        public Node Receiver { get; }

        public int Round { get; }
        public QueryVerb Verb { get; }

        public Query(T value, Node sender, Node receiver, int round, QueryVerb verb)
        {
            Value = value;
            Sender = sender;
            Receiver = receiver;
            Round = round;
            Verb = verb;
        }
    }

    public enum QueryVerb
    {
        Request,
        Response
    }
}