using System.Collections.Generic;
using System.Text.Json;
using Apocryph.Ipfs.Serialization;

namespace Apocryph.Ipfs.Test
{
    public class SerializedComparer : IEqualityComparer<object?>
    {
        private SerializedComparer()
        {
        }

        public static IEqualityComparer<object?> Instance = new SerializedComparer();

        bool IEqualityComparer<object?>.Equals(object? a, object? b)
        {
            var aString = JsonSerializer.Serialize(a, ApocryphSerializationOptions.JsonSerializerOptions);
            var bString = JsonSerializer.Serialize(b, ApocryphSerializationOptions.JsonSerializerOptions);

            return aString.Equals(bString);
        }

        public int GetHashCode(object? x)
        {
            return JsonSerializer.Serialize(x, ApocryphSerializationOptions.JsonSerializerOptions).GetHashCode();
        }
    }
}