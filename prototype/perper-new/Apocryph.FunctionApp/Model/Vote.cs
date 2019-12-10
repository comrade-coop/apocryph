namespace Apocryph.FunctionApp.Model
{
    public class Vote
    {
        public IAgentStep For { get; set; } // FIXME: Should be hash stored on IPFS

        public ValidatorKey Signer { get; set; }
        public ValidatorSignature Signature { get; set; }
    }
}