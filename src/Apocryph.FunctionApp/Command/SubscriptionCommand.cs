namespace Apocryph.FunctionApp.Command
{
    public class SubscriptionCommand : ICommand
    {
        public string Target { get; set; }
    }
}