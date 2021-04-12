using System;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace Apocryph.Ipfs.Serialization
{
    public class TypeConverter : JsonConverter<Type?>
    {
        public override Type? Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
        {
            if (reader.TokenType == JsonTokenType.Null)
            {
                reader.Skip();
                return null;
            }

            if (reader.TokenType == JsonTokenType.String)
            {
                var typeName = reader.GetString();
                var type = Type.GetType(typeName);

                if (type == null)
                {
                    throw new JsonException("Unknown type name");
                }

                return type;
            }

            throw new JsonException($"Unexpected token {reader.TokenType}");
        }

        public override void Write(Utf8JsonWriter writer, Type? type, JsonSerializerOptions options)
        {
            if (type == null)
            {
                writer.WriteNullValue();
            }
            else
            {
                writer.WriteStringValue(type.AssemblyQualifiedName);
            }
        }
    }
}