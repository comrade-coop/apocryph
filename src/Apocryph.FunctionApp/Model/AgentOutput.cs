using System.Collections.Generic;
using Apocryph.FunctionApp.Command;

namespace Apocryph.FunctionApp.Model
{
    public class AgentOutput : IAgentStep
    {
        public object State { get; set; }
        public List<ICommand> Commands { get; set; }

        public Hash Previous { get; set; }
        public List<ISigned<Commit>> PreviousCommits { get; set; }
    }
}