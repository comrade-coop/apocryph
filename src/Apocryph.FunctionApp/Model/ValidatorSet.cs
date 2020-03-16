using System;
using System.Collections.Generic;
using System.Linq;

namespace Apocryph.FunctionApp.Model
{
    public class ValidatorSet
    {
        public Dictionary<ValidatorKey, int> Weights { get; set; } = new Dictionary<ValidatorKey, int>();

        private int? _total = null;
        public int Total { get => (_total = _total ?? Weights.Values.Sum()) ?? 0; }

        public bool IsMoreThanTwoThirds(IEnumerable<ValidatorKey> distinctVoters)
        {
            var voted = distinctVoters.Select(signer => Weights[signer]).Sum();
            return 3 * voted > 2 * Total;
        }

        public void InvalidateTotal()
        {
            _total = null;
        }
    }
}