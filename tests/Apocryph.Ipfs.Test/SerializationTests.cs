using System.Text.Json;
using Apocryph.Ipfs.Serialization;
using Xunit;
using Xunit.Abstractions;

namespace Apocryph.Ipfs.Test
{
    public class SerializationTests
    {
        private readonly ITestOutputHelper _output;
        public SerializationTests(ITestOutputHelper output)
        {
            _output = output;
        }

        public string Serialize<T>(T value) =>
            JsonSerializer.Serialize(value, ApocryphSerializationOptions.JsonSerializerOptions);

        public T Deserialize<T>(string serialized) =>
            JsonSerializer.Deserialize<T>(serialized, ApocryphSerializationOptions.JsonSerializerOptions);

        [Theory]
        [MemberData(nameof(TestData.Data), MemberType = typeof(TestData))]
        public void Deserialize_WithExactType_ReturnsIdentical(object? data)
        {
            var type = data?.GetType() ?? typeof(object);

            var serializeMethod = typeof(SerializationTests).GetMethod(nameof(Serialize))!.MakeGenericMethod(type)!;
            var serialized = (string)serializeMethod.Invoke(this, new object?[] { data })!;

            _output.WriteLine("Serialized is: {0}", serialized);

            var deserializeMethod = typeof(SerializationTests).GetMethod(nameof(Deserialize))!.MakeGenericMethod(type)!;
            var deserialized = deserializeMethod.Invoke(this, new object?[] { serialized });

            Assert.Equal(deserialized, data);
        }

        [Theory]
        [MemberData(nameof(TestData.DataInterface), MemberType = typeof(TestData))]
        public void Deserialize_WithInterfaceType_ReturnsIdentical(IExample data)
        {
            var serialized = Serialize<IExample>(data);

            _output.WriteLine("Serialized is: {0}", serialized);

            var deserialized = Deserialize<IExample>(serialized);

            Assert.Equal(deserialized, data);
        }

        [Theory]
        [MemberData(nameof(TestData.DataNoInterface), MemberType = typeof(TestData))]
        public void Deserialize_WithoutInterfaceType_ThrowsJsonException(object data)
        {
            var serialized = Serialize<object>(data);

            _output.WriteLine("Serialized is: {0}", serialized);

            Assert.Throws<JsonException>(() => Deserialize<IExample>(serialized));
        }
    }
}