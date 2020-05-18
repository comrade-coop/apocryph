using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Apocryph.Agent.Command;
using Apocryph.Agent.Worker;

namespace Apocryph.Agent.Core
{
    public class Executor
    {
        private readonly string _agent;
        private readonly Func<WorkerInput, Task<WorkerOutput>> _callWorker;

        public Executor(string agent, Func<WorkerInput, Task<WorkerOutput>> callWorker)
        {
            _agent = agent;
            _callWorker = callWorker;
        }

        public async Task<(byte[]?, object[], IDictionary<Guid, (string, string[])>)> Execute(
            byte[]? state, object command, IDictionary<Guid, (string, string[])> capabilities)
        {
            var capabilityValidator = new CapabilityValidator(capabilities);

            var result = command switch
            {
                Remind cmd => await ExecuteRemindCommandAsync(state, cmd),
                Invoke cmd => await ExecuteInvokeCommandAsync(state, cmd),
                Publish cmd => await ExecutePublishCommandAsync(state, cmd),
                _ => throw new ArgumentException()
            };

            var newState = result.State;
            if (newState != null) capabilityValidator.RegisterCarrier(newState);

            var newCapabilities = result.CreatedReferences.ToDictionary(
                @ref => @ref.Key,
                @ref => (_agent, @ref.Value));
            capabilityValidator.AddCapabilities(newCapabilities);

            var newCommands = new List<object>();
            if (result.Actions != null)
            {
                foreach (var (name, @params) in result.Actions)
                {
                    object? newCommand = name switch
                    {
                        nameof(Invoke) => capabilityValidator.ValidateMessageAndRegisterAsCarrier((Guid) @params[0], ((string, byte[])) @params[1])
                            ? new Invoke((Guid) @params[0], ((string, byte[])) @params[1])
                            : null,
                        nameof(Remind) => new Remind((DateTime) @params[0], ((string, byte[])) @params[1]),
                        nameof(Publish) => new Publish(((string, byte[])) @params[0]),
                        nameof(Subscribe) => new Subscribe((string) @params[0]),
                        _ => throw new ArgumentException()
                    };
                    if (newCommand != null) newCommands.Add(newCommand);
                }
            }

            if (result.AttachedReferences != null)
            {
                foreach (var (reference, carrier) in result.AttachedReferences)
                {
                    if (capabilityValidator.ValidateAttachedReference((reference, carrier)))
                    {
                        newCapabilities[reference] = capabilities[reference];
                    }
                }
            }

            return (newState, newCommands.ToArray(), newCapabilities);
        }

        private async Task<WorkerOutput> ExecuteRemindCommandAsync(byte[]? state, Remind command)
        {
            if (command.DueDateTime > DateTime.UtcNow)
            {
                await Task.Delay(command.DueDateTime.Subtract(DateTime.UtcNow));
            }

            var input = new WorkerInput(command.Message)
            {
                State = state
            };
            return await _callWorker(input);
        }

        private async Task<WorkerOutput> ExecuteInvokeCommandAsync(byte[]? state, Invoke command)
        {
            var input = new WorkerInput(command.Message)
            {
                State = state,
                Reference = command.Reference
            };
            return await _callWorker(input);
        }

        private async Task<WorkerOutput> ExecutePublishCommandAsync(byte[]? state, Publish command)
        {
            var input = new WorkerInput(command.Message)
            {
                State = state
            };
            return await _callWorker(input);
        }
    }
}