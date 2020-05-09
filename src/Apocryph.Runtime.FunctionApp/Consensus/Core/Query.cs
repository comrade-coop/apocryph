using System;

namespace Apocryph.Runtime.FunctionApp.Consensus.Core
{
    public class Query<T>
    {
        public T Value { get; }
        public Node Sender { get; }
        public Node Receiver { get; }

        public QueryVerb Verb { get; }

        public Query(T value, Node sender, Node receiver, QueryVerb verb)
        {
            Value = value;
            Sender = sender;
            Receiver = receiver;
            Verb = verb;
        }
    }

    public enum QueryVerb
    {
        Request,
        Response,
        Accept
    }
}