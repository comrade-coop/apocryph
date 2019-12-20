using System;

namespace Apocryph.FunctionApp.Agent
{
    public interface IAgentContext<out T>
    {
        T State { get; }
        
        void IssueToken(string recipient);
        void SendMessage(string target, object message);
        void AddReminder(TimeSpan time);
        void MakePublication(object payload);
    }
}