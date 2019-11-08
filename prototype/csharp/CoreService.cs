using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Akka.Actor;

namespace Apocryph.Prototype
{
    public class CoreService : UntypedActor
    {
        public class Signed {
            public Signed(string signer)
            {
                Signer = signer;
            }

            public string Signer { get; }
        }

        public sealed class Transaction : Signed {
            public Transaction(string signer, int nonce, MessagePayload payload)
                : base(signer)
            {
                Nonce = nonce;
                Payload = payload;
            }

            public int Nonce { get; }
            public MessagePayload Payload { get; }
        }

        public sealed class Block : Signed {
            public Block(string signer, List<PreCommit> commits, Block previous)
                : base(signer)
            {
                Commits = commits;
                Previous = previous;
            }

            public List<PreCommit> Commits { get; }
            public Block Previous { get; }
        }

        public sealed class Chunk : Signed {
            public Chunk(string signer, List<MessageOrder> messageOrders, Chunk previous, Block block)
                : base(signer)
            {
                MessageOrders = messageOrders;
                Previous = previous;
                Block = block;
                Height = previous != null ? previous.Height + 1 : 1;
            }

            // public List<Transaction> Transactions { get; }
            public List<MessageOrder> MessageOrders { get; }
            public Chunk Previous { get; }
            public Block Block { get; }
            public int Height { get; }
        }

        public class Vote : Signed {
            public Vote(string signer, Chunk lastChunk, Block currentBlock)
                : base(signer)
            {
                LastChunk = lastChunk;
                Block = currentBlock;
            }

            public Chunk LastChunk { get; }
            public Block Block { get; }
        }

        public sealed class PreVote : Vote {
            public PreVote(string signer, Chunk lastChunk, Block currentBlock)
                : base(signer, lastChunk, currentBlock) {}
        }
        public sealed class PreCommit : Vote {
            // TODO: Is this really necessary? Maybe the next proposer can just issue the new block, and by including the votes, prove that it is the next proposer?
            public PreCommit(string signer, Chunk lastChunk, Block currentBlock)
                : base(signer, lastChunk, currentBlock) {}
        }

        private sealed class TimeoutFinish {
            public TimeoutFinish(Block block)
            {
                Block = block;
            }

            public Block Block { get; }
        }

        private sealed class TimeoutChunk {
            public TimeoutChunk(Block block)
            {
                Block = block;
            }

            public Block Block { get; }
        }

        // IDEA: Add slashing for multiple finish votes for different chunks of the same block?

        private readonly ICanTell _output;
        private readonly IActorRef _executor;
        private readonly ValidatorSet _validatorSet;
        private readonly string _self = Context.Self.Path.Name;
        private readonly TimeSpan _chunkTime = TimeSpan.FromSeconds(1);
        private readonly TimeSpan _proposerTimeoutTime = TimeSpan.FromSeconds(4);

        private Block _currentBlock = null;
        private Chunk _currentChunk = null;
        private bool _voted = false;
        private ICancelable _proposerTimer = null;
        private Dictionary<Block, VotingRound> _prevoteRounds = new Dictionary<Block, VotingRound>();
        private Dictionary<Block, VotingRound> _precommitRounds = new Dictionary<Block, VotingRound>();

        public CoreService(ICanTell output, Dictionary<string, int> validatorStakes)
        {
            _output = output;
            _executor = Context.ActorOf(Props.Create<CoreServiceExecutor>());
            _validatorSet = new ValidatorSet(validatorStakes);

            AdvanceProposer(new List<PreCommit>());
        }

        protected void AdvanceProposer(List<PreCommit> commits)
        {
            if (_validatorSet.Proposer == _self)
            {
                if (_proposerTimer != null) {
                    _proposerTimer.Cancel();
                    _proposerTimer = null;
                }
                _executor.Tell(new ExecutionService.SetProposerState(false));
            }

            _validatorSet.AdvanceProposer();

            if (_validatorSet.Proposer == _self)
            {
                _currentBlock = new Block(_self, commits, _currentBlock);
                _output.Tell(_currentBlock, Context.Self);
                _executor.Tell(new ExecutionService.SetProposerState(true));
                _proposerTimer = Context.System.Scheduler.ScheduleTellRepeatedlyCancelable(_chunkTime, _chunkTime, Context.Self, new TimeoutChunk(_currentBlock), Context.Self);
            }

            Context.System.Scheduler.ScheduleTellOnceCancelable(_proposerTimeoutTime, Context.Self, new TimeoutFinish(_currentBlock), Context.Self);
        }

