using System.Threading;
using System.Threading.Tasks;

namespace Apocryph.Ipfs
{
    public interface IHashResolver
    {
        Task<Hash<T>> StoreAsync<T>(T value, CancellationToken cancellationToken = default);
        Task<T> RetrieveAsync<T>(Hash<T> hash, CancellationToken cancellationToken = default);
    }
}