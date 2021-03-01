namespace Apocryph.Consensus
{
    public class Peer
    {
        public int Id { get; }
        public byte[] MinedData { get; }

        public Peer(int id, byte[] minedData)
        {
            Id = id;
            MinedData = minedData;
        }
    }
}