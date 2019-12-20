using System.Collections.Generic;

namespace Apocryph.FunctionApp.Model
{
    public interface IAgentStep
    {
        Hash Previous { get; set; }

        Dictionary<ValidatorKey, ValidatorSignature> CommitSignatures { get; set; }
    }
}