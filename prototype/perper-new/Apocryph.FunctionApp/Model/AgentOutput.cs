using System.Collections.Generic;
using Apocryph.FunctionApp.Command;

namespace Apocryph.FunctionApp.Model
{
    public class AgentOutput : IAgentStep
    {
        public object State { get; set; }
        public List<ICommand> Commands { get; set; }

        public IAgentStep Previous { get; set; }
        public Dictionary<ValidatorKey, ValidatorSignature> CommitSignatures { get; set; }

        public ValidatorKey Signer { get; set; }
        public ValidatorSignature Signature { get; set; }
    }
}