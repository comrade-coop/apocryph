using System;
using System.Security.Cryptography;
using System.Text.Json;
using Apocryph.Ipfs.Serialization;

namespace Apocryph.Ipfs
{
    public static class Hash
    {
        public static Hash<T> From<T>(T value)
        {
            using var sha256Hash = SHA256.Create();
            return Hash.From<T>(sha256Hash, value);
        }

        public static Hash<T> From<T>(SHA256 sha256Hash, T value)
        {
            var serialized = JsonSerializer.SerializeToUtf8Bytes(value, ApocryphSerializationOptions.JsonSerializerOptions);
            return FromBytes<T>(sha256Hash, serialized);
        }

        public static Hash<T> FromBytes<T>(byte[] serialized)
        {
            using var sha256Hash = SHA256.Create();
            return FromBytes<T>(sha256Hash, serialized);
        }

        public static Hash<T> FromBytes<T>(SHA256 sha256Hash, byte[] serialized)
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

    public class Hash<T> : ByteSequence<Hash<T>>
    {
        public Hash(byte[] bytes)
            : base(bytes) { }

        public Hash<S> Cast<S>()
        {
            return new Hash<S>(Bytes);
        }
    }
}