using System;
using System.Security.Cryptography;
using System.Text;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public struct ValidatorKey
    {
        public class ECParametersJsonConverter : JsonConverter<ECParameters>
        {
            public ECCurve Curve { get; set; } = ECCurve.NamedCurves.nistP521;

            public override void WriteJson(JsonWriter writer, ECParameters value, JsonSerializer serializer)
            {
                if (value.Curve.Equals(Curve))
                {
                    throw new Exception("Unexpected curve value");
                }

                serializer.Serialize(writer, value.Q, typeof(ECPoint));
            }

            public override ECParameters ReadJson(JsonReader reader, Type type, ECParameters _existingValue, bool what, JsonSerializer serializer)
            {
                var q = serializer.Deserialize<ECPoint>(reader);
                return new ECParameters
                {
                    Q = q,
                    Curve = Curve,
                };
            }
        }

        [JsonConverter(typeof(ECParametersJsonConverter))]
        public ECParameters Key { get; set; }

        public bool ValidateSignature(Hash hash, ValidatorSignature signature)
        {
            using (var ecdsa = ECDsa.Create(Key))
            {
                return ecdsa.VerifyHash(hash.Bytes, signature.Bytes);
            }
        }
    }
}