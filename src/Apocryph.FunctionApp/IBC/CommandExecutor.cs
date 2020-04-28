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

namespace Apocryph.FunctionApp.IBC
{
    public static class CommandExecutor
    {
        [FunctionName(nameof(CommandExecutor))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("commandsStream")] object[] commandsStream,
            [Perper("otherValidatorSetsStream")] object[] otherValidatorSetsStream,
            [Perper("notificationsStream")] object[] notificationsStream,
            [Perper("agentId")] string agentId,
            [Perper("services")] string[] services,
            [Perper("genesisMessage")] object genesisMessage,
            [Perper("ipfsGateway")] string ipfsGateway,
            CancellationToken cancellationToken)
        {

            await using var genesisMessageStream = await context.StreamFunctionAsync(nameof(TestDataGenerator), new
            {
                delay = TimeSpan.FromSeconds(10),
                data = genesisMessage
            });

            await using var subscriptionCommandExecutorStream = await context.StreamFunctionAsync(nameof(SubscriptionCommandExecutor), new
            {
                ipfsGateway,
                otherValidatorSetsStream,
                commandsStream
            });

            await using var reminderCommandExecutorStream = await context.StreamFunctionAsync(nameof(ReminderCommandExecutor), new
            {
                commandsStream
            });

            await using var sendMessageCommandExecutorStream = await context.StreamActionAsync(nameof(SendMessageCommandExecutor), new
            {
                agentId,
                ipfsGateway,
                otherValidatorSetsStream,
                commandsStream,
            });

            await using var callNotificationProcessorStream = await context.StreamFunctionAsync(nameof(CallNotificationProcessor), new
            {
                agentId,
                ipfsGateway,
                otherValidatorSetsStream,
                notificationsStream,
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
                subscriptionCommandExecutorStream,
                callNotificationProcessorStream
            }.Concat(serviceStreams).ToArray();


            await context.BindOutput(outputStream, cancellationToken);
        }
    }
}