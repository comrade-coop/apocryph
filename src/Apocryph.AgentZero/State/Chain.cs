namespace Apocryph.AgentZero.State
{
    public class Chain
    {
        public byte[] LatestBlock { get; set; }

        public Chain(byte[] latestBlock)
        {
            LatestBlock = latestBlock;
        }
    }
}