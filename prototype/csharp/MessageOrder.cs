using System;

namespace Apocryph.Prototype
{
    public class MessageOrder
    {
        public MessageOrder(string target, string hash)
        {
            Target = target;
            Hash = hash;
        }

        public string Target { get; }
        public string Hash { get; }
    }
}
