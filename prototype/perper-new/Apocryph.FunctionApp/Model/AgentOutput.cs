using System.Collections.Generic;

namespace Apocryph.FunctionApp.Model
{
    public class AgentOutput
    {
        public object State { get; set; }
        public IEnumerable<object> Output { get; set; }
        public string Type { get; set; }
    }
}