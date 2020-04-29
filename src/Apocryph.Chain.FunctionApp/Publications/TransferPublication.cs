using System.Numerics;

namespace Apocryph.Chain.FunctionApp.Publications
{
    public class TransferPublication
    {
        public BigInteger Amount { get; set; }

        public string From { get; set; }

        public string To { get; set; }
    }
}