using System.Linq;
using System.Numerics;
using Apocryph.Ipfs;

namespace Apocryph.KoTH
{
    public class KoTHState
    {
        public Slot?[] Slots { get; }

        public KoTHState(Slot?[] slots)
        {
            Slots = slots;
        }

        public static BigInteger GetDifficulty(Slot slot)
        {
            var hash = Hash.From(slot);
            return new BigInteger(hash.Bytes.Concat(new byte[] { 0 }).ToArray());
        }

        public bool TryInsert(Slot newSlot)
        {
            var newDifficulty = GetDifficulty(newSlot);
            var index = (int)(newDifficulty % Slots.Length);

            var currentSlot = Slots[index];
            if (currentSlot == null || GetDifficulty(currentSlot) < newDifficulty)
            {
                Slots[index] = newSlot;
                return true;
            }

            return false;
        }
    }
}