using System.Numerics;

namespace Apocryph.FunctionApp.AgentZero.Messages
{
    public class TransferMessage
    {
        public BigInteger Amount { get; set; }

        public string To { get; set; }
    }
}