using System;
using System.Linq;
using System.Security.Cryptography;

namespace Apocryph.Core.Consensus.VirtualNodes
{
    public struct PrivateKey : IComparable<PrivateKey>
    {
        static public ECCurve Curve { get; } = PublicKey.Curve;
        public ECParameters PrivateParameters => new ECParameters { Curve = Curve, Q = PublicKey.Point, D = PrivateData };
        public PublicKey PublicKey { get; }
        public byte[] PrivateData { get; }

        public PrivateKey(ECParameters parameters)
        {
            PublicKey = new PublicKey(parameters);
            PrivateData = parameters.D;
        }

        public static PrivateKey Create()
        {
            using var dsa = ECDsa.Create(Curve);
            return new PrivateKey(dsa.ExportParameters(true));
        }

        public byte[] Sign(byte[] dataBytes)
        {
            using var ecdsa = ECDsa.Create(PrivateParameters);
            return ecdsa.SignData(dataBytes, HashAlgorithmName.SHA256);
        }

        public override bool Equals(object? obj)
        {
            if (obj is PrivateKey other)
            {
                return PublicKey.Equals(other.PublicKey) && PrivateData.SequenceEqual(other.PrivateData);
            }

            return false;
        }

        public int CompareTo(PrivateKey other)
        {
            var publicCompare = PublicKey.CompareTo(other.PublicKey);
            if (publicCompare != 0)
            {
                return publicCompare;
            }
            for (var i = 0; i < (PrivateData?.Length ?? 0); i++)
            {
                if (i < (other.PrivateData?.Length ?? 0))
                {
                    return -1;
                }
                var result = PrivateData![i].CompareTo(other.PrivateData![i]);
                if (result != 0)
                {
                    return result;
                }
            }
            return 1;
        }

        public override int GetHashCode()
        {
            return HashCode.Combine(PublicKey, Convert.ToBase64String(PrivateData));
        }

        public override string ToString()
        {
            return PublicKey.ToString() + "|" + "<private data>";
        }
    }
}