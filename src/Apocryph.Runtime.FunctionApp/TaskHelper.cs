using System.Threading.Tasks;

namespace Apocryph.Runtime.FunctionApp
{
    public static class TaskHelper
    {
        public static Task WhenAllOrFail(params Task[] tasks)
        {
//             return Task.WhenAll(tasks);
            var taskCompletionSource = new TaskCompletionSource<bool>();
            foreach (var task in tasks)
            {
                task.ContinueWith(t => taskCompletionSource.SetException(t.Exception!), TaskContinuationOptions.OnlyOnFaulted);
            }
            return Task.WhenAny(taskCompletionSource.Task, Task.WhenAll(tasks));
        }
    }
}