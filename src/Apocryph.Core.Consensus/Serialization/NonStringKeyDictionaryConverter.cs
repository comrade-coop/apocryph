using System;
using System.Collections.Generic;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace Apocryph.Core.Consensus.Serialization
{
    public class NonStringKeyDictionaryConverter : JsonConverterFactory
    {
        public override bool CanConvert(Type type)
        {
            return type.IsGenericType && type.GetGenericTypeDefinition() == typeof(Dictionary<,>) && type.GenericTypeArguments[0] != typeof(string);
        }

        public override JsonConverter CreateConverter(Type type, JsonSerializerOptions options)
        {
            var convertedType = typeof(NonStringKeyDictionaryConverterInner<,>).MakeGenericType(type.GetGenericArguments());

            return (JsonConverter)Activator.CreateInstance(convertedType)!;
        }

        private class NonStringKeyDictionaryConverterInner<TKey, TValue> : JsonConverter<Dictionary<TKey, TValue>>
            where TKey : notnull
        {
            public override Dictionary<TKey, TValue> Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
            {
                if (reader.TokenType != JsonTokenType.StartObject)
                {
                    throw new JsonException();
                }

                Dictionary<TKey, TValue> dictionary = new Dictionary<TKey, TValue>();

                while (reader.Read())
                {
                    if (reader.TokenType == JsonTokenType.EndObject)
                    {
                        return dictionary;
                    }

                    if (reader.TokenType != JsonTokenType.PropertyName)
                    {
                        throw new JsonException();
                    }

                    var propertyName = reader.GetString();
                    var key = JsonSerializer.Deserialize<TKey>(propertyName, options);

                    var value = JsonSerializer.Deserialize<TValue>(ref reader, options);

                    dictionary.Add(key, value);
                }

                throw new JsonException();
            }

            public override void Write(Utf8JsonWriter writer, Dictionary<TKey, TValue> dictionary, JsonSerializerOptions options)
            {
                writer.WriteStartObject();

                foreach (var (key, value) in dictionary)
                {
                    writer.WritePropertyName(JsonSerializer.Serialize(key, options));

                    JsonSerializer.Serialize(writer, value, options);
                }

                writer.WriteEndObject();
            }
        }
    }
}