using Apocryph.HashRegistry;

namespace Apocryph.Consensus.Snowball.FunctionApp
{
    public class Query
    {
        public Hash<Block> Value { get; }
        public int Round { get; }

        public Query(Hash<Block> value, int round)
        {
            Value = value;
            Round = round;
        }
    }
}