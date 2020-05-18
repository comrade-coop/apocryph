namespace Apocryph.Runtime.FunctionApp.ValidatorSelection
{
    public struct SlotClaim
    {
        public PublicKey Key { get; set; } // Should also be signed with this key
        public byte[] AgentId { get; set; }
    }
}