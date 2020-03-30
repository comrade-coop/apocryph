namespace Apocryph.FunctionApp.Command
{
    public class SendMessageCommand : ICommand
    {
        public string Target { get; set; }
        public object Payload { get; set; }
    }
}