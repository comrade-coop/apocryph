using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace Apocryph.Runtime.FunctionApp.Utils
{
    public class AsyncDisposableList : List<IAsyncDisposable>, IAsyncDisposable
    {
        public async ValueTask DisposeAsync()
        {
            foreach (var x in this)
            {
                await x.DisposeAsync();
            }
        }
    }
}