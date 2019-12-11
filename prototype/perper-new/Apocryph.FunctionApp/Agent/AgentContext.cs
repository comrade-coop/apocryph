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
            Commands.Add(new SendMessageCommand { Target = target, Message = message});
        }

        public void AddReminder(TimeSpan time)
        {
            Commands.Add(new ReminderCommand{Time = time});
        }

        public void MakePublication(object payload)
        {
            throw new NotImplementedException();
        }
    }
}