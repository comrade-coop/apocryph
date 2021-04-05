using System;
using System.Threading.Tasks;
using Apocryph.Consensus;
using Apocryph.HashRegistry;
using Apocryph.HashRegistry.MerkleTree;
using Apocryph.PerperUtilities;
using Perper.WebJobs.Extensions.Fake;

namespace Apocryph.Executor.Test
{
    using Executor = Apocryph.Executor.FunctionApp.Executor;

    public static class ExecutorFakes
    {
        public class FakeHandler<TIn, TOut> : IHandler<TOut>
        {
            public Func<TIn, Task<TOut>> Delegate { get; }

            public FakeHandler(Func<TIn, Task<TOut>> @delegate)
            {
                Delegate = @delegate;
            }

            public Task<TOut> InvokeAsync(object? parameters)
            {
                return Delegate((TIn)parameters!);
            }
        }

        public static async Task<IHandler<(AgentState, Message[])>> GetExecutor(params (Hash<string>, Func<(AgentState, Message), Task<(AgentState, Message[])>>)[] handlers)
        {
            var executor = new Executor(new FakeState());

            foreach (var (hash, handler) in handlers)
            {
                await executor._Register((hash, new FakeHandler<(AgentState, Message), (AgentState, Message[])>(handler)));
            }

            return new FakeHandler<(AgentState, Message), (AgentState, Message[])>(executor.Execute);
        }

        public static (Hash<string>, Func<(AgentState, Message), Task<(AgentState, Message[])>>)[] TestAgents = new (Hash<string>, Func<(AgentState, Message), Task<(AgentState, Message[])>>)[]
        {
            (Hash.From("AgentInc"), ((AgentState state, Message message) input) =>
            {
                var target = input.state.Data.Deserialize<Reference>();
                var result = input.message.Data.Deserialize<int>() + 1;
                return Task.FromResult((input.state, new[] { new Message(target, ReferenceData.From(result)) }));
            }),

            (Hash.From("AgentDec"), ((AgentState state, Message message) input) =>
            {
                var target = input.state.Data.Deserialize<Reference>();
                var result = input.message.Data.Deserialize<int>() - 1;
                return Task.FromResult((input.state, new[] { new Message(target, ReferenceData.From(result)) }));
            })
        };


        public static async Task<(Chain chain, Message[] input, Message[] output)> GetTestAgentScenario(HashRegistryProxy hashRegistry, string consensusType, object? consensusParameters, int slotsCount)
        {
            var messageFilter = new string[] { typeof(int).FullName! };
            var fakeChainId = Hash.From("123").Cast<Chain>();

            var agentStates = new[] {
                new AgentState(0, ReferenceData.From(new Reference(fakeChainId, 0, messageFilter)), Hash.From("AgentInc")),
                new AgentState(1, ReferenceData.From(new Reference(fakeChainId, 1, messageFilter)), Hash.From("AgentDec"))
            };

            var agentStatesTree = await MerkleTreeBuilder.CreateRootFromValues(hashRegistry, agentStates, 2);

            var chain = new Chain(new ChainState(agentStatesTree, agentStates.Length), consensusType, consensusParameters, slotsCount);

            var chainId = await hashRegistry.StoreAsync(chain);

            var inputMessages = new Message[]
            {
                new Message(new Reference(chainId, 0, messageFilter), ReferenceData.From(4)),
                new Message(new Reference(chainId, 1, messageFilter), ReferenceData.From(3)),
            };

            var expectedOutputMessages = new Message[]
            {
                new Message(new Reference(fakeChainId, 0, messageFilter), ReferenceData.From(5)),
                new Message(new Reference(fakeChainId, 1, messageFilter), ReferenceData.From(2)),
            };

            return (chain, inputMessages, expectedOutputMessages);

        }
    }
}