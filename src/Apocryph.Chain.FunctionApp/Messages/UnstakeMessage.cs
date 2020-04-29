using System.Numerics;

namespace Apocryph.Chain.FunctionApp.Messages
{
    public class UnstakeMessage
    {
        public BigInteger Amount { get; set; }

        public string From { get; set; }
    }
}