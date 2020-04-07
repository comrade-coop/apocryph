using Apocryph.Agents.Testbed;
using Microsoft.Azure.Functions.Extensions.DependencyInjection;
using Microsoft.Extensions.DependencyInjection;

[assembly: FunctionsStartup(typeof(Apocryph.Agent.FunctionApp.Startup))]

namespace Apocryph.Agent.FunctionApp
{
    public class Startup : FunctionsStartup
    {
        public override void Configure(IFunctionsHostBuilder builder)
        {
            builder.Services.AddLogging();
            builder.Services.AddTransient(typeof(Testbed), typeof(Testbed));
        }
    }
}