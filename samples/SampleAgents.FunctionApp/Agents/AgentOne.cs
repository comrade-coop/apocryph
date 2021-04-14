using System;
using System.Collections.Generic;
using Apocryph.Consensus;
using Apocryph.Ipfs;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Triggers;

namespace SampleAgents.FunctionApp.Agents
{
    public static class AgentOne
    {
        public class AgentOneState
        {
            public int Accumulator { get; set; } = 0;
        }

        [FunctionName("AgentOne")]
        public static (AgentState, Message[]) Run([PerperTrigger] (Hash<Chain> chain, AgentState state, Message message) input)
        {
            var state = input.state.Data.Deserialize<AgentOneState>();
            var outputMessages = new List<Message>();
            if (input.message.Data.Type == typeof(PingPongMessage).FullName)
            {
                var message = input.message.Data.Deserialize<PingPongMessage>();
                Console.WriteLine("AgentOne {0}", message.Content);
                Console.WriteLine("AgentOneState {0}", state.Accumulator);

                state.Accumulator += message.Content.Length;
                outputMessages.Add(new Message(
                    message.Callback,
                    ReferenceData.From(new PingPongMessage(
                        callback: new Reference(input.chain, input.state.Nonce, new[] { typeof(PingPongMessage).FullName! }),
                        content: "PONG! " + message.Content,
                        accumulatedValue: state.Accumulator
                    ))
                ));
            }
            return (new AgentState(input.state.Nonce, ReferenceData.From(state), input.state.CodeHash), outputMessages.ToArray());
        }
    }
}