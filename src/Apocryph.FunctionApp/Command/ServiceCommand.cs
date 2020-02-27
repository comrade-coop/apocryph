namespace Apocryph.FunctionApp.Command
{
    public class ServiceCommand : ICommand
    {
        public string Service { get; set; }
        public object Parameters { get; set; }
    }
}