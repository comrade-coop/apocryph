using System;
using System.Collections.Generic;
using System.Threading.Channels;

namespace Apocryph.Agent
{
    public class AgentContext<T> : AgentContext where T : class
    {
        public T State => (T)InternalState;

        public AgentContext(T state, AgentCapability self) : base(state, self)
        {
        }
    }

    public class AgentContext
    {
        public object InternalState { get; }

        private readonly AgentCapability _self;
        private readonly List<AgentCommand> _commands;

        public AgentContext(object internalState, AgentCapability self)
        {

            InternalState = internalState;

            _self = self;
            _commands = new List<AgentCommand>();
        }

        public AgentCommands GetCommands()
        {
            return new AgentCommands { Origin = _self.Issuer, State = InternalState, Commands = _commands.ToArray() };
        }

        public AgentCapability IssueCapability(Type[] messageTypes)
        {
            var result = new AgentCapability { Issuer = _self.Issuer, MessageTypes = messageTypes };
            return result;
        }

        public void RevokeCapability(AgentCapability capability)
        {
        }

        public AgentCallTicket RequestCallTicket(AgentCapability agent)
        {
            return null;
        }

        public void CreateAgent(string id, string agent, object initMessage, AgentCallTicket callTicket)
        {
            _commands.Add(new AgentCommand
            {
                CommandType = AgentCommandType.CreateAgent,
                AgentId = id,
                Agent = agent,
                Message = initMessage
            });
        }

        public void SendMessage(AgentCapability receiver, object message, AgentCallTicket callTicket)
        {
            _commands.Add(new AgentCommand { CommandType = AgentCommandType.SendMessage, Receiver = receiver, Message = message });
        }

        public void AddReminder(TimeSpan time, object data)
        {
            _commands.Add(new AgentCommand { CommandType = AgentCommandType.Reminder, Receiver = _self, Timeout = time, Message = data });
        }

        public void MakePublication(object payload)
        {
            _commands.Add(new AgentCommand { CommandType = AgentCommandType.Publish, Message = payload });
        }

        public void AddSubscription(string target)
        {
            _commands.Add(new AgentCommand { CommandType = AgentCommandType.Subscribe, AgentId = target });
        }
    }
}