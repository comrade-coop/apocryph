using System.Numerics;

namespace Apocryph.FunctionApp.AgentZero.Messages
{
    public class StakeMessage
    {
        public BigInteger Amount { get; set; }

        public string To { get; set; }
    }
}