using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agent;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Testbed
{
    [Obsolete]
    public class Testbed
    {
        private readonly ILogger<Testbed> _logger;

        public Testbed(ILogger<Testbed> logger)
        {
            _logger = logger;
        }

        public async Task Setup(PerperStreamContext context, string agentDelegate, string runtimeDelegate, string monitorDelegate,
            CancellationToken cancellationToken)
        {
            var runtime = context.DeclareStream(runtimeDelegate);
            await context.StreamFunctionAsync(runtime, new { agentDelegate, commands = runtime });
            await context.StreamActionAsync(monitorDelegate, new { commands = runtime });
            await context.BindOutput(cancellationToken);
        }

        public async Task Agent(Func<object, AgentCapability, object, Task<AgentContext>> entryPoint,
            string agentId,
            object initMessage,
            IAsyncEnumerable<AgentCommands> commands, IAsyncCollector<AgentCommands> output,
            CancellationToken cancellationToken)
        {
            var states = new List<object>();
            await Task.WhenAll(
                InitAgent(entryPoint, states, agentId, initMessage, output),
                ExecuteAgent(entryPoint, states, agentId, commands, output, cancellationToken));
        }

        public async Task Runtime(PerperStreamContext context, string agentDelegate,
            IAsyncEnumerable<AgentCommands> commands, CancellationToken cancellationToken)
        {
            var agents = new List<IAsyncDisposable>();
            await Task.WhenAll(
                InitRuntime(context, agentDelegate, agents),
                ExecuteRuntime(context, commands, agents, cancellationToken));
        }

        public async Task Monitor(IAsyncEnumerable<AgentCommands> commands,
            CancellationToken cancellationToken)
        {
            await foreach (var commandsBatch in commands.WithCancellation(cancellationToken))
            {
                foreach (var command in commandsBatch.Commands)
                {
                    _logger.LogInformation($"{command.CommandType.ToString()} command with {command.Receiver?.Issuer} receiver");
                }
            }
        }

        private async Task InitRuntime(PerperStreamContext context, string agentDelegate,
            ICollection<IAsyncDisposable> agents)
        {
            await Task.Delay(TimeSpan.FromSeconds(1)); //Wait for Execute to engage Runtime

            var agent = await context.StreamFunctionAsync(agentDelegate, new
            {
                agentId = "AgentRoot",
                initMessage = new AgentRootInitMessage(),
                commands = context.GetStream()
            });
            agents.Add(agent);
            await context.RebindOutput(agents);
        }

        private async Task ExecuteRuntime(PerperStreamContext context, IAsyncEnumerable<AgentCommands> commands,
            ICollection<IAsyncDisposable> agents, CancellationToken cancellationToken)
        {
            await foreach (var commandsBatch in commands.WithCancellation(cancellationToken))
            {
                foreach (var command in commandsBatch.Commands)
                {
                    if (command.CommandType == AgentCommandType.CreateAgent)
                    {
                        var agent = await context.StreamFunctionAsync(command.Agent, new
                        {
                            agentId = command.AgentId,
                            initMessage = command.Message,
                            commands = context.GetStream()
                        });
                        agents.Add(agent);
                    }
                }
                await context.RebindOutput(agents);
            }
        }

        private async Task InitAgent(Func<object, AgentCapability, object, Task<AgentContext>> entryPoint,
            ICollection<object> states,
            string agentId,
            object initMessage, IAsyncCollector<AgentCommands> output)
        {
            await Task.Delay(TimeSpan.FromSeconds(1)); //Wait for Execute to engage Runtime

            var agentContext = await entryPoint(null,
                new AgentCapability
                {
                    Issuer = agentId,
                    MessageTypes = new[] { initMessage.GetType() }
                }, initMessage);
            states.Add(agentContext.InternalState);
            await output.AddAsync(agentContext.GetCommands());
        }

        private async Task ExecuteAgent(Func<object, AgentCapability, object, Task<AgentContext>> entryPoint,
            ICollection<object> states,
            string agentId,
            IAsyncEnumerable<AgentCommands> commands, IAsyncCollector<AgentCommands> output,
            CancellationToken cancellationToken)
        {
            var publishers = new HashSet<string>();
            await foreach (var commandsBatch in commands.WithCancellation(cancellationToken))
            {
                foreach (var command in commandsBatch.Commands)
                {
                    if (command.CommandType == AgentCommandType.SendMessage && command.Receiver.Issuer == agentId)
                    {
                        var agentContext = await entryPoint(states.Last(), command.Receiver, command.Message);
                        states.Add(agentContext.InternalState);
                        await output.AddAsync(agentContext.GetCommands(), cancellationToken);
                    }
                    else if (command.CommandType == AgentCommandType.Reminder && commandsBatch.Origin == agentId)
                    {
                        await Task.Delay(command.Timeout, cancellationToken);
                        var agentContext = await entryPoint(states.Last(), command.Receiver, command.Message);
                        states.Add(agentContext.InternalState);
                        await output.AddAsync(agentContext.GetCommands(), cancellationToken);
                    }
                    else if (command.CommandType == AgentCommandType.Subscribe && commandsBatch.Origin == agentId)
                    {
                        publishers.Add(command.AgentId);
                    }
                    else if (command.CommandType == AgentCommandType.Publish && publishers.Contains(commandsBatch.Origin))
                    {
                        var agentContext = await entryPoint(states.Last(), new AgentCapability { Issuer = agentId }, command.Message);
                        states.Add(agentContext.InternalState);
                        await output.AddAsync(agentContext.GetCommands(), cancellationToken);
                    }
                }
            }
        }
    }
}