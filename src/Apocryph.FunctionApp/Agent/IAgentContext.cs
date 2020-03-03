using System;

namespace Apocryph.FunctionApp.Agent
{
    public interface IAgentContext
    {
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