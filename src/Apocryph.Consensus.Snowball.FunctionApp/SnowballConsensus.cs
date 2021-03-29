using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Apocryph.HashRegistry;
using Apocryph.HashRegistry.MerkleTree;
using Apocryph.KoTH;
using Apocryph.Peering;
using Apocryph.PerperUtilities;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.Consensus.Snowball.FunctionApp
{
    public class SnowballConsensus
    {
        static public string PeeringProtocol = "snowball";
        private IContext _context;
        private IState _state;

        public SnowballConsensus(IContext context, IState state)
        {
            _context = context;
            _state = state;
        }

        [FunctionName("Apocryph-SnowballConsensus")]
        public async Task<IAsyncEnumerable<Message>> Start([PerperTrigger] (IAsyncEnumerable<Message> messages, string subscribtionsStream, HashRegistryProxy hashRegistry, Chain chain, IAgent peering, IAsyncEnumerable<(Hash<Chain>, Slot?[])> kothStates) input)
        {
            var parameters = (SnowballParameters)input.chain.ConsensusParameters!;
            var emptyMessagesTree = await MerkleTreeBuilder.CreateRootFromValues(input.hashRegistry, new Message[] { }, 3);
            var genesisBlock = new Block(null, emptyMessagesTree, emptyMessagesTree, input.chain.GenesisState);
            var genesis = await input.hashRegistry.StoreAsync(genesisBlock);
            var self = await input.hashRegistry.StoreAsync(input.chain);

            await input.peering.CallActionAsync("Register", (PeeringProtocol, _context.Agent.GetHandler<IAsyncEnumerable<object>>("PeeringResponder")));

            var messagePoolTask = _context.StreamActionAsync("MessagePool", (input.messages, self));
            var kothTask = _context.StreamActionAsync("KothProcessor", (self, input.kothStates));

            return await _context.StreamFunctionAsync<Message>("SnowballStream", (input.peering, input.hashRegistry, parameters, self, genesis));
        }

        [FunctionName("SnowballStream")]
        public async IAsyncEnumerable<Message> SnowballStream([PerperTrigger] (IAgent peering, HashRegistryProxy hashRegistry, SnowballParameters parameters, Hash<Chain> self, Hash<Block> genesis) input)
        {
            async Task<Query> SendQuery(Peer target, Query query)
            {
                var messages = new ListAsyncEnumerable<object>(new object[] { query });
                var result = await input.peering.CallFunctionAsync<IAsyncEnumerable<object>>("Connect", (target, PeeringProtocol, messages));

                return (Query)await result.FirstAsync();
            }

            async Task<(AgentState, Message[])> Execute(AgentState agentState, Message message)
            {
                var (_, result) = await _context.StartAgentAsync<(AgentState, Message[])>(agentState.Handler, (agentState, message));

                return result;
            }

            async Task<(ChainState, IMerkleTree<Message>)> ExecuteBlock(ChainState chainState, IMerkleTree<Message> inputMessages)
            {
                var agentStates = await chainState.AgentStates.EnumerateItems(input.hashRegistry).ToDictionaryAsync(x => x.Nonce, x => x);
                var outputMesages = new List<Message>();

                await foreach (var message in inputMessages.EnumerateItems(input.hashRegistry))
                {
                    Message[] resultMessages;
                    (agentStates[message.Target.AgentNonce], resultMessages) = await Execute(agentStates[message.Target.AgentNonce], message);

                    outputMesages.AddRange(resultMessages);
                }

                var outputStatesTree = await MerkleTreeBuilder.CreateRootFromValues(input.hashRegistry, agentStates.Values, 3); // FIXME: what about ordering?
                var outputState = new ChainState(outputStatesTree, chainState.NextAgentNonce);
                var outputMessagesTree = await MerkleTreeBuilder.CreateRootFromValues(input.hashRegistry, outputMesages, 3);

                return (outputState, outputMessagesTree);
            }

            async Task<bool> ValidateBlock(Block block, Hash<Block> expectedPrevious)
            {
                if (block.Previous != expectedPrevious)
                {
                    return false;
                }

                var previous = await input.hashRegistry.RetrieveAsync(expectedPrevious);

                // FIXME: messages are treated as a multiset elsewhere
                var inputMessagesSet = (await _state.GetValue<List<Message>>("messagePool")).ToHashSet();
                await foreach (var inputMessage in previous.InputMessages.EnumerateItems(input.hashRegistry))
                {
                    // if (inputMessage.Target.Chain != input.self || !inputMessage.Target.AllowedMessageTypes.Contains(inputMessage.Data.Type))
                    //     return false;

                    if (!inputMessagesSet.Remove(inputMessage))
                        return false;
                }

                var (outputState, outputMessages) = await ExecuteBlock(previous.State, block.InputMessages);

                return Hash.From(block.State) == Hash.From(outputState) && Hash.From(block.OutputMessages) == Hash.From(outputMessages);
            }

            async Task<Block?> ProposeBlock(Hash<Block> previousHash)
            {
                var previous = await input.hashRegistry.RetrieveAsync(previousHash);
                var inputMessagesList = (await _state.GetValue<List<Message>>("messagePool")).ToArray();

                if (inputMessagesList.Length == 0 && await _state.GetValue("finished", () => false))
                {
                    return null; // DEBUG: Used for testing purposes mainly
                }

                var inputMessages = await MerkleTreeBuilder.CreateRootFromValues(input.hashRegistry, inputMessagesList, 3);

                var (outputStates, outputMessages) = await ExecuteBlock(previous.State, inputMessages);

                return new Block(previousHash, inputMessages, outputMessages, outputStates);
            }

            var currentRound = await _state.GetValue<int>("currentRound", () => 0);

            if (currentRound == 0)
            {
                var snowball = await _state.GetValue<SnowballState>($"snowballState-{currentRound}");
                snowball.ProcessQuery(input.genesis);
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
                    var queryToSend = new Query(snowball.CurrentValue!.Value, currentRound);
                    replyTasks = sampledPeers.Select(peer => SendQuery(peer, queryToSend));
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
                    var finishedBlock = await input.hashRegistry.RetrieveAsync(finishedHash.Value);

                    var previousHash = await _state.GetValue<Hash<Block>>("lastBlock", () => input.genesis);
                    if (await ValidateBlock(finishedBlock, previousHash))
                    {
                        previousHash = finishedHash.Value;
                        await _state.SetValue("lastBlock", previousHash);

                        await foreach (var outputMessage in finishedBlock.OutputMessages.EnumerateItems(input.hashRegistry))
                        {
                            yield return outputMessage;
                        }

                        var messagePool = await _state.GetValue<List<Message>>("messagePool");
                        await foreach (var processedMessage in finishedBlock.InputMessages.EnumerateItems(input.hashRegistry))
                        {
                            messagePool.Remove(processedMessage);
                        }
                        await _state.SetValue("messagePool", messagePool);
                    }

                    await _state.SetValue("currentRound", ++currentRound);
                    // FIXME: Calculate proposers or proposal order from (previousHash, currentRound)
                    var newBlock = await ProposeBlock(previousHash);

                    if (newBlock == null) // DEBUG: Used for testing purposes mainly
                    {
                        break;
                    }

                    var newBlockHash = await input.hashRegistry.StoreAsync(newBlock);

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
            await _state.SetValue("finished", true);
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

        [FunctionName("PeeringResponder")]
        public async Task<IAsyncEnumerable<object>> PeeringResponder([PerperTrigger] (Peer other, IAsyncEnumerable<object> messages) input)
        {
            var request = (Query)await input.messages.FirstAsync();

            // NOTE: Should be using some locking here
            var snowballState = await _state.GetValue<SnowballState>($"snowballState-{request.Round}");
            snowballState.ProcessQuery(request.Value);
            var result = new Query(snowballState.CurrentValue!.Value, request.Round);
            await _state.SetValue($"snowballState-{request.Round}", snowballState);

            return new ListAsyncEnumerable<object>(new object[] { result });
        }
    }
}