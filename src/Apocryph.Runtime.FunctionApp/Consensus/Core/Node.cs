namespace Apocryph.Runtime.FunctionApp.Consensus.Core
{
    public class Node
    {
        public string PublicKey { get; set; }
        public int Stake { get; set; }
        public bool IsProposer { get; set; }
    }
}