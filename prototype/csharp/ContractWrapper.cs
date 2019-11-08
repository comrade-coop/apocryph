using System;
using System.Collections.Generic;
using System.Linq;
using Akka.Actor;

namespace Apocryph.Prototype
{
    public class ContractWrapper : UntypedActor
    {
        private sealed class FlushQueues {};

        private readonly Contract _contract;
        private readonly string _self = Context.Self.Path.Name;
        private readonly Dictionary<string, Message> _pendingMessages;
        private readonly Queue<string> _pendingMessageOrders;
        private int _nonce = 0;
        private bool _isProposer = false;

        public ContractWrapper(Contract contract)
        {
            _contract = contract;
            _pendingMessages = new Dictionary<string, Message>();
            _pendingMessageOrders = new Queue<string>();
        }

        protected void ProcessQueues()
        {
            while (_pendingMessageOrders.Count != 0 && _pendingMessages.ContainsKey(_pendingMessageOrders.Peek()))
            {
                var hash = _pendingMessageOrders.Dequeue();
                DispatchMessage(_pendingMessages[hash]);
                _pendingMessages.Remove(hash);
                Context.Parent.Tell(new ExecutionService.FlushResult(_pendingMessageOrders.Count));
            }

            if (_isProposer)
            {
                // Quotas??
                if (_pendingMessages.Count != 0) {
                    var hash = _pendingMessages.Keys.First();

                    var order = new MessageOrder(_self, hash);
                    Context.Parent.Tell(order);

                    DispatchMessage(_pendingMessages[hash]);
                    _pendingMessages.Remove(hash);

                    Context.Self.Tell(new FlushQueues());
                }
            }
        }

        protected void DispatchMessage(Message message)
        {
            Console.WriteLine(_self + " received a message from " + message.Sender);
            var results = _contract.Receive(message);
            foreach (var result in results)
            {
                _nonce ++;
                var messageToSend = new Message(_self, _self, _self + ":" + _nonce, result);
                Console.WriteLine(_self + " sends a message to " + messageToSend.Payload.Target);
                Context.ActorSelection("../" + messageToSend.Payload.Target).Tell(messageToSend);
            }
        }

        protected override void OnReceive(object received)
        {
            switch (received)
            {
                case Message message:
                    _pendingMessages[message.Hash] = message;
                    ProcessQueues();
                    break;
                case MessageOrder order:
                    _pendingMessageOrders.Enqueue(order.Hash);
                    ProcessQueues();
                    break;
                case ExecutionService.SetProposerState state:
                    _isProposer = state.Value;
                    ProcessQueues();
                    break;
                case ExecutionService.FlushChunk _:
                    Context.Parent.Tell(new ExecutionService.FlushResult(_pendingMessageOrders.Count));
                    break;
                case FlushQueues _:
                    ProcessQueues();
                    break;
                default:
                    break;
            }
        }
    }
}
