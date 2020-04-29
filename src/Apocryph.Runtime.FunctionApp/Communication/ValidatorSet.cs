using System;
using System.Collections.Generic;
using Apocryph.Agent;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Ipfs;

namespace Apocryph.Runtime.FunctionApp.Communication
{
    // HACK: Needed only so that code compiles while worked on
    [Obsolete("ValidatorSet is needed only so that code compiles")]
    public class ValidatorSet
    {
        public Dictionary<ValidatorKey, int> Weights { get; set; }
        public bool IsMoreThanTwoThirds(IEnumerable<ValidatorKey> keys)
        {
            return false;
        }
    }
}