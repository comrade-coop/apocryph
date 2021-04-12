using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Apocryph.Ipfs;
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
        public async Task Start([PerperTrigger] (IAgent serviceRegistry, int _) input)
        {
            await _state.SetValue("serviceRegistry", input.serviceRegistry);

            await input.serviceRegistry.CallActionAsync("RegisterHandler", ("Chain", new Handler(_context.Agent, "InstanceConsensus")));
        }

        [FunctionName("InstanceConsensus")]
        public async Task InstanceConsensus([PerperTrigger] ServiceLocator locator, IHashResolver hashResolver)
        {
            if (locator.Type != "Chain") throw new Exception("InstanceConsensus is implemented only for Chain serivces!");

            var chainId = Hash.FromString<Chain>(locator.Id);
            var chain = await hashResolver.RetrieveAsync(chainId);

            var serviceRegistry = await _state.GetValue<IAgent>("serviceRegistry", () => default!);

            var (callsStream, callsStreamName) = await _context.CreateBlankStreamAsync<Message>();
            var (subscriptionsStream, subscriptionsStreamName) = await _context.CreateBlankStreamAsync<List<Reference>>();

            var routedInput = await _context.StreamFunctionAsync<Message>("RouterInput", (callsStream, subscriptionsStream, serviceRegistry));
            var consensusOutput = await _context.StartAgentAsync<IAsyncEnumerable<Message>>(chain.ConsensusType, (routedInput, subscriptionsStreamName, chain));
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