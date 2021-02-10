using System.Text.Encodings.Web;
using System.Text.Json;

namespace Apocryph.HashRegistry.Serialization
{
    public static class ApocryphSerializationOptions
    {
        public static readonly JsonSerializerOptions JsonSerializerOptions = new JsonSerializerOptions
        {
            Encoder = JavaScriptEncoder.UnsafeRelaxedJsonEscaping,
            Converters =
            {
                { new TypeDictionaryConverter() },
                { new NonStringKeyDictionaryConverter() },
                { new ObjectParameterConstructorConverter() {
                    AllowSubtypes = true
                } }
            }
        };
    }
}