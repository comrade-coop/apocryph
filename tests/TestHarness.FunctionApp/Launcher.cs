using System;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Blocks.Command;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace TestHarness.FunctionApp
{
    public static class Launcher
    {
        [FunctionName("Launcher")]
        public static async Task RunAsync([PerperStreamTrigger(RunOnStartup = true)]
            PerperStreamContext context,
            CancellationToken cancellationToken)
        {
            await context.StreamActionAsync("ChainList", new
            {
                chains = new Dictionary<byte[], string>
                {
                    {new byte[0], typeof(ChainAgentPing).FullName!},
                    {new byte[0], typeof(ChainAgentPong).FullName!}
                }
            });

            // Exchange references between independent agents (publications?)

            // For test harness we create seed blocks with valid references in the states
            var pingReference = Guid.NewGuid();
            var pongReference = Guid.NewGuid();

            var pingBlock = new Block(
                Guid.NewGuid(),
                Guid.NewGuid(),
                new Dictionary<string, byte[]>
                {
                    {
                        typeof(ChainAgentPing).FullName!,
                        JsonSerializer.SerializeToUtf8Bytes(new ChainAgentState {OtherReference = pongReference})
                    }
                },
                new object[]{},
                new object[]{},
                new Dictionary<Guid, (string, string[])>
                {
                    {pongReference, (typeof(ChainAgentPong).FullName!, new[] {typeof(string).FullName!})}
                });
            var pongBlock = new Block(
                Guid.NewGuid(),
                Guid.NewGuid(),
                new Dictionary<string, byte[]>
                {
                    {
                        typeof(ChainAgentPong).FullName!,
                        JsonSerializer.SerializeToUtf8Bytes(new ChainAgentState {OtherReference = pingReference})
                    }
                },
                new object[]{},
                new object[]
                {
                    new Invoke(pingReference, (typeof(string).FullName!, JsonSerializer.SerializeToUtf8Bytes("Pong")))
                },
                new Dictionary<Guid, (string, string[])>
                {
                    {pingReference, (typeof(ChainAgentPing).FullName!, new[] {typeof(string).FullName!})}
                });

            await context.BindOutput(cancellationToken);
        }
    }
}