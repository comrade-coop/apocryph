using System;
using System.Collections.Generic;
using System.Linq;
using System.Linq.Expressions;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Entrypoint
{
    public static class NodeLauncher
    {
        [FunctionName(nameof(NodeLauncher))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("ipfsGateway")] string ipfsGateway,
            [Perper("privateKey")] ECParameters privateKey,
            [Perper("self")] ValidatorKey self,
            CancellationToken cancellationToken)
        {
            await using var ipfsStream = await context.StreamFunctionAsync(nameof(IpfsInput), new
            {
                ipfsGateway,
                topic = "apocryph-agent-0"
            });
            /*
            await using var _unverifiedStepsStream = await context.StreamFunctionAsync(nameof(StepHashCollector), new
            {
                inputStream = ipfsStream
            });

            await using var unverifiedStepsStream = await context.StreamFunctionAsync(nameof(IpfsRecursiveLoader), new
            {
                ipfsGateway,
                hashStream = _unverifiedStepsStream
            });

            await using var _stepValidatorSetSplitterStream = await context.StreamFunctionAsync(nameof(StepValidatorSetSplitter), new
            {
                stepsStream = unverifiedStepsStream,
            });

            await using var stepValidatorSetSplitterStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _stepValidatorSetSplitterStream
            });

            await using var stepSignatureVerifierStream = await context.StreamFunctionAsync(nameof(StepSignatureVerifier), new
            {
                stepsStream = unverifiedStepsStream,
                stepValidatorSetSplitterStream
            });

            await using var _verifiedStepsStream = await context.StreamFunctionAsync(nameof(StepVerifiedStepGetter), new
            {
                stepSignatureVerifierStream
            });

            await using var verifiedStepsStream = await context.StreamFunctionAsync(nameof(IpfsLoader), new
            {
                ipfsGateway,
                hashStream = _verifiedStepsStream
            });

            await using var validatorSetsStream = await context.StreamFunctionAsync(nameof(AgentZetoStepOrderVerifier), new
            {
                stepsStream = verifiedStepsStream,
                ipfsGateway
            });

            await using var validatorSchedulerStream = await context.StreamActionAsync(nameof(ValidatorScheduler), new
            {
                validatorSetsStream,
                ipfsGateway,
                privateKey,
                self
            });*/

            await context.BindOutput(cancellationToken);
        }
    }
}