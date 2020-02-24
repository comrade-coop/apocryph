using System;
using System.Collections.Generic;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public class Signed<T> : ISigned<T>
    {
        public Signed(T value, ValidatorKey signer, ValidatorSignature signature) {
            Value = value;
            Signer = signer;
            Signature = signature;
        }

        public T Value { get; }
        public ValidatorKey Signer { get; }
        public ValidatorSignature Signature { get; }
    }

    public static class Signed
    {
        public static ISigned<object> Create(object value, ValidatorKey signer, ValidatorSignature signature)
        {
            var type = typeof(Signed<>).MakeGenericType(value.GetType());
            return (ISigned<object>)Activator.CreateInstance(type, value, signer, signature);
        }
    }
}