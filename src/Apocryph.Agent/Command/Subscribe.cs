namespace Apocryph.Agent.Command
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