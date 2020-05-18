using System.Collections.Generic;

namespace Apocryph.AgentZero.State
{
    public class AgentsState
    {
        public IDictionary<string, byte[]> Agents { get; set; } = new Dictionary<string, byte[]>();

        public void SetAgentBlock(string agentId, byte[] block)
        {
            Agents[agentId] = block;
        }
    }
}