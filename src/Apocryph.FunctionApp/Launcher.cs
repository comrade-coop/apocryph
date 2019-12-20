using System;
using System.Collections.Generic;
using System.Linq;
using System.Linq.Expressions;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Apocryph.FunctionApp.AgentZero.Publications;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class Launcher
    {
        [FunctionName("Launcher")]
        public static async Task Run([PerperStreamTrigger("Launcher")] IPerperStreamContext context,
            [Perper("cancellationToken")] CancellationToken cancellationToken)
        {
            ECParameters privateKey;
            ValidatorKey self;

            using (var dsa = ECDsa.Create())
            {
                privateKey = dsa.ExportParameters(true);
                self = new ValidatorKey{Key = dsa.ExportParameters(false)};
            }

            var ipfsGateway = "127.0.0.1:5001";

            await using var validatorSetPublicationsStream = await context.StreamFunctionAsync("UNIMPLEMENTED-AgentPublications", new
            {
                ipfsGateway,
                agentId = "0",
                filter = (Expression<Func<object, bool>>)(x => x is ValidatorSetPublication)
            });

            await using var aggregatedValidatorSetsStream = await context.StreamFunctionAsync("UNIMPLEMENTED-AggregateValidatorSets", new
            {
                validatorSetPublicationsStream
            });

            await using var validatorSetsStream = await context.StreamFunctionAsync("UNIMPLEMENTED-FilterValidatorSets", new
            {
                aggregatedValidatorSetsStream,
                self
            });

            await context.StreamActionAsync("ValidatorScheduler", new
            {
                validatorSetsStream,
                ipfsGateway,
                privateKey,
                self
            });
        }
    }
}