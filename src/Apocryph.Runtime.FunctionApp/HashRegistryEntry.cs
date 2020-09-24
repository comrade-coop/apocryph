using Perper.WebJobs.Extensions.Config;

namespace Apocryph.Runtime.FunctionApp
{
    [PerperData]
    public class HashRegistryEntry
    {
        public string Hash { get; set; } // HACK: Using string, as Perper wouldn't allow querying with Hash or byte[]
        public object Value { get; set; }

        public HashRegistryEntry(string hash, object value)
        {
            Hash = hash;
            Value = value;
        }
    }
}