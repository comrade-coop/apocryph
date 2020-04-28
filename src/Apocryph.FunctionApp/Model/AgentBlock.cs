using System.Collections.Generic;
using Ipfs;

namespace Apocryph.FunctionApp.Model
{
    public class AgentBlock
    {
        public object State { get; set; }
        public string Sender { get; set; }
        public object Message { get; set; }
        public List<ICommand> Commands { get; set; }

        public Cid Previous { get; set; }
    }
}