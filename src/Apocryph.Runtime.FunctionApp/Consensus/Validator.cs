using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agent.Protocol;
using Apocryph.Runtime.FunctionApp.Consensus.Core;
using Apocryph.Runtime.FunctionApp.Execution;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Consensus
{
    public class Validator
    {
        private Block? _lastBlock;
        private HashSet<object>? _pendingCommands;
        private IAsyncCollector<Message<Block>>? _output;

        [FunctionName(nameof(Validator))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("lastBlock")] Block lastBlock,
            [Perper("pendingCommands")] HashSet<object> pendingCommands,
            [PerperStream("ibc")] IAsyncEnumerable<Block> ibc,
            [PerperStream("queries")] IAsyncEnumerable<Query<Block>> queries,
            [PerperStream("output")] IAsyncCollector<Message<Block>> output,
            CancellationToken cancellationToken)
        {
            _lastBlock = lastBlock;
            _pendingCommands = pendingCommands;
            _output = output;

            await Task.WhenAll(
                HandleIBC(ibc, node, cancellationToken),
                HandleQueries(context, queries, node, cancellationToken));
        }

        private async Task<bool> Validate(PerperStreamContext context, Node node, Block block)
        {
            var executor = new Executor(node?.ToString()!,
                async input => await context.CallWorkerAsync<WorkerOutput>("AgentWorker", new { input }, default));
            var (newStates, newCommands, newCapabilities) = await executor.Execute(
                _lastBlock!.States, block.InputCommands, _lastBlock!.Capabilities);
            return block.Equals(new Block(newStates, block.InputCommands, newCommands, newCapabilities));
            // Validate historical blocks as per protocol
        }

        private async Task HandleIBC(IAsyncEnumerable<Block> ibc, Node node, CancellationToken cancellationToken)
        {
            var executor = new Executor(node?.ToString()!, default!);

            await foreach (var block in ibc.WithCancellation(cancellationToken))
            {
                foreach (var command in block.Commands)
                {
                    if (executor.FilterCommand(command, _lastBlock!.Capabilities))
                    {
                        _pendingCommands!.Add(command);
                    }
                }
            }
        }

        private async Task HandleQueries(PerperStreamContext context, IAsyncEnumerable<Query<Block>> queries, Node node, CancellationToken cancellationToken)
        {
            await foreach (var query in queries.WithCancellation(cancellationToken))
            {
                if (query.Receiver != node) continue;

                var block = query.Value;
                var valid = await Validate(context, node, block);
                await _output!.AddAsync(new Message<Block>(block, valid ? MessageType.Valid : MessageType.Invalid), cancellationToken);
            }
        }
    }
}