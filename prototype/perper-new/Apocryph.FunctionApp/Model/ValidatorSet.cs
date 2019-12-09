using System.Collections.Generic;
using System.Linq;

namespace Apocryph.FunctionApp.Model
{
    public class ValidatorSet
    {
        public Dictionary<string, int> Weights { get; set; }
        public Dictionary<string, int> AccumulatedWeights { get; set; }
        public int Total { get; set; }

        public string GetMaxAccumulatedWeight()
        {
            return AccumulatedWeights.Select(kv => (kv.Value, kv.Key)).Max().Item2;
        }

        public string PopMaxAccumulatedWeight()
        {
            var maxAccumulatedWeight = GetMaxAccumulatedWeight();
            AccumulatedWeights[maxAccumulatedWeight] -= Total;
            return maxAccumulatedWeight;
        }

        public void AccumulateWeights()
        {
            foreach (var kv in Weights)
            {
                if (!AccumulatedWeights.ContainsKey(kv.Key))
                {
                    AccumulatedWeights[kv.Key] = 0;
                }
                AccumulatedWeights[kv.Key] += kv.Value;
            }
        }
    }
}