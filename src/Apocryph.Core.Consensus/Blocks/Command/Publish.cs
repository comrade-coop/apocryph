namespace Apocryph.Core.Consensus.Blocks.Command
{
    public class Publish
    {
        public (string, byte[]) Message { get; }

        public Publish((string, byte[]) message)
        {
            Message = message;
        }
    }
}