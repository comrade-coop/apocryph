namespace Apocryph.Agent.Command
{
    // Delete and use Invoke instead to AgentZero
    public class Create
    {
        public string Agent { get; }
        public (string, byte[]) Message { get; }

        public Create(string agent, (string, byte[]) message)
        {
            Agent = agent;
            Message = message;
        }
    }
}