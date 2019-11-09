using System;

namespace Mock.Perper
{
    public class PerperStreamContext<T> where T:new()
    {
        public T State { get; } = new T();

        public object CallStreamFunction(string functionName, params object[] arguments)
        {
            throw new Exception();
        }
    }
}