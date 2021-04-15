using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Ipfs.Serialization;

namespace Apocryph.Ipfs.Fake
{
    public class FakePeerConnectorProvider
    {
        private ConcurrentDictionary<(Peer, string), Func<Peer, byte[], Task<byte[]>>> _queryListeners = new ConcurrentDictionary<(Peer, string), Func<Peer, byte[], Task<byte[]>>>();
        private ConcurrentDictionary<string, Action<Peer, byte[]>?> _gossipListeners = new ConcurrentDictionary<string, Action<Peer, byte[]>?>();

        public FakePeerConnector GetPeerConnector(Peer peer)
        {
            return new FakePeerConnector(peer, this);
        }

        public FakePeerConnector GetPeerConnector()
        {
            return GetPeerConnector(GetFakePeer());
        }

        private static Random _random = new Random();
        public static Peer GetFakePeer()
        {
            var bytes = new byte[16];
            _random.NextBytes(bytes);
            return new Peer(bytes);
        }

        public class FakePeerConnector : IPeerConnector
        {
            public FakePeerConnectorProvider Factory { get; }
            public Peer Self { get; }
            public HashSet<Task> PendingHandlerTasks { get; } = new HashSet<Task>();
            Task<Peer> IPeerConnector.Self => Task.FromResult(Self);


            public FakePeerConnector(Peer self, FakePeerConnectorProvider factory)
            {
                Self = self;
                Factory = factory;
            }

            private byte[] Serialize<T>(T value) =>
                JsonSerializer.SerializeToUtf8Bytes(value, ApocryphSerializationOptions.JsonSerializerOptions);

            private T Deserialize<T>(byte[] serialized) =>
                JsonSerializer.Deserialize<T>(serialized, ApocryphSerializationOptions.JsonSerializerOptions);

            public async Task<TResult> SendQuery<TRequest, TResult>(Peer peer, string path, TRequest request, CancellationToken cancellationToken = default)
            {
                var handler = Factory._queryListeners[(peer, path)];
                var requestBytes = Serialize(request);
                var resultBytes = await handler(Self, requestBytes);
                return Deserialize<TResult>(resultBytes);
            }

            public Task ListenQuery<TRequest, TResult>(string path, Func<Peer, TRequest, Task<TResult>> handler, CancellationToken cancellationToken = default)
            {
                Factory._queryListeners.TryAdd((Self, path), async (peer, requestBytes) =>
                {
                    var request = Deserialize<TRequest>(requestBytes);
                    var result = await handler(peer, request);
                    return Serialize(result);
                });

                cancellationToken.Register(() =>
                {
                    Factory._queryListeners.TryRemove((Self, path), out var _);
                });

                return Task.CompletedTask;
            }

            public Task SendPubSub<T>(string path, T message, CancellationToken token = default)
            {
                var messageBytes = Serialize(message);

                var handler = Factory._gossipListeners[path];
                handler?.Invoke(Self, messageBytes);

                return Task.CompletedTask;
            }

            public Task ListenPubSub<T>(string path, Func<Peer, T, Task<bool>> handler, CancellationToken cancellationToken = default)
            {
                Action<Peer, byte[]> wrappedHandler = (peer, messageBytes) =>
                {
                    var message = Deserialize<T>(messageBytes);
                    var task = handler(peer, message);
                    PendingHandlerTasks.Add(task);
                    task.ContinueWith((t) => Console.WriteLine("PubSub handler '{0}' exited with exception: {1}", path, t.Exception), TaskContinuationOptions.OnlyOnFaulted);
                    task.ContinueWith((t) => PendingHandlerTasks.Remove(task));
                };

                Factory._gossipListeners.AddOrUpdate(path, _ => wrappedHandler, (_, existingHandler) => existingHandler + wrappedHandler);

                cancellationToken.Register(() =>
                {
                    while (true)
                    {
                        var current = Factory._gossipListeners[path];
                        if (Factory._gossipListeners.TryUpdate(path, current - wrappedHandler, current))
                            break;
                    }
                });

                return Task.CompletedTask;
            }
        }
    }
}