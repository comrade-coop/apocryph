using System;
using System.Threading.Tasks;

namespace Apocryph.Runtime.FunctionApp
{
    public static class TaskHelper
    {
        public static Task WhenAllOrFail(params Task[] tasks)
        {
            var taskCompletionSource = new TaskCompletionSource<bool>();
            foreach (var task in tasks)
            {
                task.ContinueWith(t => { Console.WriteLine(t.Exception); taskCompletionSource.SetException(t.Exception!); }, TaskContinuationOptions.OnlyOnFaulted);
            }
            return Task.WhenAny(taskCompletionSource.Task, Task.WhenAll(tasks));
        }
    }
}