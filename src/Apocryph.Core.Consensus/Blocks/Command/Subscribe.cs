namespace Apocryph.Core.Consensus.Blocks.Command
{
    public class Subscribe
    {
        public string Target { get; }

        public Subscribe(string target)
        {
            Target = target;
        }
    }
}