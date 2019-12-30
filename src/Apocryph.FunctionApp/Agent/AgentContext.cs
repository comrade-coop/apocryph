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

        public void IssueToken(string recipient)
        {
            throw new NotImplementedException();
        }

        public void SendMessage(string target, object message)
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
    }
}