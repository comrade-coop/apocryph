using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Apocryph.HashRegistry;
using Apocryph.PerperUtilities;
using Apocryph.ServiceRegistry;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Consensus.FunctionApp
{
    public class Consensus
    {
        private IContext _context;
        private IState _state;

        public Consensus(IContext context, IState state)
        {
            _context = context;
            _state = state;
        }

        [FunctionName("Apocryph-Consensus")]
        public async Task Start([PerperTrigger] (IAgent serviceRegistry, HashRegistryProxy hashRegistry) input)
        {
            await _state.SetValue("serviceRegistry", input.serviceRegistry);
            await _state.SetValue("hashRegistry", input.hashRegistry);

            await input.serviceRegistry.CallActionAsync("RegisterHandler", ("Chain", new Handler(_context.Agent, "InstanceConsensus")));
        }

        [FunctionName("InstanceConsensus")]
        public async Task InstanceConsensus([PerperTrigger] ServiceLocator locator)
        {
            if (locator.Type != "Chain") throw new Exception("InstanceConsensus is implemented only for Chain serivces!");

            var chainId = Hash.FromString<Chain>(locator.Id);

            var serviceRegistry = await _state.GetValue<IAgent>("serviceRegistry", () => default!);
            var hashRegistry = await _state.GetValue<HashRegistryProxy>("hashRegistry", () => default!);

            var chain = await hashRegistry.RetrieveAsync(chainId);

            var (callsStream, callsStreamName) = await _context.CreateBlankStreamAsync<Message>();
            var (subscriptionsStream, subscriptionsStreamName) = await _context.CreateBlankStreamAsync<List<Reference>>();

            var routedInput = await _context.StreamFunctionAsync<Message>("RouterInput", (callsStream, subscriptionsStream, serviceRegistry));
            var consensusOutput = await _context.StartAgentAsync<IAsyncEnumerable<Message>>(chain.ConsensusType, (routedInput, subscriptionsStreamName, hashRegistry, chain));
            var routedOutput = await _context.StreamFunctionAsync<Message>("RouterOutput", (consensusOutput, chainId, serviceRegistry));

            var service = new Service(new Dictionary<string, string>() {
                {"messages", callsStreamName}
            }, new Dictionary<string, IStream>() {
                {"messages", (IStream)routedOutput}
            });

            await serviceRegistry.CallActionAsync("Register", (new ServiceLocator("Chain", chainId.ToString()), service));
        }
    }
}