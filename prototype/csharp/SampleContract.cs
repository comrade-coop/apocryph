using System;
using System.Collections.Generic;
using System.Linq;

namespace Apocryph.Prototype
{
    public class SampleContract : Contract
    {
        public class SampleMessagePayload : MessagePayload
        {
            public SampleMessagePayload(string target, string[] nextTargets, int minCount)
                : base(target)
            {
                NextTargets = nextTargets;
                MinCount = minCount;
            }

            public string[] NextTargets { get; }
            public int MinCount { get; }
        }

        public int MessagesCount { get; set; } = 0;

        public IEnumerable<MessagePayload> Receive(Message message)
        {
            MessagesCount += 1;
            if (message.Payload is SampleMessagePayload smp)
            {
                if (smp.NextTargets.Length > 0 && smp.MinCount <= MessagesCount)
                {
                    yield return new SampleMessagePayload(
                        smp.NextTargets[0],
                        smp.NextTargets.Skip(1).ToArray(),
                        smp.MinCount);
                }
            }
        }
    }
}
