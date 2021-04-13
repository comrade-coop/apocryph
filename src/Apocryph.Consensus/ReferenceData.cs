using System;
using System.Linq;
using System.Text.Json;
using Apocryph.Ipfs.Serialization;

namespace Apocryph.Consensus
{
    public class ReferenceData : IEquatable<ReferenceData>
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

        public override bool Equals(object? other)
        {
            return other is ReferenceData otherReferenceData && Equals(otherReferenceData);
        }

        public bool Equals(ReferenceData? other)
        {
            return other != null && Type.Equals(other.Type) && Data.SequenceEqual(other.Data) && References.SequenceEqual(other.References);
        }

        public override int GetHashCode()
        {
            var hashCode = new HashCode();
            hashCode.Add(Type);
            Array.ForEach(Data, hashCode.Add);
            Array.ForEach(References, hashCode.Add);
            return hashCode.ToHashCode();
        }
    }
}