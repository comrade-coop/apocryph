namespace Apocryph.AgentZero.Messages
{
    public class SetAgentBlockMessage
    {
        public string AgentId { get; set; }
        public byte[] BlockId { get; set; }
    }
}