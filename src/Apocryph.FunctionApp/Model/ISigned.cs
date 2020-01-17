using System.Collections.Generic;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public interface ISigned<out T>
    {
        T Value { get; }
        ValidatorKey Signer { get; }
        ValidatorSignature Signature { get; }
    }
}