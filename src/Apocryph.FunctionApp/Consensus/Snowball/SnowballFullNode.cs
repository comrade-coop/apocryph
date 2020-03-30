using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Ipfs;
using Apocryph.FunctionApp.Model;
using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class SnowballConsensus
    {
        [FunctionName(nameof(SnowballConsensus))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("services")] string[] services,
            [Perper("validatorSetsStream")] object[] validatorSetsStream,
            [Perper("genesisMessage")] (string, object) genesisMessage,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("privateKey")] ECParameters privateKey,
            [Perper("self")] ValidatorKey self,
            CancellationToken cancellationToken)
        {
            var query = "/x/apocryph-agent/" + agentId;

            var ipfs = new IpfsClient(ipfsGateway);

            await ipfs.DoCommandAsync("p2p/listen", cancellationToken, query, new [] {"arg=/ip4/127.0.0.1/tcp/1234"});


            try
            {
                await context.BindOutput(cancellationToken);
            }
            finally
            {
                await ipfs.DoCommandAsync("p2p/close", cancellationToken, null, new [] {"target-address=/ip4/127.0.0.1/tcp/1234"});
            }
            // await using var ipfsStream = await context.StreamFunctionAsync(nameof(IpfsInput), new
            // {
            //     ipfsGateway,
            //     topic
            // });


            // await using var commandsStream = await context.StreamFunctionAsync(nameof(CommandSplitter), new
            // {
            //     stepsStream = stepOrderVerifierStream
            // });
            //
            // var commandExecutorStream =  await context.StreamFunctionAsync(nameof(CommandExecutor), new
            // {
            //     commandsStream,
            //     agentId,
            //     services,
            //     genesisMessage,
            //     ipfsGateway
            // });
            //
            //
            // await using var proposerRuntimeStream = await context.StreamFunctionAsync(nameof(Runtime), new
            // {
            //     self,
            //     agentId,
            //     inputStream = proposerFilterStream,
            // });
            //
            // // Output
            //
            // await using var signerStream = await context.StreamFunctionAsync(nameof(Signer), new
            // {
            //     self,
            //     privateKey,
            //     dataStream = new[] {proposerStream, validatorRuntimeVotingStream, inputValidatorStream, voteCounterStream}
            // });

            // await context.BindOutput(cancellationToken);
        }
    }
}