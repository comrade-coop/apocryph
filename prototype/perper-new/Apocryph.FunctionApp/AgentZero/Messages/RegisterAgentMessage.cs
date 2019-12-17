using System.Numerics;

namespace Apocryph.FunctionApp.AgentZero.Messages
{
    public class RegisterAgentMessage
    {
        public BigInteger InitialBalance { get; set; }

        public string AgentId { get; set; }

        // public int ExpectedValidators { get; set; }
    }
}