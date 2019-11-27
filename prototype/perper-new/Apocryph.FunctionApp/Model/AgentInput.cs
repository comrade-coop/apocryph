using Apocryph.FunctionApp.Model.Message;

namespace Apocryph.FunctionApp.Model
{
    public class AgentInput
    {
        public object State { get; set; }
        public IMessage Input { get; set; }
    }
}