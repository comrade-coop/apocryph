using System.Collections.Generic;

namespace Apocryph.FunctionApp.Model
{
    public class AgentInput : IAgentStep
    {
        public object State { get; set; }
        public string Sender { get; set; }
        public object Message { get; set; }

        public Hash Previous { get; set; }
        public Dictionary<ValidatorKey, ValidatorSignature> CommitSignatures { get; set; }
    }
}