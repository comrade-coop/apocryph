using System;
using System.Collections.Generic;
using System.Numerics;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model.Message;
using Apocryph.FunctionApp.Wetonomy.Messages;

namespace Apocryph.FunctionApp.Wetonomy
{
    public static class TokenManagerAgent
    {
        public static async Task<(object, IEnumerable<object>)> Run(object state, IMessage message)
        {
	        var balances = state as Dictionary<string, BigInteger> ?? new Dictionary<string, BigInteger>();

	        var messages = new List<object>();
	        
	        switch (message)
	        {
		        case MintMessage mintAction:
			        messages.Add(AddTokens(balances, mintAction.To, mintAction.Amount));
			        break;

		        case TransferMessage transferAction:
			        messages.Add(RemoveTokens(balances, transferAction.From, transferAction.Amount));
			        messages.Add(AddTokens(balances, transferAction.To, transferAction.Amount));
			        break;

		        case BurnMessage burnAction:
			        messages.Add(RemoveTokens(balances, burnAction.From, burnAction.Amount));
			        break;
	        }

	        await Task.CompletedTask;
	        
	        return (balances, messages);
        }

        private static object RemoveTokens(IDictionary<string, BigInteger> balances, string from, BigInteger amount)
        {
	        if (!balances.ContainsKey(from))
	        {
		        throw new Exception("Not such account");
	        }
	        
	        if (balances[from] < amount)
	        {
		        throw new Exception("Not enough funds to burn");
	        }

	        balances[from] -= amount;
	        
	        var @event = new TokensChangedMessage {Change = -amount, Total = balances[from], Target = from};
			        
	        if (balances[from] == 0)
	        {
		        balances.Remove(from);
	        }

	        return @event;
        }

        private static object AddTokens(IDictionary<string, BigInteger> balances, string to, BigInteger amount)
        {
	        if (!balances.ContainsKey(to))
	        {
		        balances[to] = 0;
	        }

	        balances[to] += amount;
	        
	        return new TokensChangedMessage {Change = amount, Total = balances[to], Target = to};
        }
    }
}