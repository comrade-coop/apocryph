using System.Collections.Generic;

namespace Apocryph.FunctionApp.Model
{
    public interface IAgentStep : ISigned, IHashed
    {
        Hash PreviousHash { get; set; }
        Dictionary<ValidatorKey, ValidatorSignature> CommitSignatures { get; set; }
    }
}