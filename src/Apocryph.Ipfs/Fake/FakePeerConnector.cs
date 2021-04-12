using System;
using System.Collections.Concurrent;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Ipfs.Serialization;

namespace Apocryph.Ipfs.Fake
{
    public class FakePeerConnectorProvider // FIXME filename
    {
        private ConcurrentDictionary<(Peer, string), Func<Peer, byte[], Task<byte[]>>> _queryListeners = new ConcurrentDictionary<(Peer, string), Func<Peer, byte[], Task<byte[]>>>();

        public FakePeerConnector GetConnector(Peer peer) // FIXME: GetPeerConnector?
        {
            return new FakePeerConnector(peer, this);
        }

        public FakePeerConnector GetConnector()
        {
            return GetConnector(GetFakePeer());
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
            Task<Peer> IPeerConnector.Self => Task.FromResult(Self);

            public FakePeerConnector(Peer self, FakePeerConnectorProvider factory)
            {
                Self = self;
                Factory = factory;
            }

            public async Task<TResult> Query<TRequest, TResult>(Peer peer, string path, TRequest request, CancellationToken cancellationToken = default)
            {
                var handler = Factory._queryListeners[(peer, path)];
                var requestBytes = JsonSerializer.SerializeToUtf8Bytes(request, ApocryphSerializationOptions.JsonSerializerOptions);
                var resultBytes = await handler(Self, requestBytes);
                return JsonSerializer.Deserialize<TResult>(resultBytes, ApocryphSerializationOptions.JsonSerializerOptions);
            }

            public Task ListenQuery<TRequest, TResult>(string path, Func<Peer, TRequest, Task<TResult>> handler, CancellationToken cancellationToken = default)
            {
                Factory._queryListeners.TryAdd((Self, path), async (peer, requestBytes) =>
                {
                    var request = JsonSerializer.Deserialize<TRequest>(requestBytes, ApocryphSerializationOptions.JsonSerializerOptions);
                    var result = await handler(peer, request);
                    return JsonSerializer.SerializeToUtf8Bytes(result, ApocryphSerializationOptions.JsonSerializerOptions);
                });

                cancellationToken.Register(() =>
                {
                    Factory._queryListeners.TryRemove((Self, path), out var _);
                });

                return Task.CompletedTask;
            }
        }
    }
}