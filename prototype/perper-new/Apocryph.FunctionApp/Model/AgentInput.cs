using System;

namespace Apocryph.FunctionApp.Model
{
    public class AgentInput : IAgentStep
    {
        public string Type { get; set; }
        
        public object State { get; set; }
        public string Sender { get; set; }
        public object Message { get; set; }
    }
}