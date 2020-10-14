using System;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Blocks.Command;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using TestHarness.FunctionApp.Mock;


namespace TestHarness.FunctionApp
{
    public static class Launcher
    {
        [FunctionName("Launcher")]
        public static async Task RunAsync([PerperModuleTrigger(RunOnStartup = true)]
            PerperModuleContext context,
            CancellationToken cancellationToken)
        {
            var slotCount = 10; // 30

            var pingChainId = Guid.NewGuid();
            var pongChainId = Guid.NewGuid();

            // For test harness we create seed blocks with valid references in the states
            var pingReference = Guid.NewGuid();
            var pongReference = Guid.NewGuid();

            var chains = new Dictionary<Guid, Chain>
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
            };

            var self = new Peer(new byte[0]);
            var slotGossipsStream = "DummyStream";
            var hashRegistryStream = typeof(HashRegistryStream).FullName! + ".Run";
            var hashRegistryWorker = typeof(HashRegistryWorker).FullName! + ".Run";
            var outsideGossipsStream = "DummyStream";
            var outsideQueriesStream = "DummyStream";

            var mode = Environment.GetEnvironmentVariable("ApocryphEnvironment");
            if (mode == "ipfs")
            {
                self = await context.CallWorkerAsync<Peer>("Apocryph.Runtime.FunctionApp.SelfPeerWorker.Run", new { }, default);
                slotGossipsStream = "Apocryph.Runtime.FunctionApp.IpfsSlotGossipStream.Run";
                hashRegistryStream = "Apocryph.Runtime.FunctionApp.HashRegistryStream.Run";
                hashRegistryWorker = "Apocryph.Runtime.FunctionApp.HashRegistryWorker.Run";
                outsideGossipsStream = "Apocryph.Runtime.FunctionApp.IpfsGossipStream.Run";
                outsideQueriesStream = "Apocryph.Runtime.FunctionApp.IpfsQueryStream.Run";
            }

            await context.StreamActionAsync("Apocryph.Runtime.FunctionApp.ChainListStream.Run", new
            {
                self,
                hashRegistryStream,
                hashRegistryWorker,
                outsideGossipsStream,
                outsideQueriesStream,
                slotGossipsStream,
                chains
            });

            await context.BindOutput(cancellationToken);
        }
    }
}