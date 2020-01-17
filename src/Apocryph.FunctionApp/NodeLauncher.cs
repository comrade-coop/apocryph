using System;
using System.Collections.Generic;
using System.Linq;
using System.Linq.Expressions;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.AgentZero.Publications;
using Apocryph.FunctionApp.Ipfs;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class NodeLauncher
    {
        [FunctionName(nameof(NodeLauncher))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            CancellationToken cancellationToken)
        {
            ECParameters privateKey;
            ValidatorKey self;

            using (var dsa = ECDsa.Create())
            {
                privateKey = dsa.ExportParameters(true);
                self = new ValidatorKey{Key = dsa.ExportParameters(false)};
            }

            var ipfsGateway = "http://127.0.0.1:5001";

            await using var agentZeroStream = await context.StreamFunctionAsync(nameof(IpfsInput), new
            {
                ipfsGateway,
                topic = "apocryph-agent-0"
            });

            await using var _inputVerifierStream = await context.StreamFunctionAsync(nameof(StepSignatureVerifier), new
            {
                stepsStream = agentZeroStream,
            });

            await using var inputVerifierStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _inputVerifierStream
            });

            await using var validatorSetsStream = await context.StreamFunctionAsync(nameof(ValidatorSets), new
            {
                inputVerifierStream
            });

            await using var validatorSchedulerStream = await context.StreamActionAsync(nameof(ValidatorScheduler), new
            {
                validatorSetsStream,
                ipfsGateway,
                privateKey,
                self
            });

            await context.BindOutput(cancellationToken);
        }
    }
}