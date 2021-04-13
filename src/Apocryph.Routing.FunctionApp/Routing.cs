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
            var (subscriptionsStream, subscriptionsStreamName) = await context.CreateBlankStreamAsync<List<Reference>>();

            var (kothStates, executor) = await state.GetValue<(IAsyncEnumerable<(Hash<Chain>, Slot?[])>, IAgent)>("input", () => default!);

            var routedInput = await context.StreamFunctionAsync<Message>("RouterInput", (callsStream, subscriptionsStream));

            var consensusParameters = (
                (IAsyncEnumerable<Message>)routedInput,
                subscriptionsStreamName,
                chain,
                kothStates,
                executor
            );

            var (_, consensusOutput) = await context.StartAgentAsync<IAsyncEnumerable<Message>>(chain.ConsensusType, consensusParameters);
            var routedOutput = await context.StreamFunctionAsync<Message>("RouterOutput", (consensusOutput, chainId));

            var value = (callsStreamName, routedOutput);
            await state.SetValue(key, value);
            return value;
        }
    }
}