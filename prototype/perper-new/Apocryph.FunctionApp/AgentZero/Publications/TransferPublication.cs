using System.Numerics;

namespace Apocryph.FunctionApp.AgentZero.Publications
{
    public class TransferPublication
    {
        public BigInteger Amount { get; set; }

        public string From { get; set; }

        public string To { get; set; }
    }
}