using System;
using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using Apocryph.Agent;
using Apocryph.AgentZero.Messages;
using Apocryph.AgentZero.Publications;
using Apocryph.AgentZero.State;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.AgentZero
{
    public static class AgentZeroWorker
    {
        public class State
        {
            public BalancesState Balances { get; set; } = new BalancesState();
            public AgentsState Agents { get; set; } = new AgentsState();
        }

        [FunctionName(nameof(AgentZeroWorker))]
        [return: Perper("$return")]
        public static AgentContext<State> Run([PerperWorkerTrigger] object workerContext,
            [Perper("self")] AgentCapability self,
            [Perper("state")] State state,
            [Perper("sender")] string sender,
            [Perper("message")] object message)
        {
            var context = new AgentContext<State>(state ?? new State(), self);
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

                case RegisterAgentMessage registerAgentMessage:
                    context.State.Balances.RemoveTokens(sender, registerAgentMessage.InitialBalance);
                    context.State.Balances.AddTokens(registerAgentMessage.AgentId, registerAgentMessage.InitialBalance);
                    context.State.Agents.SetAgentBlock(registerAgentMessage.AgentId, registerAgentMessage.BlockId);
                    break;

                case SetAgentBlockMessage setAgentBlockMessage:
                    context.State.Agents.SetAgentBlock(setAgentBlockMessage.AgentId, setAgentBlockMessage.BlockId);
                    break;
            }
            return context;
        }
    }
}