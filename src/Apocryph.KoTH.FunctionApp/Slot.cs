using Apocryph.Peering;

namespace Apocryph.KoTH
{
    public class Slot
    {
        public Peer Peer { get; set; }
        public byte[] MinedData { get; set; }

        public Slot(Peer peer, byte[] minedData)
        {
            Peer = peer;
            MinedData = minedData;
        }
    }
}