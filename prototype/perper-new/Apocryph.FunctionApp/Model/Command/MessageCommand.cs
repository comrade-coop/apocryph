using Apocryph.FunctionApp.Model.Message;

namespace Apocryph.FunctionApp.Model.Command
{
    public class MessageCommand : ICommand
    {
        public string Target { get; set; }
        public string TokenNonce { get; set; }
        public IMessage Message { get; set; }
    }
}