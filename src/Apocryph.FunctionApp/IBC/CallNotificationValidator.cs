using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Model;
using Apocryph.FunctionApp.Ipfs;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp.IBC
{
    public static class CallNotificationValidator
    {
        public class State
        {
            public Dictionary<string, IHashed<ValidatorSet>> ValidatorSets { get; set; }
        }

        [FunctionName(nameof(CallNotificationValidator))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("agentId")] string agentId,
            [PerperStream("otherValidatorSetsStream")] IAsyncEnumerable<Dictionary<string, IHashed<ValidatorSet>>> otherValidatorSetsStream,
            [PerperStream("notificationsStream")] IAsyncEnumerable<IHashed<CallNotification>> notificationsStream,
            [PerperStream("outputStream")] IAsyncCollector<IHashed<CallNotification>> outputStream,
            ILogger logger)
        {
            var cts = new CancellationTokenSource(); // Dispose!
            await using var outputStreams = new AsyncDisposableList();
            await using var utilityStreams = new AsyncDisposableList();

            var state = await context.FetchStateAsync<State>() ?? new State();

            await Task.WhenAll(
                otherValidatorSetsStream.ForEachAsync(async validatorSets =>
                {
                    state.ValidatorSets = validatorSets;
                    await context.UpdateStateAsync(state);
                }, CancellationToken.None),

                notificationsStream.ForEachAsync(async notification =>
                {
                    try
                    {
                        if (notification.Value.Command.Target != agentId)
                        {
                            return;
                        }

                        var notificationSignaturesValid = notification.Value.Commits.All(commit =>
                        {
                            if (commit.Value.For != notification.Value.Step)
                            {
                                return false;
                            }
                            var bytes = IpfsJsonSettings.ObjectToBytes<object>(commit.Value);
                            return commit.Signer.ValidateSignature(bytes, commit.Signature);
                        });

                        if (notificationSignaturesValid)
                        {
                            // TODO: Accept the previous validator set for a short period of time.
                            var validatorSet = state.ValidatorSets[notification.Value.From];
                            var signers = notification.Value.Commits.Select(commit => commit.Signer).Distinct();
                            if (validatorSet.Value.IsMoreThanTwoThirds(signers))
                            {
                                await outputStream.AddAsync(notification);
                            }
                        }
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None));
        }
    }
}