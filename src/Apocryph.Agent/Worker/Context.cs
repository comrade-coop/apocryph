using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json;
using Apocryph.Agent.Api;

namespace Apocryph.Agent.Worker
{
    public class Context<T> : IContext<T> where T:class
    {
        public T? State { get; set; }

        private readonly Dictionary<Guid, string[]> _createdReferences;
        private readonly Dictionary<Guid, string> _attachedReferences;
        private readonly List<(string, object[])> _actions;

        public Context()
        {
            _createdReferences = new Dictionary<Guid, string[]>();
            _attachedReferences = new Dictionary<string, (object, string)>();
            _actions = new List<(string, object[])>();
        }

        public void RegisterInstance<TInterface, TClass>()
        {
            throw new NotImplementedException();
        }

        public TInterface CreateInstance<TInterface>()
        {
            throw new NotImplementedException();
        }

        public Guid CreateReference(Type[] messageTypes)
        {
            var result = Guid.NewGuid();
            _createdReferences[result] = messageTypes.Select(t => t.FullName!).ToArray();
            return result;
        }

        public void Create(string agent, object message)
        {
            _actions.Add((nameof(Create), new[] {agent, message}));
        }

        public void Invoke(Guid receiver, object message)
        {
            _actions.Add((nameof(Invoke), new[] {receiver, message}));
        }

        public void Remind(DateTime dueDateTime, object message)
        {
            _actions.Add((nameof(Remind), new[] {dueDateTime, message}));
        }

        public void Publish(object message)
        {
            _actions.Add((nameof(Publish), new[] {message}));
        }

        public void Subscribe(string target)
        {
            _actions.Add((nameof(Subscribe), new object[] {target}));
        }

        public void Load(byte[]? state)
        {
            State = state is null ? null : JsonSerializer.Deserialize<T>(state!);
        }

        public (byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>) Save()
        {
            var state = State is null ? null : JsonSerializer.SerializeToUtf8Bytes(State);
            return (state, _actions.ToArray(), _createdReferences, _attachedReferences);
        }
    }
}