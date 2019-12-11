using System.Collections.Generic;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public interface IHashed
    {
        [JsonIgnore]
        Hash Hash { get; set; }
    }
}