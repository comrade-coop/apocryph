using System.Numerics;

namespace Apocryph.Chain.FunctionApp.Messages
{
    public class RegisterAgentMessage : SetAgentBlockMessage
    {
        public BigInteger InitialBalance { get; set; }
    }
}