using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace Apocryph.Ipfs.Serialization
{
    // NOTE: Can remove in .NET 5
    public class ObjectParameterConstructorConverter : JsonConverterFactory
    {
        public bool AllowSubtypes { get; set; } = false;
        public string TypeProperty { get; set; } = "$type";

        public override bool CanConvert(Type type)
        {
            return !type.IsArray && !type.IsPrimitive && type != typeof(string) && type != typeof(Guid) && !typeof(IEnumerable).IsAssignableFrom(type) && !type.GetConstructors().Any(c => c.GetParameters().Length == 0);
        }

        public override JsonConverter CreateConverter(Type type, JsonSerializerOptions options)
        {
            var convertedType = typeof(ObjectParameterConstructorConverterInner<>).MakeGenericType(type);

            return (JsonConverter)Activator.CreateInstance(convertedType, new object[] { this })!;
        }

        private class ObjectParameterConstructorConverterInner<T> : JsonConverter<T>
        {
            private readonly ObjectParameterConstructorConverter Factory;
            public ObjectParameterConstructorConverterInner(ObjectParameterConstructorConverter factory)
            {
                Factory = factory;
            }

            public override T Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
            {
                if (reader.TokenType != JsonTokenType.StartObject)
                {
                    throw new JsonException();
                }

                Dictionary<string, object?> values = new Dictionary<string, object?>();
                Dictionary<string, string> valuesLower = new Dictionary<string, string>();

                while (reader.Read())
                {
                    if (reader.TokenType == JsonTokenType.EndObject)
                    {
                        break;
                    }

                    if (reader.TokenType != JsonTokenType.PropertyName)
                    {
                        throw new JsonException();
                    }

                    var propertyName = reader.GetString();

                    if (Factory.AllowSubtypes && propertyName == Factory.TypeProperty)
                    {
                        var newType = JsonSerializer.Deserialize<Type>(ref reader, options);

                        if (!typeToConvert.IsAssignableFrom(newType))
                        {
                            throw new JsonException($"Type {newType} does not match base type {typeToConvert}");
                        }

                        typeToConvert = newType;
                        continue;
                    }

                    var property = typeToConvert.GetProperty(propertyName);

                    if (property != null && property.GetCustomAttributes(typeof(JsonIgnoreAttribute), false).Length == 0)
                    {
                        var value = JsonSerializer.Deserialize(ref reader, property.PropertyType, options);

                        values.Add(propertyName, value);
                        valuesLower[NormalizeName(propertyName)] = propertyName;
                        continue;
                    }

                    var field = typeToConvert.GetField(propertyName);
                    if (field != null && field.GetCustomAttributes(typeof(JsonIgnoreAttribute), false).Length == 0 && !field.IsStatic)
                    {
                        var value = JsonSerializer.Deserialize(ref reader, field.FieldType, options);

                        values.Add(propertyName, value);
                        valuesLower[NormalizeName(propertyName)] = propertyName;
                        continue;
                    }

                    reader.Skip();
                }

                if (reader.TokenType != JsonTokenType.EndObject)
                {
                    throw new JsonException();
                }

                foreach (var constructor in typeToConvert.GetConstructors())
                {
                    var parameters = constructor.GetParameters();
                    var matchedParameters = 0;
                    foreach (var parameter in parameters)
                    {
                        if (values.ContainsKey(parameter.Name!) || valuesLower.ContainsKey(NormalizeName(parameter.Name!)))
                        {
                            matchedParameters++;
                        }
                        else if (!parameter.IsOptional)
                        {
                            matchedParameters = -1;
                            break;
                        }
                    }
                    if (matchedParameters < 0)
                    {
                        continue;
                    }

                    var arguments = new object?[parameters.Length];
                    for (var i = 0; i < matchedParameters; i++)
                    {
                        if (values.ContainsKey(parameters[i].Name!))
                        {
                            arguments[i] = values[parameters[i].Name!];
                            values.Remove(parameters[i].Name!);
                        }
                        else
                        {
                            var key = valuesLower[NormalizeName(parameters[i].Name!)];
                            arguments[i] = values[key];
                            values.Remove(key);
                        }
                    }
                    for (var i = matchedParameters; i < parameters.Length; i++)
                    {
                        arguments[i] = parameters[i].DefaultValue;
                    }

                    var result = (T)constructor.Invoke(arguments)!;

                    foreach (var (key, value) in values)
                    {
                        var property = typeToConvert.GetProperty(key);
                        if (property != null && property.CanWrite)
                        {
                            property.SetValue(result, value);
                            continue;
                        }

                        var field = typeToConvert.GetField(key);
                        if (field != null && !field.IsStatic)
                        {
                            field.SetValue(result, value);
                        }
                    }

                    return result;
                }

                throw new JsonException($"No matching constructor for type {typeToConvert} ({typeof(T)})");
            }

            private string NormalizeName(string name)
            {
                return new string(name.Where(char.IsLetterOrDigit).Select(x => char.ToLower(x)).ToArray());
            }

            public override void Write(Utf8JsonWriter writer, T value, JsonSerializerOptions options)
            {
                writer.WriteStartObject();

                var type = typeof(T);

                if (value?.GetType() != typeof(T) && Factory.AllowSubtypes)
                {
                    type = value?.GetType()!;
                    writer.WritePropertyName(Factory.TypeProperty);
                    JsonSerializer.Serialize(writer, type, options);
                }

                foreach (var property in type.GetProperties())
                {
                    if (!property.CanRead || property.GetCustomAttributes(typeof(JsonIgnoreAttribute), false).Length > 0 || property.IsSpecialName)
                    {
                        continue;
                    }

                    writer.WritePropertyName(property.Name);

                    var propertyValue = property.GetValue(value);

                    JsonSerializer.Serialize(writer, propertyValue, property.PropertyType, options);
                }

                foreach (var field in type.GetFields())
                {
                    if (field.IsStatic || field.GetCustomAttributes(typeof(JsonIgnoreAttribute), false).Length > 0)
                    {
                        continue;
                    }

                    writer.WritePropertyName(field.Name);

                    var propertyValue = field.GetValue(value);

                    JsonSerializer.Serialize(writer, propertyValue, field.FieldType, options);
                }

                writer.WriteEndObject();
            }
        }
    }
}