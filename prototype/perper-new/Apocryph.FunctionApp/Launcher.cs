using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class Launcher
    {
        [FunctionName("Launcher")]
        public static async Task Run([Perper(Stream = "Launcher")] IPerperStreamContext context)
        {
            var generator = await context.CallStreamFunction("Generator", new {count = 100});
            var processor = await context.CallStreamFunction("Processor", new {generator, multiplier = 10});
            await context.CallStreamAction("Consumer", new {processor});
        }
    }
}