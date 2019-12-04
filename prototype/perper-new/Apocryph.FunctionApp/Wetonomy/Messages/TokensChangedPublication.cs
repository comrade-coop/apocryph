using System.Numerics;

namespace Apocryph.FunctionApp.Wetonomy.Messages
{
    public class TokensChangedPublication
    {
        public string Target { get; set; }
        
        public string TokenManager { get; set; }

        public BigInteger Change { get; set; }

        public BigInteger Total { get; set; }
    }
}