using System.Collections.Generic;
using Apocryph.FunctionApp.Command;

namespace Apocryph.FunctionApp.Model
{
    public class AgentOutput : IAgentStep
    {
        public string Type { get; set; }
        
        public object State { get; set; }
        public IEnumerable<ICommand> Commands { get; set; }
    }
}