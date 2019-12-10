using System.Collections.Generic;

namespace Apocryph.FunctionApp.Model
{
    public interface ISigned
    {
        // FIXME: Would be great if those would come directly from IPFS and thus skip the need of reimplementing the cryptography ourselves
        ValidatorKey Signer { get; set; }
        ValidatorSignature Signature { get; set; }
    }
}