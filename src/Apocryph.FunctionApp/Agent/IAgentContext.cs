using System;

namespace Apocryph.FunctionApp.Agent
{
    public interface IAgentContext
    {
        void RequestCallTicket(AgentCapability agent);
        void SendMessage(AgentCapability receiver, object message, AgentCallTicket callTicket);
        void IssueCapability(AgentCapability receiver, string[] messageTypes, AgentCallTicket callTicket);
        void RevokeCapability(AgentCapability agent);

        void AddReminder(TimeSpan time, object data);
        void MakePublication(object payload);
        void AddSubscription(string target);

        void SendServiceMessage(string command, object parameters);
    }

    public interface IAgentContext<out T> : IAgentContext
    {
        T State { get; }
    }
}