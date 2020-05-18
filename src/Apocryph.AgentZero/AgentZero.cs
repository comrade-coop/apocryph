using System;
using System.Threading.Tasks;
using Apocryph.Agent;
using Apocryph.AgentZero.Messages;
using Apocryph.AgentZero.Publications;
using Apocryph.AgentZero.State;

namespace Apocryph.AgentZero
{
    public class AgentZeroState
    {
        public BalancesState Balances { get; set; } = new BalancesState();
        public AgentsState Agents { get; set; } = new AgentsState();
    }

    public class AgentZero : IAgent<AgentZeroState>
    {
        public void Setup(IContext<AgentZeroState> context)
        {
            context.RegisterInstance<RegisterAgentMessage, RegisterAgentMessage>();
            context.RegisterInstance<SetAgentBlockMessage, SetAgentBlockMessage>();
            context.RegisterInstance<TransferMessage, TransferMessage>();
        }

        public Task Run(IContext<AgentZeroState> context, object message, Guid? reference)
        {
            switch (message)
            {
                case TransferMessage transferMessage:
                    context.State!.Balances.RemoveTokens(reference!.ToString()!, transferMessage.Amount);
                    context.State!.Balances.AddTokens(transferMessage.To, transferMessage.Amount);
                    context.Publish(new TransferPublication
                    {
                        From = reference!.ToString()!,
                        To = transferMessage.To,
                        Amount = transferMessage.Amount,
                    });
                    break;

                case RegisterAgentMessage registerAgentMessage:
                    context.State!.Balances.RemoveTokens(reference!.ToString()!, registerAgentMessage.InitialBalance);
                    context.State!.Balances.AddTokens(registerAgentMessage.AgentId, registerAgentMessage.InitialBalance);
                    context.State!.Agents.SetAgentBlock(registerAgentMessage.AgentId, registerAgentMessage.BlockId);
                    break;

                case SetAgentBlockMessage setAgentBlockMessage:
                    context.State!.Agents.SetAgentBlock(setAgentBlockMessage.AgentId, setAgentBlockMessage.BlockId);
                    break;
            }

            return Task.FromResult(context);
        }
    }
}