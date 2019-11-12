using System;
using System.Collections.Generic;

namespace Mock.Perper
{
    public class PerperStreamContext
    {
        public void CallActivityFunction(string functionName, params object[] arguments)
        {
            throw new Exception();
        }

        public IAsyncEnumerable<object> CallStreamFunction(string functionName, params object[] arguments)
        {
            throw new Exception();
        }

        public IReadOnlyDictionary<string, IAsyncEnumerable<object>> CallMultiStreamFunction(string functionName, params object[] arguments)
        {
            throw new Exception();
        }
    }

    public class PerperStreamContext<T> : PerperStreamContext
      where T : new()
    {
        public T State { get; } = new T();
    }
}
