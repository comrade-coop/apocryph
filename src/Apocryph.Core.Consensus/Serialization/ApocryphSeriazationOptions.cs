using System.Text.Json;
using System.Text.Encodings.Web;

namespace Apocryph.Core.Consensus.Serialization
{
    public class ApocryphSerializationOptions
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