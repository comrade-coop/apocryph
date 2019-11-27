using System.Numerics;
using Apocryph.FunctionApp.Model.Message;

namespace Apocryph.FunctionApp.Wetonomy.Messages
{
    public class TokensChangedMessage : IMessage //Message or Publication???
    {
        public string Target { get; set; }
        
        public string TokenManager { get; set; }

        public BigInteger Change { get; set; }

        public BigInteger Total { get; set; }
    }
}