using System.Threading.Tasks;
using System.Text.Json;
using Apocryph.Core.Consensus.Serialization;
using Microsoft.AspNetCore.Http;

namespace Apocryph.Runtime.FunctionApp
{
    static class HttpResponseExtensions
    {
        public static Task WriteJsonAsync(this HttpResponse response, object? value)
        {
            var json = JsonSerializer.Serialize(value, ApocryphSerializationOptions.JsonSerializerOptions);
            return response.WriteAsync(json);
        }
    }
}