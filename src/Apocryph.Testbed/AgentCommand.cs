using System;

namespace Apocryph.Testbed
{
    public class AgentCommand
    {
        public AgentCommandType CommandType { get; set; }

        public string AgentId { get; set; }
        public string Agent { get; set; }

        public AgentCapability Receiver { get; set; }

        public object Message { get; set; }

        public TimeSpan Timeout { get; set; }
    }

    public enum AgentCommandType
    {
        CreateAgent,
        SendMessage,
        Publish,
        Subscribe,
        Reminder
    }
}