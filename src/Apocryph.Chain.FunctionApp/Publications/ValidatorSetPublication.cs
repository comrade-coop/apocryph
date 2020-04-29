using System.Collections.Generic;
using System.Numerics;

namespace Apocryph.Chain.FunctionApp.Publications
{
    public class ValidatorSetPublication
    {
        public Dictionary<string, BigInteger> Weights { get; set; }

        public string AgentId { get; set; }
    }
}