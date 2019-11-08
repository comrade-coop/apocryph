using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Akka.Actor;

namespace Apocryph.Prototype
{
    public class CoreServiceExecutor : UntypedActor
    {
        public sealed class ChunkInvalid
        {
            public ChunkInvalid(CoreService.Chunk chunk)
            {
                Chunk = chunk;
            }

            public CoreService.Chunk Chunk { get; }
        }

        public sealed class CheckResult
        {
        }

        private readonly IActorRef _execution;
        private readonly Queue<object> _pendingOperations = new Queue<object>();

        private CoreService.Chunk _currentChunk = null;
        private bool _isProposer = false;
        private bool _waitingForExecution = false;

        public CoreServiceExecutor()
        {
            _execution = Context.ActorOf(Props.Create<ExecutionService>());
            _execution.Tell(new ExecutionService.SetProposerState(false));
        }

        protected void ProcessChunkQueue()
        {
            if (_pendingOperations.Count > 0 && !_waitingForExecution) {
                var nextOperation = _pendingOperations.Dequeue();
                switch (nextOperation) {
                    case CoreService.Chunk chunk:
                        if (chunk.Previous != _currentChunk)
                        {
                            Context.Parent.Tell(new ChunkInvalid(chunk), Context.Self);
                            return;
                        }

                        if (!_isProposer)
                        {
                            foreach (var messageOrder in chunk.MessageOrders)
                            {
                                _execution.Tell(messageOrder);
                            }

                            _execution.Tell(new ExecutionService.FlushChunk());
                            _waitingForExecution = true;
                        }

                        _currentChunk = chunk;
                        break;

                    case ExecutionService.SetProposerState state:
                        _isProposer = state.Value;
                        _execution.Tell(state);
                        break;
                }
            }
        }

        protected override void OnReceive(object received)
        {
            switch (received)
            {
                case CoreService.Chunk chunk:
                    _pendingOperations.Enqueue(chunk);
                    ProcessChunkQueue();
                    break;
                case ExecutionService.SetProposerState state:
                    _pendingOperations.Enqueue(state);
                    ProcessChunkQueue();
                    break;
                case Message message:
                    _execution.Tell(message);
                    break;
                case ExecutionService.FlushChunk flush:
                    _execution.Tell(flush);
                    break;
                case ExecutionService.FlushResult result:
                    if (result.PendingCount > 0)
                    {
                        Console.WriteLine(result.PendingCount);
                        Context.Parent.Tell(new ChunkInvalid(_currentChunk), Context.Self);
                    }
                    _waitingForExecution = false;
                    ProcessChunkQueue();
                    break;
                case List<MessageOrder> flushedMessages:
                    Context.Parent.Tell(flushedMessages);
                    break;
                default:
                    break;
            }
        }
    }
}
