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
        public static async Task RunAsync([PerperModuleTrigger(RunOnStartup = true)]
            PerperModuleContext context,
            CancellationToken cancellationToken)
        {
            var slotCount = 20; // 30

            var pingChainId = Guid.NewGuid();
            var pongChainId = Guid.NewGuid();

            // For test harness we create seed blocks with valid references in the states
            var pingReference = Guid.NewGuid();
            var pongReference = Guid.NewGuid();

            var slotGossips = await context.StreamFunctionAsync("DummyStream", new { });

            await context.StreamActionAsync("Apocryph.Runtime.FunctionApp.ChainListStream.Run", new
            {
                slotGossips,
                chains = new Dictionary<Guid, Chain>
                {
                    {pingChainId, new Chain(slotCount, new Block(
                        new Hash(new byte[] {}),
                        pingChainId,
                        null,
                        Guid.NewGuid(),
                        new Dictionary<string, byte[]>
                        {
                            {
                                typeof(ChainAgentPing).FullName! + ".Run",
                                JsonSerializer.SerializeToUtf8Bytes(new ChainAgentState {OtherReference = pongReference})
                            },
                            {
                                typeof(ChainAgentPong).FullName! + ".Run",
                                JsonSerializer.SerializeToUtf8Bytes(new ChainAgentState {OtherReference = pingReference})
                            }
                        },
                        new ICommand[] { },
                        new ICommand[]
                        {
                            new Invoke(pingReference, (typeof(string).FullName!, JsonSerializer.SerializeToUtf8Bytes("Init")))
                        },
                        new Dictionary<Guid, (string, string[])>
                        {
                            {pongReference, (typeof(ChainAgentPong).FullName! + ".Run", new[] {typeof(string).FullName!})},
                            {pingReference, (typeof(ChainAgentPing).FullName! + ".Run", new[] {typeof(string).FullName!})}
                        }))}
                }
            });

            await context.BindOutput(cancellationToken);
        }
    }
}