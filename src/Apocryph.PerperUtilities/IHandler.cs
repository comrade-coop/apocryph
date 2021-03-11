using System.Threading.Tasks;

namespace Apocryph.PerperUtilities
{
    public interface IHandler
    {
        Task InvokeAsync(object? parameters);
    }

    public interface IHandler<T>
    {
        Task<T> InvokeAsync(object? parameters);
    }
}