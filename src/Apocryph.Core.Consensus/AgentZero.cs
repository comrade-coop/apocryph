using System;
using System.Collections.Generic;
using Apocryph.Core.Consensus.Blocks.Messages;
using Apocryph.Core.Consensus.Blocks.State;

namespace Apocryph.Core.Consensus
{
    public class AgentZeroState
    {
        public Dictionary<Guid, Balance> Balances { get; set; } = new Dictionary<Guid, Balance>();
        public Dictionary<Guid, Chain> Chains { get; set; } = new Dictionary<Guid, Chain>();
        public Dictionary<(Guid, Guid), Balance> Tickets { get; set; } = new Dictionary<(Guid, Guid), Balance>();
        public Balance RewardPool { get; set; } = new Balance();
    }

    public class AgentZero
    {
        public static ValueTuple<int, int> ChainRewardCut { get; } = (1, 10);
        public static ValueTuple<int, int> CallTicketCut { get; } = (1, 10);
        public static ValueTuple<int, int> RewardCut { get; } = (1, 100);

        public AgentZeroState Run(AgentZeroState state, object message, Guid? reference)
        {
            switch (message)
            {
                case TransferMessage transferMessage:
                    {
                        var sender = state.Balances.GetOrCreate(reference!.Value);
                        var receiver = state.Balances.GetOrCreate(transferMessage.To);
                        sender.Transfer(receiver, transferMessage.Amount);
                    }
                    break;

                case CreateChainMessage createChainMessage:
                    {
                        var creator = state.Balances.GetOrCreate(reference!.Value);
                        var chain = state.Balances.GetOrCreate(createChainMessage.ChainId);
                        creator.Transfer(chain, createChainMessage.InitialBalance);
                        state.Chains.Add(createChainMessage.ChainId, new Chain
                        {
                            LatestBlock = createChainMessage.InitialBlockId
                        });
                    }
                    break;

                case IssueTicketsMessage issueTicketsMessage:
                    {
                        var creator = state.Balances.GetOrCreate(reference!.Value);
                        var ticketBalance = state.Tickets[(issueTicketsMessage.For, issueTicketsMessage.Target)];
                        var cut = issueTicketsMessage.Amount * CallTicketCut.Item1 / CallTicketCut.Item2;
                        creator.Transfer(ticketBalance, issueTicketsMessage.Amount - cut);
                        creator.Transfer(state.RewardPool, cut);
                    }
                    break;

                case SetChainBlockMessage setChainBlockMessage:
                    {
                        var chain = state.Balances.GetOrCreate(setChainBlockMessage.ChainId);
                        var agentZeroProposer = state.Balances.GetOrCreate(reference!.Value);

                        foreach (var usedTicked in setChainBlockMessage.UsedTickets)
                        {
                            var (otherChain, amount) = usedTicked;
                            var ticket = state.Tickets[(setChainBlockMessage.ChainId, otherChain)];
                            ticket.Transfer(chain, amount);
                        }

                        foreach (var unlockedTicket in setChainBlockMessage.UnlockedTickets)
                        {
                            var (otherChain, amount) = unlockedTicket;
                            var ticket = state.Tickets[(setChainBlockMessage.ChainId, otherChain)];
                            var original = state.Balances.GetOrCreate(otherChain);
                            ticket.Transfer(original, amount);
                        }

                        foreach (var processedCommand in setChainBlockMessage.ProcessedCommands)
                        {
                            var (proposerReference, amount) = processedCommand;
                            var proposer = state.Balances[proposerReference];
                            var cut = amount * ChainRewardCut.Item1 / ChainRewardCut.Item2;
                            chain.Transfer(agentZeroProposer, cut);
                            chain.Transfer(proposer, amount - cut);
                        }

                        state.Chains[setChainBlockMessage.ChainId].LatestBlock = setChainBlockMessage.BlockId;
                    }
                    break;

                case ClaimRewardMessage _:
                    {
                        var agentZeroProposer = state.Balances.GetOrCreate(reference!.Value);
                        var cut = state.RewardPool.Amount * RewardCut.Item1 / RewardCut.Item2;
                        state.RewardPool.Transfer(agentZeroProposer, cut);
                    }
                    break;
            }

            return state;
        }
    }

    internal static class DictionaryExtensions
    {
        public static TValue GetOrCreate<TKey, TValue>(this IDictionary<TKey, TValue> dictionary, TKey key)
            where TKey : notnull
            where TValue : new()
        {
            if (dictionary.TryGetValue(key, out var value))
            {
                return value;
            }

            var newValue = new TValue();
            dictionary.Add(key, newValue);
            return newValue;
        }
    }
}