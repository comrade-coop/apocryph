namespace Apocryph.Consensus.Snowball
{
    public class SnowballParameters
    {
        public int K { get; set; } // Peer sample size
        public double Alpha { get; set; } // Sample threshold
        public int Beta { get; set; } // Required confidence

        public SnowballParameters(int k, double alpha, int beta)
        {
            K = k;
            Alpha = alpha;
            Beta = beta;
        }
    }
}