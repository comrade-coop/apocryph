namespace Apocryph.Core.Consensus.VirtualNodes
{
    public struct SlotClaim
    {
        public PublicKey Key { get; set; } // Should also be signed with this key
        public byte[] ChainId { get; set; }
    }
}