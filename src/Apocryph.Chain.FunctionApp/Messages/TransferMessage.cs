using System.Numerics;

namespace Apocryph.Chain.FunctionApp.Messages
{
    public class TransferMessage
    {
        public BigInteger Amount { get; set; }

        public string To { get; set; }
    }
}