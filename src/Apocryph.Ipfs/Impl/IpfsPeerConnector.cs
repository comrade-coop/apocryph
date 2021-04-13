using System;
using System.Buffers;
using System.Collections.Generic;
using System.IO;
using System.IO.Pipelines;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Runtime.CompilerServices;
using System.Text;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Ipfs.Serialization;
using Ipfs;
using Ipfs.Http;

namespace Apocryph.Ipfs.Impl
{
    public class IpfsPeerConnector : IPeerConnector
    {
        private readonly IpfsClient _ipfs;
        private readonly Random _random = new Random();

        public string QueryProtocolPrefix { get; set; } = "/x/apocryph/v0/";
        public string PubSubTopicPrefix { get; set; } = "apocryph/v0/";
        public Task<Peer> Self { get; }

        public IpfsPeerConnector(IpfsClient ipfs)
        {
            _ipfs = ipfs;
            Self = _ipfs.IdAsync().ContinueWith(x => new Peer(x.Result.Id.ToArray()));
        }

        public Task SendPubSub<T>(string path, T message, CancellationToken token = default)
        {
            var topic = PubSubTopicPrefix + path;
            var messageBytes = JsonSerializer.SerializeToUtf8Bytes(message, ApocryphSerializationOptions.JsonSerializerOptions);

            return _ipfs.PubSub.PublishAsync(topic, messageBytes, token);
        }

        public Task ListenPubSub<T>(string path, Func<Peer, T, Task<bool>> handler, CancellationToken cancellationToken = default)
        {
            var topic = PubSubTopicPrefix + path;

            return _ipfs.PubSub.SubscribeAsync(topic, ipfsMessage =>
            {
                var peer = new Peer(ipfsMessage.Sender.Id.ToArray());
                var message = JsonSerializer.Deserialize<T>(ipfsMessage.DataBytes, ApocryphSerializationOptions.JsonSerializerOptions);
                var task = handler(peer, message);
                task.ContinueWith((t) => Console.WriteLine("PubSub handler '{0}' exited with exception: {1}", path, t.Exception), TaskContinuationOptions.OnlyOnFaulted);
            }, cancellationToken);
        }

        public async Task<TResult> SendQuery<TRequest, TResult>(Peer peer, string path, TRequest request, CancellationToken cancellationToken = default)
        {
            var protocol = QueryProtocolPrefix + path;
            var peerId = new MultiHash(peer.Bytes);
            var port = _random.Next(49152, 65535);

            var forwardEndpoint = new IPEndPoint(IPAddress.Loopback, port);

            try
            {
                await _ipfs.DoCommandAsync("p2p/forward", cancellationToken, protocol, new[]
                {
                    $"arg=/ip4/{forwardEndpoint.Address}/tcp/{forwardEndpoint.Port}",
                    $"arg=/p2p/{peerId}"
                });

                var socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
                await socket.ConnectAsync(forwardEndpoint);

                await using var stream = new NetworkStream(socket, true);
                var enumerator = ReadNewlineStream(PipeReader.Create(stream), cancellationToken).GetAsyncEnumerator();

                var requestBytes = JsonSerializer.SerializeToUtf8Bytes(request, ApocryphSerializationOptions.JsonSerializerOptions);
                await WriteNewlineStream(stream, requestBytes, cancellationToken);

                if (!await enumerator.MoveNextAsync())
                {
                    throw new Exception("Stream ended prematurely");
                }

                var result = JsonSerializer.Deserialize<TResult>(enumerator.Current.ToArray(), ApocryphSerializationOptions.JsonSerializerOptions);

                stream.Close();

                return result;
            }
            finally
            {
                await _ipfs.DoCommandAsync("p2p/close", cancellationToken, "", new[]
                {
                    $"protocol={protocol}",
                    $"listen-address=/ip4/{forwardEndpoint.Address}/tcp/{forwardEndpoint.Port}",
                    $"target-address=/p2p/{peerId}"
                });
            }
        }

        public async Task ListenQuery<TRequest, TResult>(string path, Func<Peer, TRequest, Task<TResult>> handler, CancellationToken cancellationToken = default)
        {
            var protocol = QueryProtocolPrefix + path;
            var port = _random.Next(49152, 65535);
            var listenEndpoint = new IPEndPoint(IPAddress.Loopback, port);

            var listeningSocket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            listeningSocket.Bind(listenEndpoint);
            listeningSocket.Listen(120);

            await _ipfs.DoCommandAsync("p2p/listen", cancellationToken, protocol, new[]
            {
                $"arg=/ip4/{listenEndpoint.Address}/tcp/{listenEndpoint.Port}",
                "report-peer-id=true"
            });


            cancellationToken.Register(async () =>
            {
                listeningSocket.Shutdown(SocketShutdown.Both);
                listeningSocket.Dispose();
                await _ipfs.DoCommandAsync("p2p/close", default(CancellationToken), "", new[]
                {
                    $"protocol={protocol}",
                    $"target-address=/ip4/{listenEndpoint.Address}/tcp/{listenEndpoint.Port}",
                    $"listen-address=/p2p/{new MultiHash((await Self).Bytes)}"
                });
            });

            var _ = Task.Run(async () =>
            {
                while (!cancellationToken.IsCancellationRequested)
                {
                    var socket = await listeningSocket.AcceptAsync();
                    var _ = HandleQuery(socket, handler, cancellationToken);
                }
            });
        }

        private async Task HandleQuery<TRequest, TResult>(Socket socket, Func<Peer, TRequest, Task<TResult>> handler, CancellationToken cancellationToken)
        {
            try
            {
                await using var stream = new NetworkStream(socket, true);
                var enumerator = ReadNewlineStream(PipeReader.Create(stream), cancellationToken).GetAsyncEnumerator();

                if (!await enumerator.MoveNextAsync())
                {
                    throw new Exception("Stream ended prematurely");
                }

                var peerId = Encoding.UTF8.GetString(enumerator.Current.ToArray()); // NOTE: Can drop ToArray in .NET 5
                var peer = new Peer(Base58.Decode(peerId));

                while (await enumerator.MoveNextAsync())
                {
                    var request = JsonSerializer.Deserialize<TRequest>(enumerator.Current.ToArray(), ApocryphSerializationOptions.JsonSerializerOptions);

                    var result = await handler(peer, request);

                    var resultBytes = JsonSerializer.SerializeToUtf8Bytes(result, ApocryphSerializationOptions.JsonSerializerOptions);
                    await WriteNewlineStream(stream, resultBytes, cancellationToken);
                }

                stream.Close();
            }
            catch (Exception e)
            {
                Console.WriteLine(e);
            }
        }

        private async IAsyncEnumerable<ReadOnlySequence<byte>> ReadNewlineStream(PipeReader reader, [EnumeratorCancellation] CancellationToken cancellationToken)
        {
            while (true)
            {
                var readResult = await reader.ReadAsync();
                var position = readResult.Buffer.PositionOf((byte)'\n');
                if (position != null)
                {
                    var result = readResult.Buffer.Slice(0, position.Value);
                    yield return result;
                    reader.AdvanceTo(readResult.Buffer.GetPosition(1, position.Value));
                }
                else if (readResult.IsCompleted)
                {
                    break;
                }
            }
        }

        private async Task WriteNewlineStream(Stream stream, ReadOnlyMemory<byte> bytes, CancellationToken cancellationToken)
        {
            await stream.WriteAsync(bytes, cancellationToken);
            await stream.WriteAsync(new byte[] { (byte)'\n' }, cancellationToken);
        }
    }
}