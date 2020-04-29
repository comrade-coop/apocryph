using System;

namespace Apocryph.Agent
{
    public class AgentCapability
    {
        public string Issuer { get; set; }
        public Type[] MessageTypes { get; set; }
    }
}