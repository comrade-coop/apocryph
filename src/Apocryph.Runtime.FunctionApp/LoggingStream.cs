using System;
using System.Collections.Generic;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.Blocks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class LoggingStream
    {
        [FunctionName(nameof(LoggingStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("filter")] IAsyncEnumerable<Block> filter,
            CancellationToken cancellationToken)
        {
            await foreach (var block in filter)
            {
                Console.WriteLine("New block on {0} by {1} ({2})", block.ChainId, block.Proposer, block.ProposerAccount);

                Console.WriteLine("Inputs:");
                foreach (var inputCommand in block.InputCommands)
                {
                    Console.WriteLine("\t{0}", inputCommand);
                }

                Console.WriteLine("Outputs:");
                foreach (var command in block.Commands)
                {
                    Console.WriteLine("\t{0}", command);
                }

                Console.WriteLine("States:");
                foreach (var (stateName, state) in block.States)
                {
                    Console.WriteLine("\t{0}: {1}", stateName, Encoding.UTF8.GetString(state));
                }

                Console.WriteLine("Capabilities:");
                foreach (var (capabilityId, (stateName, methods)) in block.Capabilities)
                {
                    Console.WriteLine("\t{0}: {1} -> {2}", capabilityId, stateName, string.Join(", ", methods));
                }
            }
        }
    }
}