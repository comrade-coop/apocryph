using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Akka.Actor;

namespace Apocryph.Prototype
{
    public class ExecutionService : UntypedActor
    {
        public sealed class SetProposerState
        {
            public SetProposerState(bool value)
            {
                Value = value;
            }

            public bool Value { get; }
        }

        public sealed class FlushChunk {}

        public sealed class FlushResult
        {
            public FlushResult(int pendingCount)
            {
                PendingCount = pendingCount;
            }

            public int PendingCount { get; }
        }

        public sealed class ResultTimeout {}

        private bool _isProposer = false;
        private List<MessageOrder> _chunk = new List<MessageOrder>();
        private Dictionary<string, int> _pendingMessageProgress = new Dictionary<string, int>();
        private readonly TimeSpan _resultTime = TimeSpan.FromMilliseconds(400);
        private ICancelable _resultTimer = null;

        public ExecutionService()
        {
            Context.ActorOf(Props.Create<ContractWrapper>(new SampleContract()), "A");
            Context.ActorOf(Props.Create<ContractWrapper>(new SampleContract()), "B");
            Context.ActorOf(Props.Create<ContractWrapper>(new SampleContract()), "C");
        }

        protected override void OnReceive(object received)
        {
            switch (received)
            {
                case Message message:
                    Context.ActorSelection(message.Payload.Target).Tell(message);
                    break;
                case MessageOrder order:
                    if (_isProposer)
                    {
                        _chunk.Add(order);
                    } else {
                        Context.ActorSelection(order.Target).Tell(order);
                    }
                    break;
                case FlushChunk flush:
                    if (_isProposer)
                    {
                        Context.Sender.Tell(_chunk);
                        _chunk = new List<MessageOrder>();
                    }
                    else
                    {
                        foreach (var child in Context.GetChildren())
                        {
                            child.Tell(flush);
                        }
                    }
                    break;
                case FlushResult result:
                    _pendingMessageProgress[Context.Sender.Path.Name] = result.PendingCount;

                    if (_resultTimer != null) {
                        _resultTimer.Cancel();
                    }

                    _resultTimer = Context.System.Scheduler.ScheduleTellOnceCancelable(_resultTime, Context.Self, new ResultTimeout(), Context.Self);

                    break;
                case ResultTimeout _:
                    Context.Parent.Tell(new FlushResult(_pendingMessageProgress.Values.Sum()));

                    if (_resultTimer != null) {

                        _resultTimer.Cancel();
                    }
                    _resultTimer = null;

                    break;
                case SetProposerState state:
                    _isProposer = state.Value;
                    foreach (var child in Context.GetChildren())
                    {
                        child.Tell(state);
                    }
                    break;
                default:
                    break;
            }
        }
    }
}
