using System.Linq;
using System.Text.Json;
using Apocryph.HashRegistry.Serialization;
using Xunit;
using Xunit.Abstractions;

namespace Apocryph.HashRegistry.Test
{
    public class SerializationTests
    {

        #region Sample data
        public interface IExample
        {
        }

        public class Example : IExample
        {
            public int Number { get; set; }

            public override bool Equals(object? obj) => obj is Example other && Number == other.Number;
            public override int GetHashCode() => Number.GetHashCode();
            public override string ToString() => $"Example({Number})";
        }

        public static object[][] SerializationExamples = new[] {
            new object[] { 42 },
            new object[] { "Example text" },
            new object[] { new Example { Number = 42 } },
            new object[] { (42, 103, "Tuple") }
        };

        public static object[][] SerializationExamplesInterface = SerializationExamples.Where(x => x[0] is IExample).ToArray();
        #endregion Sample data

        private readonly ITestOutputHelper _output;
        public SerializationTests(ITestOutputHelper output)
        {
            _output = output;
        }

        public string Serialize<T>(T value) => JsonSerializer.Serialize(value, ApocryphSerializationOptions.JsonSerializerOptions);

        public T Deserialize<T>(string serialized) => JsonSerializer.Deserialize<T>(serialized, ApocryphSerializationOptions.JsonSerializerOptions);

        [Theory]
        [MemberData(nameof(SerializationExamples))]
        public void Deserialize_WithExactType_ReturnsIdentical(object data)
        {
            var type = data.GetType();

            var serializeMethod = typeof(SerializationTests).GetMethod(nameof(Serialize))!.MakeGenericMethod(type)!;
            var serialized = (string)serializeMethod.Invoke(this, new object[] { data })!;

            _output.WriteLine("Serialized is: {0}", serialized);

            var deserializeMethod = typeof(SerializationTests).GetMethod(nameof(Deserialize))!.MakeGenericMethod(type)!;
            var deserialized = deserializeMethod.Invoke(this, new object[] { serialized });

            Assert.Equal(deserialized, data);
        }

        [Theory]
        [MemberData(nameof(SerializationExamplesInterface))]
        public void Deserialize_WithInterfaceType_ReturnsIdentical(IExample data)
        {
            var serialized = Serialize<IExample>(data);

            _output.WriteLine("Serialized is: {0}", serialized);

            var deserialized = Deserialize<IExample>(serialized);

            Assert.Equal(deserialized, data);
        }
    }
}