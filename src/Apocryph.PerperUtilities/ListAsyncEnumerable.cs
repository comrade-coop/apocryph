using System.Collections.Generic;
using System.Linq;
using System.Threading;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.PerperUtilities
{
    [PerperData]
    public class ListAsyncEnumerable<T> : IAsyncEnumerable<T>
    {
        public IList<T> Value { get; private set; }

        public ListAsyncEnumerable(IList<T> value)
        {
            Value = value;
        }

        public IAsyncEnumerator<T> GetAsyncEnumerator(CancellationToken cancellationToken)
            => Value.ToAsyncEnumerable().GetAsyncEnumerator(cancellationToken);
    }
}