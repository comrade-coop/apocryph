using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using System.Threading.Tasks;
using System.Text.Json;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.VirtualNodes;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Serialization;
using Ipfs.Http;
using Ipfs;
using Peer = Apocryph.Core.Consensus.VirtualNodes.Peer;

namespace Apocryph.Runtime.FunctionApp
{
    public class IpfsQueryStream
    {
        static public readonly string ProtocolName = "/x/apocryph/query/0.0";

        private Dictionary<Node, Peer> nodeMappings = new Dictionary<Node, Peer>();
        private Dictionary<Node, TaskCompletionSource<Query<Hash>>> _receiveCompletionSources = new Dictionary<Node, TaskCompletionSource<Query<Hash>>>();
        private IAsyncCollector<Query<Hash>>? _output;

        [FunctionName(nameof(IpfsQueryStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("self")] Peer self,
            [Perper("chain")] IAsyncEnumerable<(Node, Peer)> chain,
            [Perper("queries")] IAsyncEnumerable<Query<Hash>> queries,
            [Perper("output")] IAsyncCollector<Query<Hash>> output,
            CancellationToken cancellationToken)
        {
            _output = output;
            var ipfs = new IpfsClient();
            await TaskHelper.WhenAllOrFail(
                HandleChain(chain, cancellationToken),
                HandleQueries(ipfs, self, queries, cancellationToken),
                RunListener(ipfs, cancellationToken));
        }

        private async Task HandleChain(IAsyncEnumerable<(Node, Peer)> chain, CancellationToken cancellationToken)
        {
            await foreach (var (node, peer) in chain.WithCancellation(cancellationToken))
            {
                nodeMappings[node] = peer;
            }
        }

        private async Task HandleQueries(IpfsClient ipfs, Peer self, IAsyncEnumerable<Query<Hash>> queries, CancellationToken cancellationToken)
        {
            var random = new Random();
            var tasks = new List<Task>();

            await foreach (var query in queries.WithCancellation(cancellationToken))
            {
                if (query.Verb == QueryVerb.Response && _receiveCompletionSources.ContainsKey(query.Sender))
                {
                    _receiveCompletionSources[query.Sender].TrySetResult(query);
                }
                else if (query.Verb == QueryVerb.Request)
                {
                    var targetPeer = nodeMappings[query.Receiver];
                    if (!targetPeer.Equals(self))
                    {
                        var peerId = new MultiHash(targetPeer.Value);
                        var port = random.Next(49152, 65535);

                        var forwardEndpoint = new IPEndPoint(Dns.GetHostAddresses(ipfs.ApiUri.Host).First(), port);

                        await ipfs.DoCommandAsync("p2p/forward", cancellationToken, ProtocolName, new[]
                        {
                            $"arg=/ip4/{forwardEndpoint.Address}/tcp/{forwardEndpoint.Port}",
                            $"arg=/p2p/{peerId}"
                        });

                        var socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
                        await socket.ConnectAsync(forwardEndpoint);
                        tasks.Add(SendQuery(socket, query, cancellationToken));
                    }
                }
            }

            await Task.WhenAll(tasks);
        }

        private async Task SendQuery(Socket socket, Query<Hash> query, CancellationToken cancellationToken)
        {
            await using var stream = new NetworkStream(socket, true);
            await JsonSerializer.SerializeAsync(stream, query, ApocryphSerializationOptions.JsonSerializerOptions, cancellationToken);
            var response = await JsonSerializer.DeserializeAsync<Query<Hash>>(stream, ApocryphSerializationOptions.JsonSerializerOptions, cancellationToken);

            await _output!.AddAsync(response, cancellationToken);
        }

        private async Task RunListener(IpfsClient ipfs, CancellationToken cancellationToken)
        {
            var listenEndpoint = new IPEndPoint(IPAddress.Loopback, 8419);

            using var listeningSocket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            listeningSocket.Bind(listenEndpoint);
            listeningSocket.Listen(120);

            await ipfs.DoCommandAsync("p2p/listen", cancellationToken, ProtocolName, new[]
            {
                $"arg=/ip4/{listenEndpoint.Address}/tcp/{listenEndpoint.Port}"
            });

            var tasks = new List<Task>();
            while (!cancellationToken.IsCancellationRequested)
            {
                var socket = await listeningSocket.AcceptAsync();
                tasks.Add(ReceiveQuery(socket, cancellationToken));
            }

            await Task.WhenAll(tasks);
        }

        private async Task ReceiveQuery(Socket socket, CancellationToken cancellationToken)
        {
            await using var stream = new NetworkStream(socket, true);
            var query = await JsonSerializer.DeserializeAsync<Query<Hash>>(stream, ApocryphSerializationOptions.JsonSerializerOptions, cancellationToken);

            var taskCompletionSource = new TaskCompletionSource<Query<Hash>>();
            cancellationToken.Register(() => taskCompletionSource.TrySetCanceled());
            _receiveCompletionSources[query.Receiver] = taskCompletionSource;

            await _output!.AddAsync(query, cancellationToken);

            var result = await taskCompletionSource.Task;
            await JsonSerializer.SerializeAsync(stream, result, ApocryphSerializationOptions.JsonSerializerOptions, cancellationToken);
        }
    }
}