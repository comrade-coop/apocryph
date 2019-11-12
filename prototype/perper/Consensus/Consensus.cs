using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Mock.Perper;
using Apocryph.Execution;

namespace Apocryph.Consensus
{
    public static class ConsensusPre
    {
        public class State
        {
            // public VotingRound VotingRound;
            public ConsensusState ConsensusState;
        }

        [FunctionName("ConsensusPre")]
        public static async Task RunAsync(PerperStreamContext<State> context,
            Validator self,
            IAsyncEnumerable<ConsensusMessage> inputStream,
            [PerperOutput("consensusState")] IAsyncCollector<ConsensusState> consensusStream,
            [PerperOutput("execution")] IAsyncCollector<IExecutionControl> executionControlStream,
            [PerperOutput("output")] IAsyncCollector<ConsensusMessage> outputStream)
        {
            await inputStream.ForEachAsync(async input =>
            {
                switch (input)
                {
                    case VoteMessage vote:
                        if (context.State.ConsensusState.ValidatorSet.Proposer == self)
                        {
                            await executionControlStream.AddAsync(new ExecutionControlProposeEnd());
                            await executionControlStream.FlushAsync();
                        }

                        // context.State.VotingRound.AddVote(input);
                        await consensusStream.AddAsync(context.State.ConsensusState);
                        await consensusStream.FlushAsync();
                        var lastChunk = await context.State.ConsensusState.LastChunk.GetValue();
                        await executionControlStream.AddAsync(new ExecutionControlDropSnapshots()
                        {
                            FinalSnapshot = lastChunk.Subheight,
                        });

                        if (context.State.ConsensusState.ValidatorSet.Proposer == self)
                        {
                            await executionControlStream.AddAsync(new ExecutionControlProposeStart()
                            {
                                FromSnapshot = 0,
                                ChunkTimeLimit = TimeSpan.FromSeconds(10),
                                ChunkCountLimit = 1024,
                            });
                            await executionControlStream.FlushAsync();
                        }
                        break;
                    case ChunkMessage chunk:
                        if (chunk.Signer == context.State.ConsensusState.ValidatorSet.Proposer && chunk.Signer != self)
                        {
                            await executionControlStream.AddAsync(new ExecutionControlValidate()
                            {
                                FromSnapshot = chunk.Subheight - 1, // TODO: Ensure that the subheight is valid
                                SnapshotIdentifier = chunk.PreviousChunk.ToString(),
                                Messages = chunk.Messages,
                            });
                            await executionControlStream.FlushAsync();
                        }
                        break;
                    default:
                        // Imagine we never get here
                        break;
                }
            });
        }
    }

    public static class ConsensusPost
    {
        public class State
        {
            // public VotingRound VotingRound;
            public ConsensusState ConsensusState;
            public int Subheight;
            public Hash<ChunkMessage> LastChunk;
        }

        [FunctionName("ConsensusPost")]
        public static async Task RunAsync(PerperStreamContext<State> context,
            Validator self,
            IAsyncEnumerable<ConsensusState> consensusStream,
            IAsyncEnumerable<IExecutionOutput> executionOutputStream,
            IAsyncCollector<ConsensusMessage> outputStream)
        {
            var _ = consensusStream.ForEachAsync(newConsensus => // TODO
            {
                context.State.Subheight = 0;
                context.State.ConsensusState = newConsensus;
            });

            await executionOutputStream.ForEachAsync(async output =>
            {
                switch (output)
                {
                    case ExecutionOutputValidate validate:
                        if (!validate.Valid) {
                            var newVote = new VoteMessage()
                            {
                                Signer = self,
                                Height = context.State.ConsensusState.Height,
                                LastValid = Hash.FromString<ChunkMessage>(validate.SnapshotIdentifier),
                            };
                            await outputStream.AddAsync(newVote);
                        }
                        break;
                    case ExecutionOutputPropose propose:
                        if (context.State.ConsensusState.ValidatorSet.Proposer == self) {
                            var newChunk = new ChunkMessage()
                            {
                                Signer = self,
                                Subheight = context.State.Subheight + 1,
                                PreviousChunk = context.State.LastChunk,
                                Messages = propose.Messages,
                            };
                            context.State.LastChunk = Hash.Create(newChunk);
                            context.State.Subheight = newChunk.Subheight;

                            Debug.Assert(propose.Snapshot == context.State.Subheight);

                            await outputStream.AddAsync(newChunk);
                            await outputStream.FlushAsync();
                        }
                        break;
                    default:
                        // Shouldn't get here
                        break;
                }
            });
        }
    }
}
