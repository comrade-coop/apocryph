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
        public IDictionary<string, byte[]> Agents { get; set; } = new Dictionary<string, byte[]>();

        public void SetAgentBlock(string agentId, byte[] block)
        {
            Agents[agentId] = block;
        }
    }
}