using System.Collections.Generic;

namespace Apocryph.FunctionApp.Model
{
    public interface IAgentStep : ISigned
    {
        IAgentStep Previous { get; set; } // FIXME: Should be hash stored on IPFS
        Dictionary<ValidatorKey, ValidatorSignature> CommitSignatures { get; set; }
    }
}