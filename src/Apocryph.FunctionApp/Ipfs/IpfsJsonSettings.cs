using System;
using System.Collections.Generic;
using System.Linq;
using Apocryph.FunctionApp.Model;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Ipfs
{
    public static class IpfsJsonSettings
    {
        // Courtesy of https://stackoverflow.com/a/59310390
        public class CustomDictionaryConverter<TKey, TValue> : JsonConverter
        {
            public override bool CanConvert(Type objectType)
            {
                return objectType == typeof(Dictionary<TKey, TValue>);
            }

            public override void WriteJson(JsonWriter writer, object value, JsonSerializer serializer)
            {
                serializer.Serialize(writer, ((Dictionary<TKey, TValue>)value).ToList());
            }

            public override object ReadJson(JsonReader reader, Type objectType, object existingValue, JsonSerializer serializer)
            {
                return serializer.Deserialize<KeyValuePair<TKey, TValue>[]>(reader).ToDictionary(kv => kv.Key, kv => kv.Value);
            }
        }

        public static JsonSerializerSettings DefaultSettings { get; } = new JsonSerializerSettings
        {
            Converters = {
                new CustomDictionaryConverter<ValidatorKey, ValidatorSignature>(),
            },
            TypeNameHandling = TypeNameHandling.Auto,
            // TODO: Must set SerializationBinder as well!
        };
    }
}