using System;
using System.Linq;

namespace Apocryph.Ipfs.Test
{
    public interface IExample
    {
    }

    public class Example1 : IExample
    {
        public int Number { get; set; }

        public override bool Equals(object? obj) => obj is Example1 other && Number == other.Number;
        public override int GetHashCode() => Number.GetHashCode();
        public override string ToString() => $"Example1({Number})";
    }

    public static class TestData
    {
        public class Example2 : IExample // NOTE: Nested in order to test nested type serialization
        {
            public string String { get; set; } = "default";

            public override bool Equals(object? obj) => obj is Example2 other && String == other.String;
            public override int GetHashCode() => String.GetHashCode();
            public override string ToString() => $"Example2({String})";
        }

        public static object?[][] Data = new[] {
            new object?[] { null },
            new object?[] { 42 },
            new object?[] { "Example text" },
            new object?[] { new Example1 { Number = 42 } },
            new object?[] { new Example2 { String = "Example text" } },
            new object?[] { (42, 103, "Tuple") },
            new object?[] { new byte[] { 0, 1, 2, 255, 254, 253, 42 } },
        };

        public static object?[][] DataInterface = Data.Where(x => x[0] is IExample).ToArray();
        public static object?[][] DataNoNull = Data.Where(x => x[0] != null).ToArray();
        public static object?[][] DataNoInterface = DataNoNull.Where(x => !(x[0] is IExample)).ToArray();
    }
}