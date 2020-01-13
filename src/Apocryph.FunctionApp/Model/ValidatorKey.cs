using System;
using System.Linq;
using System.Security.Cryptography;
using System.Text;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public struct ValidatorKey : IComparable<ValidatorKey>
    {
        public ECParameters Key { get; set; }

        public bool ValidateSignature(Hash hash, ValidatorSignature signature)
        {
            using (var ecdsa = ECDsa.Create(Key))
            {
                return ecdsa.VerifyHash(hash.Bytes, signature.Bytes);
            }
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
            return GetHashCode().CompareTo(other.GetHashCode());
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