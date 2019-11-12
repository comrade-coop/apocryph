using System.Threading.Tasks;
using System.Collections.Generic;

namespace System.Linq
{
    public static class AsyncEnumerable {

        public static Task ForEachAsync<T>(this IAsyncEnumerable<T> enumerable, Action<T> f)
        {
            throw new Exception();
        }

        public static Task ForEachAsync<T>(this IAsyncEnumerable<T> enumerable, Func<T, Task> f)
        {
            throw new Exception();
        }

        public static IAsyncEnumerable<T> Race<T>(this IAsyncEnumerable<T> a, IAsyncEnumerable<T> b)
        {
            throw new Exception();
        }

        public static IAsyncEnumerable<TR> Race<TA, TB, TR>(this IAsyncEnumerable<TA> a, IAsyncEnumerable<TB> b)
            where TA : TR
            where TB : TR
        {
            throw new Exception();
        }
    }
}
