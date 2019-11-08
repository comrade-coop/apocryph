using System;
using System.Threading.Tasks;
using System.Collections.Generic;
using System.Linq;
using Akka.Actor;

namespace Apocryph.Prototype
{
    class Program
    {
        static void Main(string[] args)
        {
            using (var system = ActorSystem.Create("apocryph-prototype"))
            {
                var validatorStakes = new Dictionary<string, int>()
                {
                    ["Na"] = 1,
                    ["Nb"] = 1,
                };

                var broadcaster = (ICanTell) system.ActorSelection("user/broadcast");
                var cores = new List<ICanTell>();

                foreach (var key in validatorStakes.Keys) {
                    cores.Add(system.ActorOf(Props.Create<CoreService>(
                        broadcaster,
                        validatorStakes
                    ), key));
                }

                broadcaster = system.ActorOf(Props.Create<Broadcaster>(cores), "broadcast");

                broadcaster.Tell(new CoreService.Transaction("_", 0, new SampleContract.SampleMessagePayload(
                    "B",
                    new string[] { "C", "A" },
                    0)), system.DeadLetters);

                Console.ReadLine();
            }
        }
    }
}
