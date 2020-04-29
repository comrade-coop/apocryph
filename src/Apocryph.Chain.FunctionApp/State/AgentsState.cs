using System;
using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using Apocryph.Agent;
using Apocryph.Chain.FunctionApp.Messages;
using Apocryph.Chain.FunctionApp.Publications;
using Microsoft.Azure.WebJobs;

namespace Apocryph.Chain.FunctionApp.State
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