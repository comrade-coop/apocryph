using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp.Service
{
    public static class IpfsInputSubscriber
    {
        private class State
        {
            public ISet<ValidatorKey> SubscribedKeys { get; set; } = new HashSet<ValidatorKey>();
        }

        [FunctionName(nameof(IpfsInputSubscriber))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("commandsStream")] IAsyncEnumerable<object> commandsStream,
            [PerperStream("ipfsStream")] IAsyncEnumerable<ISigned<object>> ipfsStream,
            [PerperStream("outputStream")] IAsyncCollector<(string, object)> outputStream,
            ILogger logger)
        {
            var state = await context.FetchStateAsync<State>() ?? new State();
            await Task.WhenAll(
                commandsStream.ForEachAsync(async command =>
                {
                    try
                    {
                        if (command is ValidatorKey key)
                        {
                            state.SubscribedKeys.Add(key);
                        }
                        await context.UpdateStateAsync(state);
                    }
                    catch (Exception e)
                    {
                        logger.LogError(e.ToString());
                    }
                }, CancellationToken.None),

                ipfsStream.ForEachAsync(async signed =>
                {
                    try
                    {
                        if (state.SubscribedKeys.Contains(signed.Signer))
                        {
                            await outputStream.AddAsync(("IpfsInput", signed.Value));
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