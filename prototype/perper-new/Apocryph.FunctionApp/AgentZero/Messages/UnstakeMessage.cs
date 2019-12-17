using System.Numerics;

namespace Apocryph.FunctionApp.AgentZero.Messages
{
    public class UnstakeMessage
    {
        public BigInteger Amount { get; set; }

        public string From { get; set; }
    }
}