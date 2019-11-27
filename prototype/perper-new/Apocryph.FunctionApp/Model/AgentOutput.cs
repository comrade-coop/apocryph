using System.Collections.Generic;
using Apocryph.FunctionApp.Model.Command;

namespace Apocryph.FunctionApp.Model
{
    public class AgentOutput
    {
        public object State { get; set; }
        public IEnumerable<ICommand> Output { get; set; }
        public string Type { get; set; }
    }
}