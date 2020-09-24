using System;
using System.Linq;
using System.Security.Cryptography;
using System.Text.Json;
using Apocryph.Core.Consensus.Serialization;

namespace Apocryph.Core.Consensus.Blocks
{
    public struct Hash : IEquatable<Hash>
    {
        public byte[] Value { get; set; }

        public Hash(byte[] value)
        {
            Value = value;
        }

        public bool Equals(Hash other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return Value.SequenceEqual(other.Value);
        }

        public override bool Equals(object? obj)
        {
            if (ReferenceEquals(null, obj)) return false;
            if (ReferenceEquals(this, obj)) return true;
            if (obj.GetType() != this.GetType()) return false;
            return Equals((Hash)obj);
        }

        public override int GetHashCode()
        {
            var hash = new HashCode();
            Array.ForEach(Value, hash.Add);
            return hash.ToHashCode();
        }

        // Via https://stackoverflow.com/a/311179
        public override string ToString()
        {
            return BitConverter.ToString(Value).Replace("-","");
        }

        // Via https://stackoverflow.com/a/311179
        public static Hash Parse(string hex)
        {
            byte[] bytes = new byte[hex.Length / 2];
            for (int i = 0; i < hex.Length; i += 2)
            {
                bytes[i / 2] = Convert.ToByte(hex.Substring(i, 2), 16);
            }

            return new Hash(bytes);
        }

        public static Hash From(SHA256 sha256Hash, object value)
        {
            var serialized = JsonSerializer.SerializeToUtf8Bytes(value, ApocryphSerializationOptions.JsonSerializerOptions);
            return new Hash(sha256Hash.ComputeHash(serialized));
        }

        public static Hash From(object value)
        {
            using var sha256Hash = SHA256.Create();
            return Hash.From(sha256Hash, value);
        }
    }
}