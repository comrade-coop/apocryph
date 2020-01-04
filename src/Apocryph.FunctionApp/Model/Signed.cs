using System.Collections.Generic;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public class Signed<T> : ISigned<T>
    {
        public Signed(Hashed<T> from, ValidatorKey signer, ValidatorSignature signature) {
            _hashed = from;
            Signer = signer;
            Signature = signature;
        }

        [JsonPropertyAttribute("Hashed")]
        readonly protected Hashed<T> _hashed;

        [JsonIgnore]
        public IHashed<T> Hashed { get => _hashed; }
        public ValidatorKey Signer { get; }
        public ValidatorSignature Signature { get; }
    }
}