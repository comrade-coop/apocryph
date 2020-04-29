using Apocryph.Agent;

namespace SampleAgents.FunctionApp.Agents
{
    public class PingPongMessage
    {
        public AgentCapability AgentOne { get; set; }
        public AgentCapability AgentTwo { get; set; }

        public string Content { get; set; }
    }
}