using System;
using System.Collections.Generic;
using System.Net.WebSockets;

namespace Apocryph.Agents.Testbed.Api
{
    public class AgentContext
    {
        public object State { get; }

        private readonly AgentCapability _self;
        private readonly List<AgentCommand> _commands;

        public AgentContext(object state, AgentCapability self)
        {
            State = state;

            _self = self;
            _commands = new List<AgentCommand>();
        }

        public AgentCommands GetCommands()
        {
            return new AgentCommands {State = State, Commands = _commands.ToArray()};
        }

        public AgentCapability IssueCapability(Type[] messageTypes)
        {
            var result = new AgentCapability {Issuer = _self.Issuer, MessageTypes = messageTypes};
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
            _commands.Add(new AgentCommand{CommandType = AgentCommandType.SendMessage, Receiver = receiver, Message = message});
        }

        public void AddReminder(TimeSpan time, object data)
        {
        }

        public void MakePublication(object payload)
        {
        }

        public void AddSubscription(string target)
        {
        }

        public void SendServiceMessage(string service, object parameters)
        {
        }
    }
}