using Apocryph.Consensus;

namespace SampleAgents.FunctionApp.Agents
{
    public class PingPongMessage
    {
        public Reference Callback { get; set; }
        public string Content { get; set; }
        public int AccumulatedValue { get; set; }

        public PingPongMessage(Reference callback, string content, int accumulatedValue)
        {
            Callback = callback;
            Content = content;
            AccumulatedValue = accumulatedValue;
        }
    }
}