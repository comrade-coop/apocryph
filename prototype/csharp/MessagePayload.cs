using System;

namespace Apocryph.Prototype
{
    public class MessagePayload
    {
        public MessagePayload(string target)
        {
            Target = target;
        }

        public string Target { get; }
    }
}
