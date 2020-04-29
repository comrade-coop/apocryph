// NOTE: File is ignored by .csproj file

using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Apocryph.Runtime.FunctionApp.Utils;
using Ipfs;
using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp.Consensus
{
    public static class SnowballListener
    {
        public class State
        {
            public Cid CurrentColor { get; set; }
            public Dictionary<Cid, int> CurrentCounts { get; set; } = new Dictionary<Cid, int>();
        }

        [FunctionName(nameof(SnowballListener))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("port")] int port,
            [Perper("protocol")] string protocol,
            [Perper("ipfsGateway")] string ipfsGateway,
            [PerperStream("initialColorStream")] IAsyncEnumerable<Cid> initialColorStream,
            // [Perper("privateKey")] ECParameters privateKey,
            // [Perper("self")] ValidatorKey self,
            CancellationToken cancellationToken)
        {
            var k = 5;
            var ownAddress = IPAddress.Loopback;
            var gatewayAddress = IPAddress.Loopback;
            var listenEndpoint = new IPEndPoint(ownAddress, port);
            var listenBacklog = 128;
            var state = await context.FetchStateAsync<State>() ?? new State();
            // await context.UpdateStateAsync(state);

            var ipfs = new IpfsClient(ipfsGateway);

            await ipfs.DoCommandAsync("p2p/listen", cancellationToken, protocol, new []
            {
                $"arg=/ip4/{listenEndpoint.Address}/tcp/{listenEndpoint.Port}"
            });

            await Task.WhenAll(
                Task.Run(async () =>
                {
                    using var socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.IP);
                    socket.Bind(listenEndpoint);
                    socket.Listen(listenBacklog);

                    async Task acceptSocket(Socket socket)
                    {
                        using (socket)
                        {
                            await using var stream = new NetworkStream(socket, true);
                            state.CurrentColor.Write(stream);
                        }
                    };

                    var pendingListeners = new List<Task>();

                    while (!cancellationToken.IsCancellationRequested)
                    {
                        pendingListeners.Add(acceptSocket(await socket.AcceptAsync().WithCancellation(cancellationToken)));
                    }

                    await Task.WhenAll(pendingListeners);
                }),
                Task.Run(async () =>
                {
                    var random = new Random();
                    var addresses = await ipfs.Swarm.AddressesAsync(cancellationToken);

                    var sampled = addresses.OrderBy(n => random.Next()).Take(k);

                    var tasks = sampled.Select(async peer => {
                        var forwardPort = random.Next(49152, 65535);
                        var forwardEndpoint = new IPEndPoint(gatewayAddress, forwardPort);
                        try
                        {
                            await ipfs.DoCommandAsync("p2p/forward", cancellationToken, protocol, new []
                            {
                                $"arg=/ip4/{forwardEndpoint.Address}/tcp/{forwardEndpoint.Port}",
                                "arg=/p2p/" + peer.Id.ToString()
                            });

                            using var socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.IP);
                            socket.Connect(forwardEndpoint);

                            await using var stream = new NetworkStream(socket, true);

                            var hash = Cid.Read(stream);
                            state.CurrentCounts[hash] += 1;
                        }
                        finally
                        {
                            await ipfs.DoCommandAsync("p2p/close", cancellationToken, null, new [] {"target-address=/ip4/127.0.0.1/tcp/" + forwardPort});
                        }
                    });

                    await Task.WhenAll(tasks);
                }));

            try
            {
                await context.BindOutput(cancellationToken);
            }
            finally
            {
                await ipfs.DoCommandAsync("p2p/close", cancellationToken, null, new [] {"target-address=/ip4/127.0.0.1/tcp/" + port});
            }
        }
    }
}