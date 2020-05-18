using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agent.Command;
using Apocryph.Agent.Core;
using Apocryph.Agent.Worker;
using Apocryph.Runtime.FunctionApp.Consensus.Core;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Consensus
{
    public class Validator
    {
        [FunctionName(nameof(Validator))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("node")] Node node,
            [Perper("lastBlock")] Block lastBlock,
            [PerperStream("queries")] IAsyncEnumerable<Query<Block>> queries,
            [PerperStream("output")] IAsyncCollector<Message<Block>> output,
            CancellationToken cancellationToken)
        {
            await foreach (var query in queries.WithCancellation(cancellationToken))
            {
                if (query.Receiver != node) continue;

                var block = query.Value;
                var valid = await Validate(context, node, lastBlock, block);
                await output.AddAsync(new Message<Block>(block, valid ? MessageType.Valid : MessageType.Invalid), cancellationToken);
            }
        }

        private async Task<bool> Validate(PerperStreamContext context, Node node, Block lastBlock, Block block)
        {
            var executor = new Executor(node?.ToString()!,
                async input => await context.CallWorkerAsync<WorkerOutput>("AgentWorker", new {input}, default));
            var command = lastBlock.Commands.FirstOrDefault(o => o is Invoke || o is Publish || o is Remind);
            var (newState, newCommands, newCapabilities) = await executor.Execute(
                lastBlock.State, command, lastBlock.Capabilities);
            return block.State.SequenceEqual(newState)
                   && block.Commands.SequenceEqual(newCommands)
                   && block.Capabilities.Count == newCapabilities.Count
                   && !block.Capabilities.Except(newCapabilities).Any();
        }
    }
}