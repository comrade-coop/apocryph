using System;
using System.Collections.Generic;
using Apocryph.FunctionApp.Agent;

namespace Apocryph.FunctionApp.Service
{
    public static class SampleServiceExtensions
    {
        public static void SampleStore(this IAgentContext context, string key, object data)
        {
            context.SendServiceMessage("Sample", Tuple.Create(key, data));
        }

        public static void SampleRestore(this IAgentContext context, string key)
        {
            context.SendServiceMessage("Sample", Tuple.Create(key));
        }
    }
}