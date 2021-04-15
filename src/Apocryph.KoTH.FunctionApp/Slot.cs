using Apocryph.Ipfs;

namespace Apocryph.KoTH
{
    public class Slot
    {
        public Peer Peer { get; private set; }
        public byte[] MinedData { get; private set; }

        public Slot(Peer peer, byte[] minedData)
        {
            Peer = peer;
            MinedData = minedData;
        }
    }
}