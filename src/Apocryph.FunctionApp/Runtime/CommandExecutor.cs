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
    public static class CommandExecutor
    {
        [FunctionName(nameof(CommandExecutor))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [Perper("services")] string[] services,
            [Perper("commandsStream")] object[] commandsStream,
            [Perper("genesisMessage")] object genesisMessage,
            [Perper("ipfsGateway")] string ipfsGateway,
            CancellationToken cancellationToken)
        {

            await using var genesisMessageStream = await context.StreamFunctionAsync(nameof(TestDataGenerator), new
            {
                delay = TimeSpan.FromSeconds(10),
                data = genesisMessage
            });

            /* await using var agentZeroStream = await context.StreamFunctionAsync(nameof(IpfsInput), new
            {
                ipfsGateway,
                topic = 0 //"apocryph-agent-0"
            });

            await using var _inputVerifierStream = await context.StreamFunctionAsync(nameof(StepSignatureVerifier), new
            {
                validatorSetsStream, // TODO: Should give agent 0's validator set instead !!
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

            await using var subscriptionCommandExecutorStream = await context.StreamFunctionAsync(nameof(SubscriptionCommandExecutor), new
            {
                ipfsGateway,
                commandsStream,
                validatorSetsStream
            }); */

            await using var reminderCommandExecutorStream = await context.StreamFunctionAsync(nameof(ReminderCommandExecutor), new
            {
                commandsStream
            });

            await using var serviceFilterStreams = new AsyncDisposableList();
            await using var serviceStreams = new AsyncDisposableList();
            foreach (var serviceName in services) {
                var functionName = "Service" + serviceName;
                var filteredCommandsStream = await context.StreamFunctionAsync(nameof(ServiceCommandFilter), new
                {
                    commandsStream,
                    serviceName
                });
                serviceFilterStreams.Add(filteredCommandsStream);
                var serviceStream = await context.StreamFunctionAsync(functionName, new
                {
                    commandsStream = filteredCommandsStream,
                    agentId,
                    ipfsGateway
                });
                serviceStreams.Add(serviceStream);
            }

            var outputStream = new []
            {
                reminderCommandExecutorStream,
                genesisMessageStream,
                // subscriptionCommandExecutorStream
            }.Concat(serviceStreams).ToArray();


            await context.BindOutput(outputStream, cancellationToken);
        }
    }
}