using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.CompilerServices;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Ipfs;
using Apocryph.Ipfs.MerkleTree;
using Apocryph.KoTH;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Consensus.Snowball.FunctionApp
{
    public class SnowballConsensus
    {
        private IContext _context;
        private IState _state;

        public SnowballConsensus(IContext context, IState state)
        {
            _context = context;
            _state = state;
        }

        [FunctionName("Apocryph-SnowballConsensus")]
        public async Task<IAsyncEnumerable<Message>> Start([PerperTrigger] (
                IAsyncEnumerable<Message> messages,
                string subscriptionsStream,
                Chain chain,
                IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates,
                IAgent executor) input,
            IHashResolver hashResolver)
        {
            var parameters = (SnowballParameters)input.chain.ConsensusParameters!;
            var emptyMessagesTree = await MerkleTreeBuilder.CreateRootFromValues(hashResolver, new Message[] { }, 3);
            var genesisBlock = new Block(null, emptyMessagesTree, emptyMessagesTree, input.chain.GenesisState);
            var genesis = await hashResolver.StoreAsync(genesisBlock);
            var self = await hashResolver.StoreAsync(input.chain);

            var messagePoolTask = _context.StreamActionAsync("MessagePool", (input.messages, self));
            var kothTask = _context.StreamActionAsync("KothProcessor", (self, input.kothStates));

            return await _context.StreamFunctionAsync<Message>("SnowballStream", (self, genesis, parameters, input.executor));
        }

        [FunctionName("SnowballStream")]
        public async IAsyncEnumerable<Message> SnowballStream([PerperTrigger] (
                Hash<Chain> self,
                Hash<Block> genesis,
                SnowballParameters parameters,
                IAgent executor) input,
            IHashResolver hashResolver,
            IPeerConnector peerConnector,
            [EnumeratorCancellation] CancellationToken cancellationToken = default)
        {
            var queryPath = $"snowball/{input.self}";

            await peerConnector.ListenQuery<Query, Query>(queryPath, async (peer, request) =>
            {
                // NOTE: Should be using some locking here
                var snowballState = await _state.GetValue<SnowballState>($"snowballState-{request.Round}");
                snowballState.ProcessQuery(request.Value);
                var result = new Query(snowballState.CurrentValue!, request.Round);
                await _state.SetValue($"snowballState-{request.Round}", snowballState);

                return result;
            }, cancellationToken);

            async Task<(ChainState, IMerkleTree<Message>)> ExecuteBlock(ChainState chainState, IMerkleTree<Message> inputMessages)
            {
                var agentStates = await chainState.AgentStates.EnumerateItems(hashResolver).ToDictionaryAsync(x => x.Nonce, x => x);
                var outputMesages = new List<Message>();

                await foreach (var message in inputMessages.EnumerateItems(hashResolver))
                {
                    var state = agentStates[message.Target.AgentNonce];

                    var (newState, resultMessages) = await input.executor.CallFunctionAsync<(AgentState, Message[])>("Execute", (state, message));

                    agentStates[message.Target.AgentNonce] = state;
                    outputMesages.AddRange(resultMessages);
                }

                var outputStatesTree = await MerkleTreeBuilder.CreateRootFromValues(hashResolver, agentStates.Values, 3); // FIXME: what about ordering?
                var outputState = new ChainState(outputStatesTree, chainState.NextAgentNonce);
                var outputMessagesTree = await MerkleTreeBuilder.CreateRootFromValues(hashResolver, outputMesages, 3);

                return (outputState, outputMessagesTree);
            }

            async Task<bool> ValidateBlock(Block block, Hash<Block> expectedPrevious)
            {
                if (block.Previous != expectedPrevious)
                {
                    return false;
                }

                var inputMessagesSet = await _state.GetValue<List<Message>>("messagePool");
                await foreach (var inputMessage in block.InputMessages.EnumerateItems(hashResolver))
                {
                    if (!inputMessagesSet.Remove(inputMessage))
                        return false;
                }

                var previous = await hashResolver.RetrieveAsync(expectedPrevious);
                var (outputState, outputMessages) = await ExecuteBlock(previous.State, block.InputMessages);

                return Hash.From(block.State) == Hash.From(outputState) && Hash.From(block.OutputMessages) == Hash.From(outputMessages);
            }

            async Task<Block?> ProposeBlock(Hash<Block> previousHash)
            {
                var previous = await hashResolver.RetrieveAsync(previousHash);
                var inputMessagesList = (await _state.GetValue<List<Message>>("messagePool")).ToArray();

                if (inputMessagesList.Length == 0 && await _state.GetValue("finished", () => false))
                {
                    return null; // DEBUG: Used for testing purposes mainly
                }

                var inputMessages = await MerkleTreeBuilder.CreateRootFromValues(hashResolver, inputMessagesList, 3);

                var (outputStates, outputMessages) = await ExecuteBlock(previous.State, inputMessages);

                return new Block(previousHash, inputMessages, outputMessages, outputStates);
            }

            var currentRound = await _state.GetValue<int>("currentRound", () => 0);

            // TODO: Sync to current state first

            if (currentRound == 0)
            {
                var newBlock = await ProposeBlock(input.genesis);
                if (newBlock == null) yield break; // DEBUG: Used for testing purposes mainly
                var newBlockHash = await hashResolver.StoreAsync(newBlock);
                var snowball = await _state.GetValue<SnowballState>($"snowballState-{currentRound}");
                snowball.ProcessQuery(newBlockHash);
                await _state.SetValue($"snowballState-{currentRound}", snowball);
            }


            // NOTE: Might benefit from locking
            while (true)
            {
                IEnumerable<Task<Query>> replyTasks;
                {
                    var snowball = await _state.GetValue<SnowballState>($"snowballState-{currentRound}");
                    if (snowball.CurrentValue == null)
                    {
                        await Task.Delay(100);
                        continue;
                    }
                    var kothPeers = await _state.GetValue<Peer[]>("kothPeers", () => new Peer[] { });
                    if (kothPeers.Length == 0)
                    {
                        await Task.Delay(100);
                        continue;
                    }
                    var sampledPeers = snowball.SamplePeers(input.parameters, kothPeers);
                    var query = new Query(snowball.CurrentValue, currentRound);
                    replyTasks = sampledPeers.Select(peer => peerConnector.Query<Query, Query>(peer, queryPath, query));
                }

                var responses = (await Task.WhenAll(replyTasks)).Select(reply => reply.Value);

                Hash<Block>? finishedHash = null;
                {
                    var snowball = await _state.GetValue<SnowballState>($"snowballState-{currentRound}");
                    var finished = snowball.ProcessResponses(input.parameters, responses);

                    if (finished)
                    {
                        finishedHash = snowball.CurrentValue;
                    }

                    await _state.SetValue($"snowballState-{currentRound}", snowball);
                }

                if (finishedHash != null)
                {
                    var finishedBlock = await hashResolver.RetrieveAsync(finishedHash);

                    var previousHash = await _state.GetValue<Hash<Block>>("lastBlock", () => input.genesis);
                    if (await ValidateBlock(finishedBlock, previousHash))
                    {
                        previousHash = finishedHash;
                        await _state.SetValue("lastBlock", previousHash);

                        await foreach (var outputMessage in finishedBlock.OutputMessages.EnumerateItems(hashResolver))
                        {
                            yield return outputMessage;
                        }

                        var messagePool = await _state.GetValue<List<Message>>("messagePool");
                        await foreach (var processedMessage in finishedBlock.InputMessages.EnumerateItems(hashResolver))
                        {
                            messagePool.Remove(processedMessage);
                        }
                        await _state.SetValue("messagePool", messagePool);
                    }

                    await _state.SetValue("currentRound", ++currentRound);
                    // FIXME: Calculate proposers or proposal order from (previousHash, currentRound)
                    var newBlock = await ProposeBlock(previousHash);

                    if (newBlock == null) yield break; // DEBUG: Used for testing purposes mainly

                    var newBlockHash = await hashResolver.StoreAsync(newBlock);

                    var snowball = await _state.GetValue<SnowballState>($"snowballState-{currentRound}");
                    snowball.ProcessQuery(newBlockHash);
                    await _state.SetValue($"snowballState-{currentRound}", snowball);
                }
            }
        }

        [FunctionName("MessagePool")]
        public async Task MessagePool([PerperTrigger] (IAsyncEnumerable<Message> inputMessages, Hash<Chain> self) input)
        {
            await foreach (var message in input.inputMessages)
            {
                // FIXME: Handled by routing?
                if (message.Target.Chain != input.self || !message.Target.AllowedMessageTypes.Contains(message.Data.Type))
                    continue;

                var messagePool = await _state.GetValue<List<Message>>("messagePool");
                messagePool.Add(message);
                await _state.SetValue("messagePool", messagePool);
            }
            await _state.SetValue("finished", true); // DEBUG: Used for testing purposes mainly
        }

        [FunctionName("KothProcessor")]
        public async Task KothProcessor([PerperTrigger] (Hash<Chain> chain, IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates) input)
        {
            await foreach (var (chain, slots) in input.kothStates)
            {
                if (chain != input.chain)
                    continue;

                var peers = slots.Where(s => s != null).Select(s => s!.Peer).ToArray();

                await _state.SetValue("kothPeers", peers);
            }
        }
    }
}