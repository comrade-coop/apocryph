using System;
using System.Linq;

namespace Apocryph.FunctionApp.Model
{
    public class Hash
    {
        public byte[] Bytes { get; set; }

        public override bool Equals(object? other)
        {
            if (other is Hash otherHash)
            {
                return Bytes.SequenceEqual(otherHash.Bytes);
            }
            else
            {
                return false;
            }
        }

        public static bool operator ==(Hash a, Hash b)
        {
            return a.Bytes.SequenceEqual(b.Bytes);
        }

        public static bool operator !=(Hash a, Hash b)
        {
            return !a.Bytes.SequenceEqual(b.Bytes);
        }

        public override int GetHashCode()
        {
            unchecked
            {
                var result = 0;

                foreach (var b in Bytes)
                {
                    result *= 31;
                    result ^= b;
                }

                return result;
            }
        }

        public override string ToString()
        {
            return Convert.ToBase64String(Bytes);
        }
    }
}