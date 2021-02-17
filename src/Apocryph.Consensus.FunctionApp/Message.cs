namespace Apocryph.Consensus
{
    public class Message
    {
        // public Reference Source { get; }
        public Reference Target { get; }
        public ReferenceData Data { get; }

        public Message(Reference target, ReferenceData data)
        {
            Target = target;
            Data = data;
        }
    }
}