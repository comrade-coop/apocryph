using System.Collections.Generic;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public interface IHashed<out T>
    {
        T Value { get; }
        Hash Hash { get; }
    }
}