using System;
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

        private byte[] _chainId;
        private Slot[] _slots;

        private Action<int, PrivateKey> _addPrivateKey;
        private Action<int, PrivateKey> _removePrivateKey;

        public Assigner(int slotCount, byte[] chainId, Action<int, PrivateKey> addPrivateKey, Action<int, PrivateKey> removePrivateKey)
        {
            _slots = new Slot[slotCount];
            _chainId = chainId;
            _addPrivateKey = addPrivateKey;
            _removePrivateKey = removePrivateKey;
        }

        public void SetSalt(int slot, byte[] salt)
        {
            _slots[slot].Salt = salt;
        }

        public bool AddKey(PublicKey key, PrivateKey? privateKey)
        {
            var slotIndex = GetSlotForKey(key);
            var slot = _slots[slotIndex];

            if (slot.PublicKey is PublicKey slotKey
                && slotKey.GetDifficulty(_chainId, slot.Salt) > key.GetDifficulty(_chainId, slot.Salt))
            {
                return false;
            }

            if (slot.PrivateKey != null)
            {
                _removePrivateKey.Invoke(slotIndex, slot.PrivateKey.Value);
            }

            slot.PublicKey = key;
            slot.PrivateKey = privateKey;

            if (slot.PrivateKey != null)
            {
                _addPrivateKey.Invoke(slotIndex, slot.PrivateKey.Value);
            }

            return true;
        }

        private int GetSlotForKey(PublicKey key)
        {
            return (int)(key.GetPosition() % _slots.Length);
        }
    }
}