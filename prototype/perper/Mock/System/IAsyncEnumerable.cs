using System.Threading;
using System.Threading.Tasks;

namespace System.Collections.Generic
{
    public interface IAsyncDisposable
    {
        ValueTask DisposeAsync();
    }

    public interface IAsyncEnumerable<out T>
    {
        IAsyncEnumerator<T> GetAsyncEnumerator(CancellationToken cancellationToken = new CancellationToken());
    }

    public interface IAsyncEnumerator<out T> : IAsyncDisposable
    {
        ValueTask<bool> MoveNextAsync();
        T Current { get; }
    }
}
