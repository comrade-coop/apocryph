using System.Text.Encodings.Web;
using System.Text.Json;

namespace Apocryph.Ipfs.Serialization
{
    public static class ApocryphSerializationOptions
    {
        public static readonly JsonSerializerOptions JsonSerializerOptions = new JsonSerializerOptions
        {
            Encoder = JavaScriptEncoder.UnsafeRelaxedJsonEscaping,
            Converters =
            {
                { new TypeConverter() },
                { new NonStringKeyDictionaryConverter() },
                { new ObjectParameterConstructorConverter() {
                    AllowSubtypes = true
                } }
            }
        };
    }
}