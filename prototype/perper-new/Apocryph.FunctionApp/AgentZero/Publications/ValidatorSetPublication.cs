using System.Collections.Generic;
using System.Numerics;
using Apocryph.FunctionApp.Model;

namespace Apocryph.FunctionApp.AgentZero.Publications
{
    public class ValidatorSetPublication
    {
        public Dictionary<string, BigInteger> Weights { get; set; } // ValidatorKey

        public string AgentId { get; set; }
    }
}