using System;
using System.Collections.Generic;
using Apocryph.FunctionApp.Command;

namespace Apocryph.FunctionApp.Agent
{
    public class AgentContext<T> : IAgentContext<T>
    {
        private readonly IList<ICommand> _commands;

        public AgentContext(IList<ICommand> commands)
        {
            _commands = commands;
        }
        
        public T State { get; }
        
        public void IssueToken(string recipient)
        {
            throw new NotImplementedException();
        }

        public void SendMessage(string target, object message)
        {
            _commands.Add(new SendMessageCommand { Target = target, Message = message});
        }

        public void AddReminder(TimeSpan time)
        {
            _commands.Add(new ReminderCommand{Time = time});
        }

        public void MakePublication(object payload)
        {
            throw new NotImplementedException();
        }
    }
}