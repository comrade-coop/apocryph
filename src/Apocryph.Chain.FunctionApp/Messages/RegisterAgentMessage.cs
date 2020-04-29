using System.Numerics;

namespace Apocryph.Chain.FunctionApp.Messages
{
    public class RegisterAgentMessage
    {
        public BigInteger InitialBalance { get; set; }

        public string AgentId { get; set; }

        // public int ExpectedValidators { get; set; }
    }
}