namespace Apocryph.FunctionApp.Model.Command
{
    public class IssueTokenCommand : ICommand
    {
        public string Nonce { get; set; }
    }
}