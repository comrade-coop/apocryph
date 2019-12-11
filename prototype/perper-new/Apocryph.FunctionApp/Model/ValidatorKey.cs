using System.Security.Cryptography;
using System.Text;
using Newtonsoft.Json;

namespace Apocryph.FunctionApp.Model
{
    public struct ValidatorKey
    {
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