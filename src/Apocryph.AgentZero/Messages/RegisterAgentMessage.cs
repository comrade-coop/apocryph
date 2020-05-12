using System.Numerics;

namespace Apocryph.AgentZero.Messages
{
    public class RegisterAgentMessage : SetAgentBlockMessage
    {
        public BigInteger InitialBalance { get; set; }
    }
}