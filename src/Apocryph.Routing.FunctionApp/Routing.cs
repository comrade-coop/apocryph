using System.Collections.Generic;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Apocryph.KoTH;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Routing.FunctionApp
{
    public static class Routing
    {
        [FunctionName("Apocryph-Routing")]
        public static async Task Start([PerperTrigger] (IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates, IAgent executor) input, IState state)
        {
            await state.SetValue("input", input);
        }

        [FunctionName("GetChainInstance")]
        public static async Task<(string, IStream<Message>)> GetChainInstance([PerperTrigger] Hash<Chain> chainId, IContext context, IState state, IHashResolver hashResolver)
        {
            // NOTE: Can benefit from locking
            var key = $"{chainId}";
            var (currentCallsStreamName, currentRoutedOutput) = await state.GetValue<(string, IStream<Message>?)>(key, () => ("", null));
            if (currentCallsStreamName != "" && currentRoutedOutput != null)
            {
                return (currentCallsStreamName, currentRoutedOutput);
            }

            var chain = await hashResolver.RetrieveAsync(chainId);

            var (callsStream, callsStreamName) = await context.CreateBlankStreamAsync<Message>();
            var (publicationsStream, publicationsStreamName) = await context.CreateBlankStreamAsync<Message>();
            var (subscriptionsStream, subscriptionsStreamName) = await context.CreateBlankStreamAsync<List<Reference>>();

            var (kothStates, executor) = await state.GetValue<(IAsyncEnumerable<(Hash<Chain>, Slot?[])>, IAgent)>("input", () => default!);

            var routedInput = await context.StreamFunctionAsync<Message>("RouterInput", (callsStream, subscriptionsStream));

            var (_, consensusOutput) = await context.StartAgentAsync<IAsyncEnumerable<Message>>(chain.ConsensusType, (
                (IAsyncEnumerable<Message>)routedInput,
                subscriptionsStreamName,
                chain,
                kothStates,
                executor
            ));

            var task = context.StreamActionAsync("RouterOutput", (publicationsStreamName, consensusOutput, chainId));
            var _ = task.ContinueWith(x => System.Console.WriteLine(x.Exception), TaskContinuationOptions.OnlyOnFaulted); // DEBUG: FakeStream does not log errors

            var resultValue = (callsStreamName, publicationsStream);
            await state.SetValue(key, resultValue);
            return resultValue;
        }
    }
}