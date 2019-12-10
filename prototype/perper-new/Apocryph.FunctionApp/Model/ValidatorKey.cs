using System.Security.Cryptography;
using System.Text;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public struct ValidatorKey
    {
        public ECParameters Key { get; set; }

        public bool ValidateSignature(object item, ValidatorSignature signature)
        {
            var bytes = Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(item));

            using (var ecdsa = ECDsa.Create(Key))
            {
                return ecdsa.VerifyData(bytes, signature.Bytes, HashAlgorithmName.SHA256);
            }
        }
    }
}