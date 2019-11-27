namespace Apocryph.FunctionApp.Model
{
    public class Vote
    {
        public AgentInput Input { get; set; }
        public AgentOutput Output { get; set; }

        public string Signer { get; set; }
    }
}