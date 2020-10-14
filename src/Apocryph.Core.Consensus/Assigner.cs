using System;
using System.Linq;
using System.Collections.Generic;
using System.Numerics;
using System.Security.Cryptography;
using Apocryph.Core.Consensus.VirtualNodes;

namespace Apocryph.Core.Consensus
{
    public class Assigner
    {
        private struct SlotOccupant
        {
            public BigInteger Difficulty { get; }
            public BigInteger Position { get => Difficulty; }
            public Peer Peer { get; }
            public byte[] Proof { get; }

            public SlotOccupant(BigInteger difficulty, Peer peer, byte[] proof)
            {
                Difficulty = difficulty;
                Peer = peer;
                Proof = proof;
            }

            public SlotOccupant(Peer peer, byte[] proof)
            {
                using var sha256Hash = SHA256.Create();
                var hash = sha256Hash.ComputeHash(peer.Value.Concat(new byte[] { 0 }).Concat(proof).ToArray());
                Peer = peer;
                Proof = proof;
                Difficulty =  new BigInteger(hash.Concat(new byte[] { 0 }).ToArray());
            }
        }

        private class Slot
        {
            public byte[] Salt { get; set; } = new byte[0];
            public SlotOccupant? Occupant { get; set; }
        }

        private Dictionary<Guid, Slot[]> _slots = new Dictionary<Guid, Slot[]>();

        public event Action<Node, Peer, byte[]>? SlotOccupantChanged;

        public void SetSalt(Guid chainId, int slot, byte[] salt)
        {
            _slots[chainId][slot].Salt = salt;
        }

        public void AddChain(Guid chainId, int slotCount)
        {
            _slots[chainId] = Enumerable.Range(0, slotCount).Select(x => new Slot()).ToArray();
        }

        public bool ProcessClaim(SlotClaim claim)
        {
            return ProcessClaim(claim.ChainId, claim.Peer, claim.Proof);
        }

        public void ProcessClaim(Peer peer, byte[] proof)
        {
            var occupant = new SlotOccupant(peer, proof);
            foreach (var (chainId, _) in _slots)
            {
                ProcessClaim(chainId, occupant);
            }
        }

        public bool ProcessClaim(Guid chainId, Peer peer, byte[] proof)
        {
            return ProcessClaim(chainId, new SlotOccupant(peer, proof));
        }

        // TODO: Ensure thread safety
        private bool ProcessClaim(Guid chainId, SlotOccupant occupant)
        {
            var slotIndex = (int)(occupant.Position % _slots[chainId].Length);
            var slot = _slots[chainId][slotIndex];

            if (slot.Occupant is SlotOccupant existingOccupant
                ) // && existingOccupant.Difficulty > occupant.Difficulty)
            {
                return false;
            }

            slot.Occupant = occupant;
            SlotOccupantChanged?.Invoke(new Node(chainId, slotIndex), occupant.Peer, occupant.Proof);

            return true;
        }

        public Node?[] GetNodes(Guid chainId)
        {
            return Enumerable.Range(0, _slots[chainId].Length).Select(x => new Node(chainId, x)).ToArray();
        }

        public Dictionary<Guid, Node?[]> GetNodes()
        {
            return _slots.Keys.ToDictionary(chainId => chainId, chainId => GetNodes(chainId));
        }
    }
}