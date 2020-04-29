using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Security.Cryptography;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;

namespace Apocryph.Runtime.FunctionApp.Ipfs
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

            public override void WriteJson(JsonWriter writer, object? value, JsonSerializer serializer)
            {
                serializer.Serialize(writer, ((Dictionary<TKey, TValue>)value).ToList());
            }

            public override object? ReadJson(JsonReader reader, Type objectType, object existingValue, JsonSerializer serializer)
            {
                return serializer.Deserialize<KeyValuePair<TKey, TValue>[]>(reader)?.ToDictionary(kv => kv.Key, kv => kv.Value);
            }
        }

        public class ECParametersJsonConverter : JsonConverter<ECParameters>
        {
            public ECCurve Curve { get; set; } = ECCurve.NamedCurves.nistP521;

            public override void WriteJson(JsonWriter writer, ECParameters value, JsonSerializer serializer)
            {
                // Serialization of ECCurve messes up, so the next if is disabled
                /*if (!value.Curve.Equals(Curve))
                {
                    throw new Exception("Unexpected curve value");
                }*/

                serializer.Serialize(writer, value.Q, typeof(ECPoint));
            }

            public override ECParameters ReadJson(JsonReader reader, Type type, ECParameters _existingValue, bool what, JsonSerializer serializer)
            {
                var q = serializer.Deserialize<ECPoint>(reader);
                return new ECParameters
                {
                    Q = q,
                    Curve = Curve,
                };
            }
        }

        public static JsonSerializerSettings DefaultSettings { get; } = new JsonSerializerSettings
        {
            Converters = {
                new CustomDictionaryConverter<ValidatorKey, int>(),
                new ECParametersJsonConverter()
            },
            TypeNameHandling = TypeNameHandling.Auto,
            // TODO: Must set SerializationBinder as well!
        };

        public static JToken JTokenFromObject<T>(T value)
        {
            var serializer = JsonSerializer.Create(DefaultSettings);
            using var tokenWriter = new JTokenWriter();
            serializer.Serialize(tokenWriter, value, typeof(T));
            return tokenWriter.Token;
        }

        public static T ObjectFromJToken<T>(JToken token)
        {
            var serializer = JsonSerializer.Create(DefaultSettings);
            return token.ToObject<T>(serializer);
        }

        public static byte[] ObjectToBytes<T>(T value)
        {
            var json = JsonConvert.SerializeObject(value, typeof(T), DefaultSettings);
            return Encoding.UTF8.GetBytes(json);
        }

        public static T BytesToObject<T>(byte[] bytes)
        {
            var json = Encoding.UTF8.GetString(bytes);
            return JsonConvert.DeserializeObject<T>(json, IpfsJsonSettings.DefaultSettings);
        }
    }
}