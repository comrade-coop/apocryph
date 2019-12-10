namespace Apocryph.FunctionApp.Model
{
    public class Commit
    {
        public IAgentStep For { get; set; }

        public ValidatorKey Signer { get; set; }
        public ValidatorSignature Signature { get; set; }
    }
}