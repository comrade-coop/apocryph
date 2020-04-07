using System;

namespace Apocryph.Agents.Testbed.Api
{
    public class AgentCapability
    {
        public string Issuer { get; set; }
        public Type[] MessageTypes { get; set; }
    }
}