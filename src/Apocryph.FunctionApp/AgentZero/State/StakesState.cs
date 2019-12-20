using System;
using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.AgentZero.Messages;
using Apocryph.FunctionApp.AgentZero.Publications;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;

namespace Apocryph.FunctionApp.AgentZero.State
{
    public class StakesState
    {
        public IDictionary<string, IDictionary<string, BigInteger>> Amounts { get; set; } = new Dictionary<string, IDictionary<string, BigInteger>>();

        public void RemoveStake(string staker, string stakee, BigInteger amount)
        {
            if (!Amounts.ContainsKey(stakee) || !Amounts[stakee].ContainsKey(staker) || Amounts[stakee][staker] < amount)
            {
                throw new Exception("Not enough stake");
            }

            Amounts[stakee][staker] -= amount;

            if (Amounts[stakee][staker] == 0)
            {
                Amounts[stakee].Remove(staker);
            }

            if (Amounts[stakee].Count == 0)
            {
                Amounts.Remove(stakee);
            }
        }

        public void AddStake(string staker, string stakee, BigInteger amount)
        {
            if (!Amounts.ContainsKey(stakee))
            {
                Amounts[stakee] = new Dictionary<string, BigInteger>();
            }

            if (!Amounts[stakee].ContainsKey(staker))
            {
                Amounts[stakee][staker] = 0;
            }

            Amounts[stakee][staker] += amount;
        }
    }
}