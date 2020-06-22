using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.Blocks.Command;

namespace Apocryph.Core.Consensus
{
    public class Executor
    {
        private readonly string _chain;
        private readonly Func<string, (byte[]?, (string, byte[]), Guid?), Task<(byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>)>> _callWorker;

        public Executor(byte[] chain, Func<string, (byte[]?, (string, byte[]), Guid?), Task<(byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>)>> callWorker)
        {
            _chain = Convert.ToBase64String(chain); // TODO: Remove conversion
            _callWorker = callWorker;
        }

        public async Task<(IDictionary<string, byte[]>, object[], IDictionary<Guid, (string, string[])>)> Execute(
            IDictionary<string, byte[]> states, object[] commands, IDictionary<Guid, (string, string[])> capabilities)
        {
            var capabilityValidator = new CapabilityValidator(capabilities);

            var newCommands = new List<object>();

            foreach (var command in commands)
            {
                var targetReference = GetTargetReference(command);

                var targetState = capabilities[targetReference].Item1;

                var (state, actions, createdReferences, attachedReferences) = command switch
                {
                    Invoke cmd => await _callWorker(targetState, (states[targetState], cmd.Message, cmd.Reference)),
                    _ => throw new ArgumentException()
                };

                if (state != null)
                {
                    capabilityValidator.RegisterCarrier(state);
                    states[targetState] = state;
                }
                else
                {
                    states.Remove(targetState);
                }

                var newCapabilities = createdReferences.ToDictionary(
                    @ref => @ref.Key,
                    @ref => (_chain, @ref.Value));
                capabilityValidator.AddCapabilities(newCapabilities);

                if (actions != null)
                {
                    foreach (var (name, @params) in actions)
                    {
                        object? newCommand = name switch
                        {
                            nameof(Invoke) => capabilityValidator.ValidateMessageAndRegisterAsCarrier((Guid)@params[0], ((string, byte[]))@params[1])
                                ? new Invoke((Guid)@params[0], ((string, byte[]))@params[1])
                                : null,
                            nameof(Remind) => new Remind((DateTime)@params[0], ((string, byte[]))@params[1]),
                            nameof(Publish) => new Publish(((string, byte[]))@params[0]),
                            nameof(Subscribe) => new Subscribe((string)@params[0]),
                            _ => throw new ArgumentException()
                        };
                        if (newCommand != null) newCommands.Add(newCommand);
                    }
                }

                if (attachedReferences != null)
                {
                    foreach (var (reference, carrier) in attachedReferences)
                    {
                        if (capabilityValidator.ValidateAttachedReference((reference, carrier)))
                        {
                            newCapabilities[reference] = capabilities[reference];
                        }
                    }
                }

                capabilities = newCapabilities;
            }

            return (states, newCommands.ToArray(), capabilities);
        }

        public bool FilterCommand(object command, IDictionary<Guid, (string, string[])> capabilities)
        {
            var targetReference = GetTargetReference(command);
            return capabilities.ContainsKey(targetReference);
        }

        private Guid GetTargetReference(object command)
        {
            return command switch
            {
                Invoke cmd => cmd.Reference,
                _ => throw new ArgumentException()
            };
        }
    }
}