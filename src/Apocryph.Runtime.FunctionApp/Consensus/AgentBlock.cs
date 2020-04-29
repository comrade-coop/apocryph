using System.Collections.Generic;
using Apocryph.Agent;
using Ipfs;

namespace Apocryph.Runtime.FunctionApp.Consensus
{
    public class AgentBlock
    {
        public object State { get; set; }
        public string Sender { get; set; }
        public object Message { get; set; }
        public AgentCommand[] Commands { get; set; }

        public Cid Previous { get; set; }
    }
}