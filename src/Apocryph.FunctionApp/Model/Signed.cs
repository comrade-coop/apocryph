using System.Collections.Generic;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public class Signed<T> : ISigned<T>
    {
        public Signed(IHashed<T> hashed, ValidatorKey signer, ValidatorSignature signature) {
            Hashed = hashed;
            Signer = signer;
            Signature = signature;
        }

        public IHashed<T> Hashed { get; }
        public ValidatorKey Signer { get; }
        public ValidatorSignature Signature { get; }
        public T Value => Hashed.Value;
        public Hash Hash => Hashed.Hash;
    }
}