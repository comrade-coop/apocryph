using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Akka.Actor;

namespace Apocryph.Prototype
{
    public class Broadcaster : UntypedActor
    {
        private List<ICanTell> _targets;

        public Broadcaster(List<ICanTell> targets) {
            _targets = targets;
        }

        protected override void OnReceive(object received)
        {
            foreach (var target in _targets) {
                if (target != Context.Sender) {
                    target.Tell(received, Context.Self);
                }
            }
        }
    }
}
