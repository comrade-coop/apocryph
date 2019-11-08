using System;
using System.Collections.Generic;
using System.Linq;
using Akka.Actor;

namespace Apocryph.Prototype
{
    public class ValidatorSet {
        public ValidatorSet(Dictionary<string, int> validatorStakes)
        {
            Stakes = validatorStakes;
            AccumulatedStakes = new Dictionary<string, int>();
            Proposer = null;
        }

        public Dictionary<string, int> Stakes { get; }
        public Dictionary<string, int> AccumulatedStakes { get; }
        public string Proposer { get; private set; }

        public void SetValidatorStake(string key, int stake)
        {
            if (stake != 0)
            {
                Stakes[key] = stake;
            }
            else
            {
                Stakes.Remove(key);
            }
        }

        public void AdvanceProposer()
        {
            var totalStake = 0;

            foreach (var key in Stakes.Keys)
            {
                if (!AccumulatedStakes.ContainsKey(key))
                {
                    AccumulatedStakes[key] = 0;
                }
                AccumulatedStakes[key] += Stakes[key];
                totalStake += Stakes[key];
            }

            var maxProposer = Stakes.Keys.First();
            var maxStake = int.MinValue;

            foreach (var key in Stakes.Keys)
            {
                if (AccumulatedStakes[key] > maxStake || (AccumulatedStakes[key] == maxStake && string.Compare(key, maxProposer) > 0))
                {
                    maxProposer = key;
                    maxStake = AccumulatedStakes[key];
                }
            }

            AccumulatedStakes[maxProposer] -= totalStake;
            Proposer = maxProposer;
        }
    }
}
