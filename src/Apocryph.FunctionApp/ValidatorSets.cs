using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Apocryph.FunctionApp.AgentZero.Publications;

namespace Apocryph.FunctionApp
{
    public static class ValidatorSets
    {
        public class State
        {
            public Dictionary<string, ValidatorSet> ValidatorSets { get; set; }
        }

        [FunctionName("ValidatorSets")]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("stepVerifierStream")] IAsyncEnumerable<Hashed<AgentOutput>> stepsStream,
            [PerperStream("outputStream")] IAsyncCollector<Dictionary<string, ValidatorSet>> outputStream)
        {
            var state = await context.FetchStateAsync<State>();

            await stepsStream.ForEachAsync(async output =>
            {
                foreach (var command in output.Value.Commands)
                {
                    if (command is PublicationCommand publication)
                    {
                        if (publication.Payload is ValidatorSetPublication validatorSetPublication)
                        {
                            // TODO: Fix conversion
                            var validatorSet = new ValidatorSet {
                                Weights = validatorSetPublication.Weights.ToDictionary(
                                    kv => (ValidatorKey)(object)kv.Key,
                                    kv => (int)kv.Value),
                            };
                            state.ValidatorSets[validatorSetPublication.AgentId] = validatorSet;
                            await outputStream.AddAsync(state.ValidatorSets);
                        }
                    }
                }
            });
        }
    }
}