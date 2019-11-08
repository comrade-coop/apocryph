using System;
using System.Collections.Generic;

namespace Apocryph.Prototype
{
    public interface Contract
    {
        IEnumerable<MessagePayload> Receive(Message message);
    }
}
