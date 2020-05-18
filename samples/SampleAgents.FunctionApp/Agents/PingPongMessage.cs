using System;

namespace SampleAgents.FunctionApp.Agents
{
    public class PingPongMessage : IPingPongMessage
    {
        public Guid? AgentOne { get; set; }
        public Guid? AgentTwo { get; set; }
        public string? Content { get; set; }
    }
}