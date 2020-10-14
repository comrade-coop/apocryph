using System;

namespace Apocryph.Core.Consensus.VirtualNodes
{
    public struct SlotClaim
    {
        public Guid ChainId { get; set; }
        public Peer Peer { get; set; }
        public byte[] Proof { get; set; }

        public SlotClaim(Guid chainId, Peer peer, byte[] proof)
        {
            ChainId = chainId;
            Peer = peer;
            Proof = proof;
        }
    }
}