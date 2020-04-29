using System;
using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using Apocryph.Agent;
using Apocryph.Chain.FunctionApp.Messages;
using Apocryph.Chain.FunctionApp.Publications;
using Apocryph.Chain.FunctionApp.State;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Chain.FunctionApp
{
    public static class AgentZeroWorker
    {
        public class State
        {
            public BalancesState Balances { get; set; } = new BalancesState();
            public StakesState Stakes { get; set; } = new StakesState();
            public AgentsState Agents { get; set; } = new AgentsState();
        }

        [FunctionName(nameof(AgentZeroWorker))]
        [return: Perper("$return")]
        public static AgentContext<State> Run([PerperWorkerTrigger] object workerContext,
            [Perper("state")] State state,
            [Perper("sender")] string sender,
            [Perper("message")] object message)
        {
            var context = new AgentContext<State>(state ?? new State());
            switch (message)
            {
                case TransferMessage transferMessage:
                    context.State.Balances.RemoveTokens(sender, transferMessage.Amount);
                    context.State.Balances.AddTokens(transferMessage.To, transferMessage.Amount);
                    context.MakePublication(new TransferPublication
                    {
                        From = sender,
                        To = transferMessage.To,
                        Amount = transferMessage.Amount,
                    });
                    break;

                case StakeMessage stakeMessage:
                    context.State.Balances.RemoveTokens(sender, stakeMessage.Amount);
                    context.State.Stakes.AddStake(sender, stakeMessage.To, stakeMessage.Amount);
                    context.MakePublication(new StakePublication
                    {
                        From = sender,
                        To = stakeMessage.To,
                        Amount = stakeMessage.Amount,
                    });
                    break;

                case UnstakeMessage unstakeMessage:
                    context.State.Stakes.RemoveStake(sender, unstakeMessage.From, unstakeMessage.Amount);
                    context.State.Balances.AddTokens(sender, unstakeMessage.Amount);
                    context.MakePublication(new StakePublication
                    {
                        From = sender,
                        To = unstakeMessage.From,
                        Amount = -unstakeMessage.Amount,
                    });
                    break;

                case RegisterAgentMessage registerAgentMessage:
                    context.State.Balances.RemoveTokens(sender, registerAgentMessage.InitialBalance);
                    context.State.Balances.AddTokens(registerAgentMessage.AgentId, registerAgentMessage.InitialBalance);
                    context.State.Agents.RegisterAgent(registerAgentMessage.AgentId);
                    context.MakePublication(new ValidatorSetPublication
                    {
                        AgentId = registerAgentMessage.AgentId,
                        Weights = context.State.Stakes.Amounts.ToDictionary(kv => kv.Key, kv => kv.Value.Values.Aggregate(BigInteger.Add)),
                    });
                    break;
            }
            return context;
        }
    }
}