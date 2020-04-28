namespace Apocryph.Agents.Testbed.Api
{
    public class AgentCommands
    {
        public string Origin;

        public object State { get; set; }
        public AgentCommand[] Commands { get; set; }
    }
}