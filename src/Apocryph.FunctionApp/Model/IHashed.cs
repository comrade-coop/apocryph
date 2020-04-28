using System.Collections.Generic;
using Newtonsoft.Json;
using Ipfs;

namespace Apocryph.FunctionApp.Model
{
    public interface IHashed<out T>
    {
        T Value { get; }
        Cid Hash { get; }
    }

    public static class Hashed
    {
        private class Hashed<T> : IHashed<T>
        {
            public Hashed(T value, Cid hash) {
                Value = value;
                Cid = hash;
            }

            public T Value { get; }
            public Cid Hash { get; }
        }

        public static IHashed<T> Create<T>(T value, Cid hash)
        {
            var type = typeof(Hashed<>).MakeGenericType(value.GetType());
            return (IHashed<T>)Activator.CreateInstance(type, value, hash);
        }
    }
}