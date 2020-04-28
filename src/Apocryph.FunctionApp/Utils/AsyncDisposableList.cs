using System;
using System.Collections.Generic;
using System.Threading.Tasks;
// using System.Linq;

namespace Apocryph.FunctionApp.Utils
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