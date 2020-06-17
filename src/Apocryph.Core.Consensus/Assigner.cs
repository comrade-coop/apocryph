using System;
using System.Linq;
using System.Collections.Generic;
using Apocryph.Core.Consensus.VirtualNodes;

namespace Apocryph.Core.Consensus
{
    public class Assigner
    {
        private class Slot
        {
            public byte[] Salt { get; set; } = new byte[0];
            public Node? Occupant { get; set; }
            public PublicKey? PublicKey { get; set; }
            public PrivateKey? PrivateKey { get; set; }
        }

        private Dictionary<byte[], Slot[]> _slots = new Dictionary<byte[], Slot[]>();

        private Func<byte[], int, PublicKey, PrivateKey?, Node> _createNode;

        public Assigner(Func<byte[], int, PublicKey, PrivateKey?, Node> createNode)
        {
            _createNode = createNode;
        }

        public void SetSalt(byte[] chainId, int slot, byte[] salt)
        {
            _slots[chainId][slot].Salt = salt;
        }

        public void AddChain(byte[] chainId, int slotCount)
        {
            _slots[chainId] = new Slot[slotCount];
        }

        public void AddKey(PublicKey key, PrivateKey? privateKey)
        {
            foreach (var (chainId, _) in _slots)
            {
                AddKey(chainId, key, privateKey);
            }
        }

        public bool AddKey(byte[] chainId, PublicKey key, PrivateKey? privateKey)
        {
            var slotIndex = (int)(key.GetPosition() % _slots[chainId].Length);
            var slot = _slots[chainId][slotIndex];

            if (slot.PublicKey is PublicKey slotKey
                && slotKey.GetDifficulty(chainId, slot.Salt) > key.GetDifficulty(chainId, slot.Salt))
            {
                return false;
            }

            slot.PublicKey = key;
            slot.PrivateKey = privateKey;

            slot.Occupant = _createNode.Invoke(chainId, slotIndex, slot.PublicKey.Value, slot.PrivateKey);

            return true;
        }

        public Node?[] GetNodes(byte[] chainId)
        {
            return _slots[chainId].Select(x => x.Occupant).ToArray();
        }

        public Dictionary<byte[], Node?[]> GetNodes()
        {
            return _slots.ToDictionary(kv => kv.Key, kv => kv.Value.Select(x => x.Occupant).ToArray());
        }
    }
}