        protected void StartVote()
        {
            if (!_voted)
            {
                var pv = new PreVote(_self, _currentChunk, _currentBlock);
                Context.Self.Tell(pv);
                _output.Tell(pv, Context.Self);
                _voted = true;
            }
        }

        protected void UpdateToChunk(Chunk chunk)
        {
            var chunks = new List<Chunk>();
            var iteratedChunk = chunk;
            while (iteratedChunk != null && iteratedChunk.Height > (_currentChunk != null ? _currentChunk.Height : 0))
            {
                chunks.Add(iteratedChunk);
                iteratedChunk = iteratedChunk.Previous;
            }

            if (iteratedChunk != _currentChunk)
            {
                StartVote();
                return;
            }

            _currentChunk = chunk;

            foreach (var processedChunk in chunks)
            {
                _executor.Tell(processedChunk);
            }
        }

        protected override void OnReceive(object received)
        {
            Console.WriteLine("{3} {1}->{0} {2}", Context.Self.Path.Name, (received as Signed)?.Signer, received, _validatorSet.Proposer);
            switch (received)
            {
                case Transaction transaction:
                    _executor.Tell(new Message(
                        transaction.Signer,
                        transaction.Signer,
                        transaction.Signer + ":" + transaction.Nonce,
                        transaction.Payload));
                    break;
                case Chunk chunk:
                    if (chunk.Signer == _validatorSet.Proposer)
                    {
                        UpdateToChunk(chunk);
                    }
                    break;
                case Block block:
                    if (block.Signer == _validatorSet.Proposer && (_currentBlock == null || _currentBlock.Signer != _validatorSet.Proposer))
                    {
                        _currentBlock = block;
                    }
                    break;
                case CoreServiceExecutor.ChunkInvalid chunkInvalid:
                    _currentChunk = chunkInvalid.Chunk.Previous;
                    StartVote();
                    break;
                case TimeoutChunk timeout:
                    if (timeout.Block == _currentBlock) {
                        _executor.Tell(new ExecutionService.FlushChunk());
                    }
                    break;
                case List<MessageOrder> flushedMessages:
                    if (_validatorSet.Proposer == _self) {
                        _currentChunk = new Chunk(_self, flushedMessages, _currentChunk, _currentBlock);
                        _output.Tell(_currentChunk, Context.Self);
                        _executor.Tell(_currentChunk, Context.Self);
                    }
                    break;
                case TimeoutFinish timeout:
                    if (timeout.Block == _currentBlock || timeout.Block == _currentBlock.Previous) {
                        StartVote();
                    }
                    break;
                case PreVote prevote:
                    if (prevote.LastChunk.Block == prevote.Block || prevote.LastChunk.Block == prevote.Block.Previous) {
                        if(!_prevoteRounds.ContainsKey(prevote.Block)) {
                            _prevoteRounds[prevote.Block] = new VotingRound(_validatorSet.Stakes);
                        }

                        var round = _prevoteRounds[prevote.Block];
                        round.AddVote(prevote);

                        if (prevote.Block == _currentBlock && round.HasEnoughVotes()) {
                            var pc = new PreCommit(_self, round.GetOptimalChunk(), _currentBlock);
                            round.Finish();

                            Context.Self.Tell(pc, Context.Self);
                            _output.Tell(pc, Context.Self);
                        }
                    }

                    break;
                case PreCommit precommit:
                    if (precommit.LastChunk.Block == precommit.Block || precommit.LastChunk.Block == precommit.Block.Previous) {
                        if(!_precommitRounds.ContainsKey(precommit.Block)) {
                            _precommitRounds[precommit.Block] = new VotingRound(_validatorSet.Stakes);
                        }

                        var round = _precommitRounds[precommit.Block];
                        round.AddVote(precommit);

                        if (precommit.Block == _currentBlock && round.HasEnoughVotes()) {
                            _currentChunk = round.GetOptimalChunk();
                            round.Finish();
                            AdvanceProposer(round.GetVotes().OfType<PreCommit>().ToList());
                        }
                    }
                    break;
                default:
                    break;
            }
        }
    }
}
