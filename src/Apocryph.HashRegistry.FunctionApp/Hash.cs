using System;
using System.Linq;
using System.Security.Cryptography;
using System.Text.Json;
using Apocryph.HashRegistry.Serialization;

namespace Apocryph.HashRegistry
{
    // NOTE: interface named without I in front, as we are simulating a class hierarchy, but the inheriting "class" is a struct
    public interface Hash : IEquatable<Hash>
    {
        byte[] Bytes { get; }

        string ToString();

        public static Hash<T> From<T>(T value)
        {
            using var sha256Hash = SHA256.Create();
            return Hash.From<T>(sha256Hash, value);
        }

        public static Hash<T> From<T>(SHA256 sha256Hash, T value)
        {
            var serialized = JsonSerializer.SerializeToUtf8Bytes(value, ApocryphSerializationOptions.JsonSerializerOptions);
            return FromSerialized<T>(sha256Hash, serialized);
        }

        public static Hash<T> FromSerialized<T>(byte[] serialized)
        {
            using var sha256Hash = SHA256.Create();
            return Hash.FromSerialized<T>(sha256Hash, serialized);
        }

        public static Hash<T> FromSerialized<T>(SHA256 sha256Hash, byte[] serialized)
        {
            return new Hash<T>(sha256Hash.ComputeHash(serialized));
        }

        public static Hash<T> FromString<T>(string input)
        {
            var output = new byte[input.Length / 2];
            for (var i = 0; i < output.Length; i++)
            {
                output[i] = Convert.ToByte(input.Substring(i * 2, 2), 16);
            }
            return new Hash<T>(output);
        }
    }

    public struct Hash<T> : Hash
    {
        public byte[] Bytes { get; }

        public Hash(byte[] bytes)
        {
            Bytes = bytes;
        }

        public Hash<S> Cast<S>()
        {
            return new Hash<S>(Bytes);
        }

        public bool Equals(Hash? other)
        {
            if (ReferenceEquals(null, other)) return false;
            if (ReferenceEquals(this, other)) return true;
            return Bytes.SequenceEqual(other.Bytes);
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
            Array.ForEach(Bytes, hash.Add);
            return hash.ToHashCode();
        }

        // Via https://stackoverflow.com/a/311179
        public override string ToString()
        {
            return BitConverter.ToString(Bytes).Replace("-", "");
        }

        public static bool operator ==(Hash<T> left, Hash<T> right) => left.Bytes.SequenceEqual(right.Bytes);
        public static bool operator !=(Hash<T> left, Hash<T> right) => !left.Bytes.SequenceEqual(right.Bytes);
    }
}