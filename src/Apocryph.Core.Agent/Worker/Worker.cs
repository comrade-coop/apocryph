using System;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading.Tasks;

namespace Apocryph.Core.Agent.Worker
{
    public class Worker<T> where T : class
    {
        private readonly IAgent<T> _agent;
        private readonly Context<T> _context;

        public Worker(IAgent<T> agent)
        {
            _agent = agent;
            _context = new Context<T>();

            _agent.Setup(_context);
        }

        public async Task<(byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>)> Run((byte[]?, (string, byte[]), Guid?) input)
        {
            var (state, (messageType, messagePayload), reference) = input;
            _context.Load(state);
            var message = JsonSerializer.Deserialize(messagePayload, Type.GetType(messageType));
            await _agent.Run(_context, message, reference);
            return _context.Save();
        }
    }
}