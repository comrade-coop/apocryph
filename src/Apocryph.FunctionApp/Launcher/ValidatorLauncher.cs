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
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class ValidatorLauncher
    {
        [FunctionName(nameof(ValidatorLauncher))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("services")] string[] services,
            [Perper("validatorSetsStream")] object[] validatorSetsStream,
            [Perper("otherValidatorSetsStream")] object[] otherValidatorSetsStream,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("privateKey")] ECParameters privateKey,
            [Perper("self")] ValidatorKey self,
            CancellationToken cancellationToken)
        {
            var genesisMessage = ("", (object)new InitMessage());

            await using var validatorSchedulerStream = await context.StreamActionAsync(nameof(PBFTFullNode), new
            {
                agentId,
                services,
                validatorSetsStream,
                otherValidatorSetsStream,
                genesisMessage,
                ipfsGateway,
                privateKey,
                self
            });

            await context.BindOutput(cancellationToken);
        }
    }
}