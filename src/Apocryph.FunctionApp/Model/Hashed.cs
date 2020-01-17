using System;
using System.Collections.Generic;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public class Hashed<T> : IHashed<T>
    {
        public Hashed(T value, Hash hash) {
            Value = value;
            Hash = hash;
        }

        public T Value { get; }
        public Hash Hash { get; }
    }

    public static class Hashed
    {
        public static IHashed<object> Create(object value, Hash hash)
        {
            var type = typeof(Hashed<>).MakeGenericType(value.GetType());
            return (IHashed<object>)Activator.CreateInstance(type, value, hash);
        }
    }
}