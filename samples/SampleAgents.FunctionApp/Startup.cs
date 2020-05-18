using System;
using Apocryph.Testbed;
using Microsoft.Azure.Functions.Extensions.DependencyInjection;
using Microsoft.Extensions.DependencyInjection;

[assembly: FunctionsStartup(typeof(SampleAgents.FunctionApp.Startup))]

namespace SampleAgents.FunctionApp
{
    [Obsolete]
    public class Startup : FunctionsStartup
    {
        public override void Configure(IFunctionsHostBuilder builder)
        {
            builder.Services.AddLogging();
            builder.Services.AddTransient(typeof(Testbed), typeof(Testbed));
        }
    }
}