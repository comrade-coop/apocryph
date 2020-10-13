using System;
using System.Collections.Generic;
using System.Collections.Concurrent;
using System.Linq;
using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.Serialization;
using Apocryph.Core.Consensus.VirtualNodes;

namespace Apocryph.Runtime.Reporting
{
    public class Program
    {
        static readonly HttpClient HttpClient = new HttpClient();

        public static async Task Main(string[] args)
        {
            HttpClient.BaseAddress = new Uri(args.Length > 0 ? args[0] : "http://localhost:8901/");
            HttpClient.Timeout = TimeSpan.FromSeconds(40);

            var tasks = new List<Task>();
            var blocks = new ConcurrentDictionary<Guid, Block>();
            var nodes = new ConcurrentDictionary<Guid, Node[]>();
            var reports = new ConcurrentDictionary<Node, Dictionary<Type, Report>>();

            var chainIds = await Request<Guid[]>("chain");
            foreach (var chainId in chainIds)
            {
                tasks.Add(WrapTask(async () =>
                {
                    var hash = await Request<Hash>($"chain/{chainId}");
                    blocks[chainId] = await Request<Block>($"block/{hash}");
                }));
                tasks.Add(WrapTask(async () =>
                {
                    nodes[chainId] = await Request<Node[]>($"chain/{chainId}/node");
                }));
            }

            var reportingNodes = await Request<Node[]>($"node");
            foreach (var node in reportingNodes)
            {
                tasks.Add(WrapTask(async () =>
                {
                    reports[node] = await Request<Dictionary<Type, Report>>($"chain/{node.ChainId}/node/{node.Id}");
                }));
            }

            await Task.WhenAll(tasks);

            var columns = "{0,10} {1,10}{2} {3,10}{4} {5,10} {6,10}";

            foreach (var (chainId, lastChainBlock) in blocks)
            {
                Console.WriteLine("Chain: {0}", chainId);
                Console.WriteLine(columns, "Node ID", "Round", " ", "Last Block", " ", "Next Block ", "Confidence");

                var lastChainBlockHash = Hash.From(lastChainBlock).ToString().Substring(0, 8);
                var beta = (int)(Math.Log(nodes[chainId].Length) * 5) + 1; // TODO: Move constants to parameters

                foreach (var node in nodes[chainId])
                {
                    var values = new object[7] { "N/A", "N/A", " ", "N/A", " ", "N/A", "N/A" };
                    if (node != null && reports.ContainsKey(node))
                    {
                        values[0] = node;
                        var nodeReports = reports[node];
                        if (nodeReports.ContainsKey(typeof(ConsensusReport)))
                        {
                            var consensusReport = (ConsensusReport)nodeReports[typeof(ConsensusReport)];
                            var proposers = consensusReport.Proposers;
                            var blockHash = consensusReport.LastBlock.ToString().Substring(0, 8);
                            values[1] = consensusReport.Round;
                            values[2] = proposers[0] == node ? "*" : proposers.Contains(node) ? "+" : " ";
                            values[3] = blockHash;
                            values[4] = blockHash != lastChainBlockHash ? "!" : " ";
                        }

                        if (nodeReports.ContainsKey(typeof(SnowballReport)))
                        {
                            var snowballReport = (SnowballReport)nodeReports[typeof(SnowballReport)];
                            if (snowballReport.BlockCounts.Count > 0)
                            {
                                var bestBlock = snowballReport.BlockCounts.Aggregate((p, n) => n.Value > p.Value ? n : p);
                                values[5] = bestBlock.Key.ToString().Substring(0, 8);
                                values[6] = $"{bestBlock.Value}/{beta}";
                            }
                        }
                    }
                    Console.WriteLine(columns, values);
                }

                Console.WriteLine("Last block: {0} by {1}", lastChainBlockHash, lastChainBlock.Proposer);

                foreach (var inputCommand in lastChainBlock.InputCommands)
                {
                    Console.WriteLine("  Input: {0}", inputCommand);
                }

                foreach (var command in lastChainBlock.Commands)
                {
                    Console.WriteLine("  Output: {0}", command);
                }

                foreach (var (stateName, state) in lastChainBlock.States)
                {
                    Console.WriteLine("  State: {0} = {1}", stateName, Encoding.UTF8.GetString(state));
                }

                foreach (var (capabilityId, (stateName, methods)) in lastChainBlock.Capabilities)
                {
                    Console.WriteLine("  Capability: {0} = {1} -> {2}", capabilityId, stateName, string.Join(", ", methods));
                }
            }
        }

        private static async Task<T> Request<T>(string path, CancellationToken cancellationToken = default)
        {
            var result = await HttpClient.GetStreamAsync(path);
            return await JsonSerializer.DeserializeAsync<T>(result, ApocryphSerializationOptions.JsonSerializerOptions, cancellationToken);
            // var result = await HttpClient.GetByteArrayAsync(path);
            // Console.WriteLine("{0} -- {1}", path, Encoding.UTF8.GetString(result));
            // return JsonSerializer.Deserialize<T>(result, ApocryphSerializationOptions.JsonSerializerOptions);
        }

        private static Task WrapTask(Func<Task> f)
        {
            return f();
        }
    }
}