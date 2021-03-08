using Apocryph.HashRegistry;

namespace Apocryph.Peering
{
    public class Peer
    {
        public Hash<object> Id { get; }

        public Peer(Hash<object> id)
        {
            Id = id;
        }
    }
}