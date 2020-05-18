using System;

namespace Apocryph.Agent
{
    public interface IContext<T> where T : class
    {
        T? State { get; set; }

        void RegisterInstance<TInterface, TClass>();
        TInterface CreateInstance<TInterface>(Action<TInterface> initializer = null);

        Guid CreateReference(Type[] messageTypes);

        void Create(string agent, object message);
        void Invoke(Guid receiver, object message);
        void Remind(DateTime dueDateTime, object message);
        void Publish(object message);
        void Subscribe(string target);
    }
}