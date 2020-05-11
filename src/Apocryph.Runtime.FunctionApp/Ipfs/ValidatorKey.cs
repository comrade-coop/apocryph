using System;
using System.Linq;
using System.Security.Cryptography;
using Newtonsoft.Json;
using Ipfs;

namespace Apocryph.Runtime.FunctionApp.Ipfs
{
    public struct ValidatorKey : IComparable<ValidatorKey>
    {
        public ECParameters Key { get; set; }

        public static ValidatorSignature GenerateSignature(ECParameters privateKey, byte[] dataBytes)
        {
            using var ecdsa = ECDsa.Create(privateKey);
            return new ValidatorSignature{
                Bytes = ecdsa.SignData(dataBytes, HashAlgorithmName.SHA256)
            };
        }

        public bool ValidateSignature(byte[] dataBytes, ValidatorSignature signature)
        {
            using var ecdsa = ECDsa.Create(Key);
            return ecdsa.VerifyData(dataBytes, signature.Bytes, HashAlgorithmName.SHA256);
        }

        public override bool Equals(object? obj)
        {
            if (obj is ValidatorKey other)
            {
                return Key.Q.X.SequenceEqual(other.Key.Q.X) && Key.Q.Y.SequenceEqual(other.Key.Q.Y);
            }

            return false;
        }

        public int CompareTo(ValidatorKey other)
        {
            for (var i = 0; i < (Key.Q.Y?.Length ?? 0); i++)
            {
                if (i < (other.Key.Q.Y?.Length ?? 0))
                {
                    return -1;
                }
                var result = Key.Q.Y[i].CompareTo(other.Key.Q.Y[i]);
                if (result != 0)
                {
                    return result;
                }
            }
            return 1;
        }

        public byte[] GetPosition()
        {
            return Key.Q.X.Concat(new byte[]{ 0 }).ToArray();
        }

        public byte[] GetDifficulty(Cid agentId, byte[] salt)
        {
            var concatenated = (Key.Q.Y ?? new byte[]{}).Concat(salt ?? new byte[]{}).Concat(agentId.ToArray()).ToArray();
            using var sha256Hash = SHA256.Create();
            return sha256Hash.ComputeHash(concatenated).Concat(new byte[]{ 0 }).ToArray();
        }

        public override int GetHashCode()
        {
            return HashCode.Combine(Convert.ToBase64String(Key.Q.X), Convert.ToBase64String(Key.Q.Y));
        }

        public override string ToString()
        {
            return Convert.ToBase64String(Key.Q.X) + "|" + Convert.ToBase64String(Key.Q.Y);
        }
    }
}