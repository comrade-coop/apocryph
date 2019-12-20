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
    public class BalancesState
    {
        public IDictionary<string, BigInteger> Amounts { get; set; } = new Dictionary<string, BigInteger>();

        public void RemoveTokens(string from, BigInteger amount)
        {
            if (!Amounts.ContainsKey(from) || Amounts[from] < amount)
            {
                throw new Exception("Not enough funds");
            }

            Amounts[from] -= amount;

            if (Amounts[from] == 0)
            {
                Amounts.Remove(from);
            }
        }

        public void AddTokens(string to, BigInteger amount)
        {
            if (!Amounts.ContainsKey(to))
            {
                Amounts[to] = 0;
            }

            Amounts[to] += amount;
        }
    }
}