using System;
using System.Collections.Generic;
using Apocryph.FunctionApp.Command;

namespace Apocryph.FunctionApp.Agent
{
    public class AgentContext<T> : IAgentContext<T>
    {
        public List<ICommand> Commands { get; } = new List<ICommand>();

        public AgentContext(T state)
        {
            State = state;
        }

        public T State { get; }


        public void RequestCallTicket(AgentCapability agent)
        {
            throw new NotImplementedException();
        }

        public void SendMessage(AgentCapability receiver, object message, AgentCallTicket callTicket)
        {
            throw new NotImplementedException();
        }

        public void IssueCapability(AgentCapability receiver, string[] messageTypes, AgentCallTicket callTicket)
        {
            throw new NotImplementedException();
        }

        public void RevokeCapability(AgentCapability agent)
        {
            throw new NotImplementedException();
        }

        public void AddReminder(TimeSpan time, object data)
        {
            Commands.Add(new ReminderCommand{Time = time, Data = data});
        }

        public void MakePublication(object payload)
        {
            Commands.Add(new PublicationCommand{Payload = payload});
        }

        public void AddSubscription(string target)
        {
            Commands.Add(new SubscriptionCommand{Target = target});
        }

        public void SendServiceMessage(string service, object parameters)
        {
            Commands.Add(new ServiceCommand{Service = service, Parameters = parameters});
        }
    }
}