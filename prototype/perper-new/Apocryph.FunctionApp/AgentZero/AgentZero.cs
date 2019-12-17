using System;
using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.AgentZero.Messages;
using Apocryph.FunctionApp.AgentZero.Publications;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;

namespace Apocryph.FunctionApp.AgentZero
{
    public static class AgentZero
    {
        public class State
        {
            public IDictionary<string, BigInteger> Balances { get; set; } = new Dictionary<string, BigInteger>();
            public IDictionary<string, IDictionary<string, BigInteger>> Stakes { get; set; } = new Dictionary<string, IDictionary<string, BigInteger>>();

            public void RemoveTokens(string from, BigInteger amount)
            {
                if (!Balances.ContainsKey(from) || Balances[from] < amount)
                {
                    throw new Exception("Not enough funds");
                }

                Balances[from] -= amount;

                if (Balances[from] == 0)
                {
                    Balances.Remove(from);
                }
            }

            public void AddTokens(string to, BigInteger amount)
            {
                if (!Balances.ContainsKey(to))
                {
                    Balances[to] = 0;
                }

                Balances[to] += amount;
            }

            public void RemoveStake(string staker, string stakee, BigInteger amount)
            {
                if (!Stakes.ContainsKey(stakee) || !Stakes[stakee].ContainsKey(staker) || Stakes[stakee][staker] < amount)
                {
                    throw new Exception("Not enough stake");
                }

                Stakes[stakee][staker] -= amount;

                if (Stakes[stakee][staker] == 0)
                {
                    Stakes[stakee].Remove(staker);
                }

                if (Stakes[stakee].Count == 0)
                {
                    Stakes.Remove(stakee);
                }
            }

            public void AddStake(string staker, string stakee, BigInteger amount)
            {
                if (!Stakes.ContainsKey(stakee))
                {
                    Stakes[stakee] = new Dictionary<string, BigInteger>();
                }

                if (!Stakes[stakee].ContainsKey(staker))
                {
                    Stakes[stakee][staker] = 0;
                }

                Stakes[stakee][staker] += amount;
            }
        }

        public static void Run(IAgentContext<State> context, string sender, object message)
        {
            switch (message)
            {
                case TransferMessage transferMessage:
                    context.State.RemoveTokens(sender, transferMessage.Amount);
                    context.State.AddTokens(transferMessage.To, transferMessage.Amount);
                    context.MakePublication(new TransferPublication
                    {
                        From = sender,
                        To = transferMessage.To,
                        Amount = transferMessage.Amount,
                    });
                    break;

                case StakeMessage stakeMessage:
                    context.State.RemoveTokens(sender, stakeMessage.Amount);
                    context.State.AddStake(sender, stakeMessage.To, stakeMessage.Amount);
                    context.MakePublication(new StakePublication
                    {
                        From = sender,
                        To = stakeMessage.To,
                        Amount = stakeMessage.Amount,
                    });
                    break;

                case UnstakeMessage unstakeMessage:
                    context.State.RemoveStake(sender, unstakeMessage.From, unstakeMessage.Amount);
                    context.State.AddTokens(sender, unstakeMessage.Amount);
                    context.MakePublication(new StakePublication
                    {
                        From = sender,
                        To = unstakeMessage.From,
                        Amount = -unstakeMessage.Amount,
                    });
                    break;

                case RegisterAgentMessage registerAgentMessage:
                    context.State.RemoveTokens(sender, registerAgentMessage.InitialBalance);
                    context.State.AddTokens(registerAgentMessage.AgentId, registerAgentMessage.InitialBalance);
                    context.MakePublication(new ValidatorSetPublication
                    {
                        AgentId = registerAgentMessage.AgentId,
                        Weights = context.State.Stakes.ToDictionary(kv => kv.Key, kv => kv.Value.Values.Aggregate(BigInteger.Add)),
                    });
                    break;
            }

        }
    }
}