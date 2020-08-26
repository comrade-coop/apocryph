using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.Core.Consensus;
using Apocryph.Core.Consensus.Blocks;
using Apocryph.Core.Consensus.Communication;
using Apocryph.Core.Consensus.VirtualNodes;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using System.Security.Cryptography;

namespace Apocryph.Runtime.FunctionApp
{
    public class LoggingStream
    {
        private Dictionary<(Node, Type), Report> _reports = new Dictionary<(Node, Type), Report>();
        private Dictionary<Guid, Node?[]>? _nodes;
        private Dictionary<Guid, Block> _blocks = new Dictionary<Guid, Block>();

        [FunctionName(nameof(LoggingStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("nodes")] Dictionary<Guid, Node?[]> nodes,
            [Perper("chain")] IAsyncEnumerable<Message<(Guid, Node?[])>> chain,
            [Perper("filter")] IAsyncEnumerable<Block> filter,
            [Perper("reports")] IAsyncEnumerable<Report> reports,
            CancellationToken cancellationToken)
        {
            _nodes = nodes;

            await TaskHelper.WhenAllOrFail(
                HandleChain(chain, cancellationToken),
                HandleFilter(filter, cancellationToken),
                HandleReports(reports, cancellationToken));
        }

        private async Task HandleFilter(IAsyncEnumerable<Block> filter, CancellationToken cancellationToken)
        {
            await foreach (var block in filter)
            {
                _blocks[block.ChainId] = block;
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
                _reports[(report.Source, report.GetType())] = report;
                UpdateReportDisplay();
            }
        }

        private Node?[] GetBlockProposers(Block lastBlock)
        {
            var serializedBlock = JsonSerializer.SerializeToUtf8Bytes(lastBlock, ConsensusStream.JsonSerializerOptions);
            var nodes = _nodes![lastBlock.ChainId];
            var proposerCount = nodes!.Length / 10 + 1; // TODO: Move constant to parameter
            return RandomWalk.Run(serializedBlock).Select(selected => nodes[(int)(selected.Item1 % nodes.Length)]).Take(proposerCount).ToArray();
        }

        private string GetBlockShorthash(Block lastBlock)
        {
            var serializedBlock = JsonSerializer.SerializeToUtf8Bytes(lastBlock, ConsensusStream.JsonSerializerOptions);
            using var sha256Hash = SHA256.Create();
            var hash = sha256Hash.ComputeHash(serializedBlock);
            return String.Concat(hash.Select(x => x.ToString("x2")).Take(5));
        }

        private enum NodeStatus
        {
            Offline = 0,
            Forked = 1,

        }

        private void UpdateReportDisplay()
        {
            var lineStart = "\x1b[K";
            var nodeColumn = "{0,10} ";
            var roundColumn = "{0,10}{1} ";
            var lastBlockColumn = "{0,10}{1} ";
            var nextBlockColumn = "{0,10} ";
            var confidenceColumn = "{0,10} ";
            var missingValue = "N/A";

            Console.WriteLine(lineStart);
            var linesWritten = 1;

            foreach (var (chainId, nodes) in _nodes!)
            {
                Console.WriteLine("{0}Chain: {1}", lineStart, chainId);
                Console.Write(lineStart);
                Console.Write(nodeColumn, "Node ID");
                Console.Write(roundColumn, "Round", " ");
                Console.Write(lastBlockColumn, "Last Block", " ");
                Console.Write(nextBlockColumn, "Next Block ");
                Console.Write(confidenceColumn, "Confidence");
                Console.WriteLine();
                linesWritten += 2;

                var lastChainBlock = _blocks[chainId];
                var lastChainBlockHash = GetBlockShorthash(lastChainBlock);
                var beta = nodes!.Length * 2; // TODO: Move constants to parameters

                foreach (var node in nodes)
                {
                    if (node != null)
                    {
                        Console.Write(lineStart);
                        Console.Write(nodeColumn, node);
                        if (_reports.TryGetValue((node, typeof(ConsensusReport)), out var _consensusReport))
                        {
                            var consensusReport = (ConsensusReport)_consensusReport;
                            var proposers = GetBlockProposers(consensusReport.LastBlock);
                            var blockHash = GetBlockShorthash(consensusReport.LastBlock);
                            Console.Write(roundColumn, consensusReport.Round, proposers[0] == node ? "*" : proposers.Contains(node) ? "+" : " ");
                            Console.Write(lastBlockColumn, blockHash, blockHash != lastChainBlockHash ? "!" : " ");
                        }
                        else
                        {
                            Console.Write(roundColumn, missingValue, " ");
                            Console.Write(lastBlockColumn, missingValue, " ");
                        }

                        if (_reports.TryGetValue((node, typeof(SnowballReport)), out var _snowballReport) && ((SnowballReport)_snowballReport).BlockCounts.Count != 0)
                        {
                            var snowballReport = (SnowballReport)_snowballReport;
                            var bestBlock = snowballReport.BlockCounts.Aggregate((p, n) => n.Value > p.Value ? n : p);
                            Console.Write(nextBlockColumn, GetBlockShorthash(bestBlock.Key));
                            Console.Write(confidenceColumn, $"{bestBlock.Value}/{beta}");
                        }
                        else
                        {
                            Console.Write(nextBlockColumn, missingValue);
                            Console.Write(confidenceColumn, missingValue);
                        }
                        Console.WriteLine();
                        linesWritten += 1;
                    }
                }

                Console.WriteLine("Last block: {0} by {1}", lastChainBlockHash, lastChainBlock.Proposer);
                linesWritten += 1;

                foreach (var inputCommand in lastChainBlock.InputCommands)
                {
                    Console.WriteLine("{0}  Input: {1}", lineStart, inputCommand);
                    linesWritten += 1;
                }

                foreach (var command in lastChainBlock.Commands)
                {
                    Console.WriteLine("{0}  Output: {1}", lineStart, command);
                    linesWritten += 1;
                }

                foreach (var (stateName, state) in lastChainBlock.States)
                {
                    Console.WriteLine("{0}  State: {1} = {2}", lineStart, stateName, Encoding.UTF8.GetString(state));
                    linesWritten += 1;
                }

                foreach (var (capabilityId, (stateName, methods)) in lastChainBlock.Capabilities)
                {
                    Console.WriteLine("{0}  Capability: {1} = {2} -> {3}", lineStart, capabilityId, stateName, string.Join(", ", methods));
                    linesWritten += 1;
                }
            }
            if (Console.CursorTop - linesWritten > 0)
            {
                Console.SetCursorPosition(0, Console.CursorTop - linesWritten);
            }
        }
    }
}