using Apocryph.Ipfs;

namespace Apocryph.Consensus.Snowball.FunctionApp
{
    public class Query
    {
        public Hash<Block>? Value { get; private set; }
        public int Round { get; private set; }

        public Query(Hash<Block>? value, int round)
        {
            Value = value;
            Round = round;
        }
    }
}