using System;
using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.AgentZero.Messages;
using Apocryph.FunctionApp.AgentZero.Publications;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;

namespace Apocryph.FunctionApp.AgentZero.State
{
    public class AgentsState
    {
        public struct EmptyStruct
        {
        }

        public IDictionary<string, EmptyStruct> Agents { get; set; } = new Dictionary<string, EmptyStruct>();

        public void RegisterAgent(string agentId)
        {
            Agents[agentId] = new EmptyStruct();
        }
    }
}