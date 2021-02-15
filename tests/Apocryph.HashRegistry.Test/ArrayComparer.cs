using System;
using System.Collections.Generic;

namespace Apocryph.HashRegistry.Test
{
    public class ArrayComparer<T> : IEqualityComparer<T[]>
    {
        public IEqualityComparer<T> ValueComparer { get; set; } = EqualityComparer<T>.Default;

        public bool Equals(T[]? a, T[]? b)
        {
            if (a == null || b == null) return a == b;
            if (a.Length != b.Length) return false;

            for (var i = 0; i < a.Length; i++)
            {
                if (!ValueComparer.Equals(a[i], b[i]))
                {
                    return false;
                }
            }

            return true;
        }

        public int GetHashCode(T[] x)
        {
            throw new NotImplementedException();
        }
    }
}