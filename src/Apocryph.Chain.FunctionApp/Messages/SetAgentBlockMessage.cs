using System.Numerics;

namespace Apocryph.Chain.FunctionApp.Messages
{
    public class SetAgentBlockMessage
    {
        public string AgentId { get; set; }
        public byte[] BlockId { get; set; }
    }
}