using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Routing;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class ReportsStream
    {
        private Dictionary<Node, Dictionary<Type, Report>> _reports = new Dictionary<Node, Dictionary<Type, Report>>();
        private Dictionary<Guid, Node?[]>? _nodes;
        private Dictionary<Guid, Hash> _blocks = new Dictionary<Guid, Hash>();
        private Func<Hash, Task<Block>>? _hashRegistryWorker;

        [FunctionName(nameof(ReportsStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("nodes")] Dictionary<Guid, Node?[]> nodes,
            [Perper("chain")] IAsyncEnumerable<Message<(Guid, Node?[])>> chain,
            [Perper("filter")] IAsyncEnumerable<Hash> filter,
            [Perper("reports")] IAsyncEnumerable<Report> reports,
            [Perper("hashRegistryWorker")] string hashRegistryWorker,
            CancellationToken cancellationToken)
        {
            _nodes = nodes;
            _hashRegistryWorker = hash => context.CallWorkerAsync<Block>(hashRegistryWorker, new { hash }, default);

            await TaskHelper.WhenAllOrFail(
                RunServer(cancellationToken),
                HandleChain(chain, cancellationToken),
                HandleFilter(filter, cancellationToken),
                HandleReports(reports, cancellationToken));
        }

        private async Task HandleFilter(IAsyncEnumerable<Hash> filter, CancellationToken cancellationToken)
        {
            await foreach (var hash in filter)
            {
                var block = await _hashRegistryWorker!(hash);
                _blocks[block!.ChainId] = hash;
            }
        }

        private async Task HandleChain(IAsyncEnumerable<Message<(Guid, Node?[])>> chain, CancellationToken cancellationToken)
        {
            await foreach (var message in chain.WithCancellation(cancellationToken))
            {
                var (chainId, nodes) = message.Value;

                _nodes![chainId] = nodes;
            }
        }

        private async Task HandleReports(IAsyncEnumerable<Report> reports, CancellationToken cancellationToken)
        {
            await foreach (var report in reports)
            {
                if (!_reports.ContainsKey(report.Source))
                {
                    _reports[report.Source] = new Dictionary<Type, Report>();
                }
                _reports[report.Source][report.GetType()] = report;
            }
        }

        private Task RunServer(CancellationToken cancellationToken)
        {
            var host = Host.CreateDefaultBuilder().ConfigureWebHostDefaults(webBuilder =>
            {
                webBuilder.UseUrls("http://localhost:8901");
                webBuilder.Configure(app =>
                {
                    app.UseRouting();

                    app.UseEndpoints(endpoints =>
                    {
                        endpoints.MapGet("/chain", WrapEndpoint(() => _nodes!.Keys.ToList()));
                        endpoints.MapGet("/chain/{Id:guid}", WrapEndpoint((values) =>
                        {
                            var id = new Guid((string)values["Id"]);
                            return _blocks[id];
                        }));
                        endpoints.MapGet("/chain/{Id:guid}/node", WrapEndpoint((values) =>
                        {
                            var id = new Guid((string)values["Id"]);
                            return _nodes![id];
                        }));
                        endpoints.MapGet("/chain/{Id:guid}/node/{Index:int}", WrapEndpoint((values) =>
                        {
                            var id = new Guid((string)values["Id"]);
                            var index = int.Parse((string)values["Index"]);
                            var node = new Node(id, index);
                            return _reports[node];
                        }));
                        endpoints.MapGet("/block/{Hash}", WrapEndpoint(async (values) =>
                        {
                            var hash = Hash.Parse((string)values["Hash"]);
                            var block = await _hashRegistryWorker!(hash);
                            return block;
                        }));
                        endpoints.MapGet("/node", WrapEndpoint((values) =>
                        {
                            return _reports.Keys.ToList();
                        }));
                    });
                });

                webBuilder.ConfigureLogging(logging =>
                {
                    logging.SetMinimumLevel(LogLevel.Error);
                });
            }).Build();

            return host.RunAsync(cancellationToken);
        }

        private RequestDelegate WrapEndpoint(Func<RouteValueDictionary, object?> wrapped)
        {
            return (context) => context.Response.WriteJsonAsync(wrapped(context.Request.RouteValues));
        }

        private RequestDelegate WrapEndpoint(Func<object?> wrapped)
        {
            return (context) => context.Response.WriteJsonAsync(wrapped());
        }

        private RequestDelegate WrapEndpoint(Func<RouteValueDictionary, Task<object?>> wrapped)
        {
            return async (context) => await context.Response.WriteJsonAsync(await wrapped(context.Request.RouteValues));
        }

        private RequestDelegate WrapEndpoint(Func<Task<object?>> wrapped)
        {
            return async (context) => await context.Response.WriteJsonAsync(await wrapped());
        }
    }
}