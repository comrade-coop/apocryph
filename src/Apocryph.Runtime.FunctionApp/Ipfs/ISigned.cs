using System.Collections.Generic;
using Newtonsoft.Json;

namespace Apocryph.Runtime.FunctionApp.Ipfs
{
    public interface ISigned<out T>
    {
        T Value { get; }
        ValidatorKey Signer { get; }
        ValidatorSignature Signature { get; }
    }

    public static class Signed
    {
        private class SignedImpl<T> : ISigned<T>
        {
            public SignedImpl(T value, ValidatorKey signer, ValidatorSignature signature) {
                Value = value;
                Signer = signer;
                Signature = signature;
            }

            public T Value { get; }
            public ValidatorKey Signer { get; }
            public ValidatorSignature Signature { get; }
        }

        public static ISigned<T> Create<T>(T value, ValidatorKey signer, ValidatorSignature signature)
        {
            var type = typeof(SignedImpl<>).MakeGenericType(value.GetType());
            return (ISigned<T>)Activator.CreateInstance(type, value, signer, signature);
        }
    }
}