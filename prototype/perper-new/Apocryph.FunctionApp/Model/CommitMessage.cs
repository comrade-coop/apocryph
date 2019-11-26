namespace Apocryph.FunctionApp.Model
{
    public class CommitMessage
    {
        public AgentInput Input { get; set; }
        public AgentOutput Output { get; set; }

        public string Signer { get; set; }
    }
}