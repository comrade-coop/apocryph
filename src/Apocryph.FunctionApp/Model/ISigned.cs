using System.Collections.Generic;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public interface ISigned<out T>
    {
        IHashed<T> Hashed { get; }
        ValidatorKey Signer { get; }
        ValidatorSignature Signature { get; }

        [JsonIgnore]
        T Value { get => Hashed.Value; }
        [JsonIgnore]
        Hash Hash { get => Hashed.Hash; }
    }
}