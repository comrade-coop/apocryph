using System;

namespace Apocryph.FunctionApp.Agent
{
    public interface IAgentContext
    {
        void RequestCallToken(string target);
        void SendMessage(AgentCapability receiver, object message, AgentCallToken callToken);

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