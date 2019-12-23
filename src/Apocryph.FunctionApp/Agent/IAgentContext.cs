using System;

namespace Apocryph.FunctionApp.Agent
{
    public interface IAgentContext<out T>
    {
        T State { get; }

        // void IssueToken(string recipient);
        // void SendMessage(string target, object message);
        void AddReminder(TimeSpan time, object? data = null);
        void MakePublication(object payload);
        void AddSubscription(string target);
    }
}