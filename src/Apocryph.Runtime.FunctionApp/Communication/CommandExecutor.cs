using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Agent;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Apocryph.Runtime.FunctionApp.Utils;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Communication
{
    public static class CommandExecutor
    {
        [FunctionName(nameof(CommandExecutor))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("commandsStream")] object[] commandsStream,
            [Perper("otherValidatorSetsStream")] object[] otherValidatorSetsStream,
            [Perper("notificationsStream")] object[] notificationsStream,
            [Perper("agentId")] string agentId,
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

            var outputStream = new []
            {
                reminderCommandExecutorStream,
                genesisMessageStream,
                subscriptionCommandExecutorStream,
                callNotificationProcessorStream
            };


            await context.BindOutput(outputStream, cancellationToken);
        }
    }
}