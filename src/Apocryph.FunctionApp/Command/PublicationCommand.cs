namespace Apocryph.FunctionApp.Command
{
    public class PublicationCommand : ICommand
    {
        public object Payload { get; set; }
    }
}