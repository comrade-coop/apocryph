using System.Collections.Generic;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public class Hashed<T>
    {
        public Hashed(T from, Hash hash) {
            Value = from;
            Hash = hash;
        }

        public T Value { get; }
        public Hash Hash { get; }
    }
}