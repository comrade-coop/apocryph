using System.Collections.Generic;

namespace Apocryph.FunctionApp.Model
{
    public class Signed<T>
    {
        public Signed(Hashed<T> from, ValidatorKey signer, ValidatorSignature signature) {
            Hashed = from;
            Signer = signer;
            Signature = signature;
        }

        public Hashed<T> Hashed { get; }
        public ValidatorKey Signer { get; }
        public ValidatorSignature Signature { get; }

        public T Value { get => Hashed.Value; }
        public Hash Hash { get => Hashed.Hash; }
    }
}