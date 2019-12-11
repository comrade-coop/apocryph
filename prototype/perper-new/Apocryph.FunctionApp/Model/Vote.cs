namespace Apocryph.FunctionApp.Model
{
    public class Vote
    {
        public Hash ForHash { get; set; }

        public ValidatorKey Signer { get; set; }
        public ValidatorSignature Signature { get; set; }
    }
}