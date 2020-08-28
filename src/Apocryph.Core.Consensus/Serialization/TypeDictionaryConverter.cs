using System;
using System.Collections.Generic;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace Apocryph.Core.Consensus.Serialization
{
    public class TypeDictionaryConverter : JsonConverterFactory
    {
        public override bool CanConvert(Type type)
        {
            return type.IsGenericType && type.GetGenericTypeDefinition() == typeof(Dictionary<,>) && type.GenericTypeArguments[0] == typeof(Type);
        }

        public override JsonConverter CreateConverter(Type type, JsonSerializerOptions options)
        {
            var convertedType = typeof(TypeDictionaryConverterInner<>).MakeGenericType(type.GenericTypeArguments[1]);

            return (JsonConverter)Activator.CreateInstance(convertedType)!;
        }

        private class TypeDictionaryConverterInner<TValue> : JsonConverter<Dictionary<Type, TValue>>
        {
            public override Dictionary<Type, TValue> Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
            {
                if (reader.TokenType != JsonTokenType.StartObject)
                {
                    throw new JsonException();
                }

                Dictionary<Type, TValue> dictionary = new Dictionary<Type, TValue>();

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
                    var type = Type.GetType(propertyName)!;

                    if (!typeof(TValue).IsAssignableFrom(type))
                    {
                        throw new JsonException();
                    }

                    var value = (TValue)JsonSerializer.Deserialize(ref reader, type, options);

                    dictionary.Add(type, value);
                }

                throw new JsonException();
            }

            public override void Write(Utf8JsonWriter writer, Dictionary<Type, TValue> dictionary, JsonSerializerOptions options)
            {
                writer.WriteStartObject();

                foreach (var (type, value) in dictionary)
                {
                    writer.WritePropertyName(type.AssemblyQualifiedName);

                    JsonSerializer.Serialize(writer, value, type, options);
                }

                writer.WriteEndObject();
            }
        }
    }
}