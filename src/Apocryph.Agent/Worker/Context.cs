using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Security.Cryptography;
using System.Text.Json;

namespace Apocryph.Agent.Worker
{
    public class Context<T> : IContext<T> where T : class
    {
        public T? State { get; set; }

        private readonly Dictionary<Guid, string[]> _createdReferences;
        private readonly Dictionary<Guid, object> _attachedReferences;
        private readonly List<(string, object[])> _actions;
        private readonly Dictionary<Type, Type> _registrations;

        public Context()
        {
            _createdReferences = new Dictionary<Guid, string[]>();
            _attachedReferences = new Dictionary<Guid, object>();
            _actions = new List<(string, object[])>();
            _registrations = new Dictionary<Type, Type>();
        }

        public void RegisterInstance<TInterface, TClass>()
        {
            _registrations.Add(typeof(TInterface), typeof(TClass));
        }

        public TInterface CreateInstance<TInterface>(Action<TInterface> initializer = null)
        {
            var result = DispatchProxy.Create<TInterface, ReferenceProxy>();
            var classType = _registrations[typeof(TInterface)];
            (result as ReferenceProxy)?.Init(Activator.CreateInstance(classType)!, _attachedReferences);
            return result;
        }

        public Guid CreateReference(Type[] messageTypes)
        {
            var result = Guid.NewGuid();
            _createdReferences[result] = messageTypes.Select(t => t.FullName!).ToArray();
            return result;
        }

        public void Create(string agent, object message)
        {
            Invoke(Guid.Empty, message);
        }

        public void Invoke(Guid receiver, object message)
        {
            _actions.Add((nameof(Invoke), new object[]
            {
                receiver,
                (message.GetType().FullName!, JsonSerializer.SerializeToUtf8Bytes(message))
            }));
        }

        public void Remind(DateTime dueDateTime, object message)
        {
            _actions.Add((nameof(Remind), new object[]
            {
                dueDateTime,
                (message.GetType().FullName!, JsonSerializer.SerializeToUtf8Bytes(message))
            }));
        }

        public void Publish(object message)
        {
            _actions.Add((nameof(Publish), new object[]
            {
                (message.GetType().FullName!, JsonSerializer.SerializeToUtf8Bytes(message))
            }));
        }

        public void Subscribe(string target)
        {
            _actions.Add((nameof(Subscribe), new object[] { target }));
        }

        public void Load(byte[]? state)
        {
            State = state is null ? null : JsonSerializer.Deserialize<T>(state!);
        }

        public (byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>) Save()
        {
            var state = State is null ? null : JsonSerializer.SerializeToUtf8Bytes(State);
            var attachedReferences = _attachedReferences.ToDictionary(pair => pair.Key, pair =>
             {
                 var payload = JsonSerializer.SerializeToUtf8Bytes(pair.Value);
                 using var sha1 = new SHA1CryptoServiceProvider();
                 return Convert.ToBase64String(sha1.ComputeHash(payload));
             });
            return (state, _actions.ToArray(), _createdReferences, attachedReferences);
        }
    }
}