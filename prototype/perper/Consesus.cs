using Microsoft.Azure.WebJobs;
using Mock.Perper;

public static class Consesus
{
    [FunctionName("Consesus")]
    public static void Run(PerperStreamContext<object> context)
    {
        //TODO: @Bob - split into more functions / classes.
        var executionStream = context.CallStreamFunction("Execution");

        var generatorStream = context.CallStreamFunction("Generator");
        var validatorStream = context.CallStreamFunction("Execution", generatorStream);
    }
}