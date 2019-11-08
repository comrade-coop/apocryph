using System;

namespace Apocryph.Prototype
{
    public class Message
    {
        public Message(string origin, string sender, string hash, MessagePayload payload)
        {
            Origin = origin;
            Sender = sender;
            Hash = hash;
            Payload = payload;
        }

        public string Origin { get; }
        public string Sender { get; }
        public string Hash { get; }
        public MessagePayload Payload { get; }
    }
}
