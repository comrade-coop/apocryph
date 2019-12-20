using System;
using System.Collections.Generic;
using System.Linq.Expressions;
using System.Numerics;
using System.Reflection;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Model;
using Apocryph.FunctionApp.Wetonomy.Messages;
using Microsoft.Azure.WebJobs;
using Newtonsoft.Json.Linq;

namespace Apocryph.FunctionApp.Wetonomy
{
    public static class TokenManagerAgent
    {
	    //TODO: What can be accessed from the TriggerBinding - function attributes???, paramater attributes????

	    [PerperRouter(typeof(MintMessage))]
	    [FunctionName("Transfer")]
	    public static void Transfer(
		    [PerperWorkerTrigger] IAgentContext<IDictionary<string, BigInteger>> context,
		    [PerperWorker("sender")] string sender,
		    [PerperWorker("message")] MintMessage message)
	    {
		    var publication = AddTokens(context.State, message.To, message.Amount);
		    context.MakePublication(publication);
	    }
	    
	    [PerperRouter(typeof(BurnMessage))]
	    [FunctionName("Burn")]
	    public static void Burn(
		    [PerperWorkerTrigger] IAgentContext<IDictionary<string, BigInteger>> context,
		    [PerperWorker("sender")] string sender,
		    [PerperWorker("message")] BurnMessage message)
	    {
		    var publication = RemoveTokens(context.State, message.From, message.Amount);
		    context.MakePublication(publication);
	    }
	    
	    
	    [PerperRouter(typeof(TransferMessage))]
	    [FunctionName("Transfer")]
	    public static void Transfer(
		    [PerperWorkerTrigger] IAgentContext<IDictionary<string, BigInteger>> context,
		    [PerperWorker("sender")] string sender,
		    [PerperWorker("message")] TransferMessage message)
	    {
		    var removeTokensPublication = RemoveTokens(context.State, message.From, message.Amount);
		    var addTokensPublication = AddTokens(context.State, message.To, message.Amount);
		    context.MakePublication(removeTokensPublication);
		    context.MakePublication(addTokensPublication);
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
	        
	        var @event = new TokensChangedPublication {Change = -amount, Total = balances[from], Target = from};
			        
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
	        
	        return new TokensChangedPublication {Change = amount, Total = balances[to], Target = to};
        }
    }

    public class PerperWorkerAttribute : Attribute
    {
	    public PerperWorkerAttribute(string sender)
	    {
		    throw new NotImplementedException();
	    }
    }

    public class PerperWorkerTriggerAttribute : Attribute
    {
    }

    public class PerperRouterAttribute : Attribute
    {
	    public PerperRouterAttribute(Type expression)
	    {
		    throw new NotImplementedException();
	    }
    }
}