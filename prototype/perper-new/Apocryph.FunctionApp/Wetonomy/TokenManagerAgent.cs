using System;
using System.Collections.Generic;
using System.Linq.Expressions;
using System.Numerics;
using System.Reflection;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Apocryph.FunctionApp.Wetonomy.Messages;
using Microsoft.Azure.WebJobs;
using Newtonsoft.Json.Linq;

namespace Apocryph.FunctionApp.Wetonomy
{
    public static class TokenManagerAgent
    {

	    //TODO: What can be accessed from the TriggerBinding - function attributes???, paramater attributes????

	    [PerperRouter(typeof(BurnMessage))]
	    [FunctionName("Burn")]
	    public static void Burn(
		    [PerperWorkerTrigger] IApocryphContext<IDictionary<string, BigInteger>> context,
		    [PerperWorker("sender")] string sender,
		    [PerperWorker("message")] BurnMessage message)
	    {
		    var publication = RemoveTokens(context.State, message.From, message.Amount);
		    context.MakePublication(publication);
	    }


	    [TargetException=PublicationMessage(message)]; AdentState=A1(a1State)
	    public static void DoSomething(ApocryphContext context, 
		    [AgentState]a1State
		    PublicationMessage message)
	    {
		    context.SendMessage();
		    context.AddReminder();
		    context.MakePublication();
	    }
	    
	    public static async Task<(object, Commands)> Run(object state, PublicationMessage message)
        {
	        
	        
	        var balances = state as Dictionary<string, BigInteger> ?? new Dictionary<string, BigInteger>();

	        var commands = new Commands();
	        
	        switch (message)
	        {
		        case MintMessage mintAction:
			        commands.Add(AddTokens(balances, mintAction.To, mintAction.Amount));
			        break;

		        case TransferMessage transferAction:
			        commands.Add(RemoveTokens(balances, transferAction.From, transferAction.Amount));
			        commands.Add(AddTokens(balances, transferAction.To, transferAction.Amount));
			        break;

		        case BurnMessage burnAction:
			        commands.Add(RemoveTokens(balances, burnAction.From, burnAction.Amount));
			        break;
	        }

	        await Task.CompletedTask;
	        
	        return (balances, commands);
        }
	    
	    
	    public static async Task<(object, Commands)> Run(object state, ServiceResultMessage message)
	    {
		    var balances = state as Dictionary<string, BigInteger> ?? new Dictionary<string, BigInteger>();

		    var commands = new Commands();
	        
		    switch (message)
		    {
			    case MintMessage mintAction:
				    commands.Add(AddTokens(balances, mintAction.To, mintAction.Amount));
				    break;

			    case TransferMessage transferAction:
				    commands.Add(RemoveTokens(balances, transferAction.From, transferAction.Amount));
				    commands.Add(AddTokens(balances, transferAction.To, transferAction.Amount));
				    break;

			    case BurnMessage burnAction:
				    commands.Add(RemoveTokens(balances, burnAction.From, burnAction.Amount));
				    break;
		    }

		    await Task.CompletedTask;
	        
		    return (balances, commands);
	    }
	    
	    public static async Task<(object, Commands)> Run(object state, MintMessage message)
	    {
		    var balances = state as Dictionary<string, BigInteger> ?? new Dictionary<string, BigInteger>();

		    var commands = new Commands();
	        
		    switch (message)
		    {
			    case MintMessage mintAction:
				    commands.Add(AddTokens(balances, mintAction.To, mintAction.Amount));
				    break;

			    case TransferMessage transferAction:
				    commands.Add(RemoveTokens(balances, transferAction.From, transferAction.Amount));
				    commands.Add(AddTokens(balances, transferAction.To, transferAction.Amount));
				    break;

			    case BurnMessage burnAction:
				    commands.Add(RemoveTokens(balances, burnAction.From, burnAction.Amount));
				    break;
		    }

		    await Task.CompletedTask;
	        
		    return (balances, commands);
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