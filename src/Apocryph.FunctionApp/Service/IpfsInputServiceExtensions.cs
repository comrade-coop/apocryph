using System;
using System.Collections.Generic;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Model;

namespace Apocryph.FunctionApp.Service
{
    public static class IpfsInputServiceExtensions
    {
        public static void SubscribeUserInput(this IAgentContext context, ValidatorKey key)
        {
            context.SendServiceMessage("IpfsInput", key);
        }
    }
}