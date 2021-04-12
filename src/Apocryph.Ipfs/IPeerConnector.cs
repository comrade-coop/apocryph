using System;
using System.Threading;
using System.Threading.Tasks;

namespace Apocryph.Ipfs
{
    public interface IPeerConnector
    {
        Task<Peer> Self { get; } // NOTE: Task<> is needed as Ipfs.Http.IpfsClient.IdAsync returns Task and we cannot inject the result asynchronously

        Task<TResult> Query<TRequest, TResult>(Peer peer, string path, TRequest request, CancellationToken token = default);
        Task ListenQuery<TRequest, TResult>(string path, Func<Peer, TRequest, Task<TResult>> handler, CancellationToken cancellationToken = default);

        // Task<IAsyncEnumerable<T>> Stream<T>(Peer peer, string path, IAsyncEnumerable<T>, CancellationToken token);
        // IDisposable ListenStream(string path, Func<Peer, IAsyncEnumerable<T>, IAsyncEnumerable<T>> handler);
    }
}