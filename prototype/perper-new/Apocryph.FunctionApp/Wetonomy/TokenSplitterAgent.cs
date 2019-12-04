using System.Collections.Generic;
using System.Linq;
using System.Numerics;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Wetonomy.Messages;

namespace Apocryph.FunctionApp.Wetonomy
{
    public static class TokenSplitterAgent
    {
        public class State
        {
            public List<string> Targets { get; set; }
        }
        
        public static async Task<(object, IEnumerable<object>)> Run(object _state, object action)
        {
            var state = _state as State ?? new State();
            var messages = new List<object>();
            switch (action)
            {
                case TokensChangedPublication tokensChangedEvent:
                    if (tokensChangedEvent.Change > 0)
                    {
                        var individualAmount = tokensChangedEvent.Total / state.Targets.Count;
                        messages.Add(state.Targets.Select(x => new TransferMessage
                        {
                            Amount = individualAmount,
                            From = tokensChangedEvent.Target,
                            To = x
                        }));
                    }
                    break;
            }

            await Task.CompletedTask;

            return (state, messages);
        }
        
    }
}