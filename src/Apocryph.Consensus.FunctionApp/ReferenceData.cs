using System.Text.Json;
using Apocryph.HashRegistry.Serialization;

namespace Apocryph.Consensus
{
    public class ReferenceData
    {
        public string Type { get; }
        public byte[] Data { get; }
        public Reference[] References { get; }

        public ReferenceData(string type, byte[] data, Reference[] references)
        {
            Type = type;
            Data = data;
            References = references;
        }

        public static ReferenceData From(object? value, Reference[]? references = null)
        {
            return new ReferenceData(
                value?.GetType()?.FullName ?? "",
                JsonSerializer.SerializeToUtf8Bytes(value, ApocryphSerializationOptions.JsonSerializerOptions),
                references ?? new Reference[] { });
        }

        public T Deserialize<T>()
        {
            return JsonSerializer.Deserialize<T>(Data, ApocryphSerializationOptions.JsonSerializerOptions);
        }


    }
}