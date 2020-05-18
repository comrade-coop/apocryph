using System;
using System.Text.Json;
using System.Threading.Tasks;
using Apocryph.Agent.Protocol;

namespace Apocryph.Agent.Worker
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

        public async Task<WorkerOutput> Run(WorkerInput input)
        {
            _context.Load(input.State);
            var (messageType, messagePayload) = input.Message;
            var message = JsonSerializer.Deserialize(messagePayload, Type.GetType(messageType));
            await _agent.Run(_context, message, input.Reference);
            var (state, actions, createdReferences, attachedReferences) = _context.Save();
            return new WorkerOutput(state, actions, createdReferences, attachedReferences);
        }
    }
}