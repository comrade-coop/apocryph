using System;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace TestHarness.FunctionApp
{
    public class ChainAgentPong
    {
        [FunctionName(nameof(ChainAgentPong))]
        public Task<(byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>)> Run([PerperWorkerTrigger] PerperWorkerContext context,
            [Perper("input")] (byte[]?, (string, byte[]), Guid?) input, CancellationToken cancellationToken)
        {
            var (state, (_, _), _) = input;
            var actions = new List<(string, object[])>
            {
                ("", new object[]
                {
                    JsonSerializer.Deserialize<ChainAgentState>(state).OtherReference,
                    (typeof(string).FullName!, JsonSerializer.SerializeToUtf8Bytes("Pong"))
                })
            };
            return Task.FromResult<(byte[]?, (string, object[])[], IDictionary<Guid, string[]>, IDictionary<Guid, string>)>((state, actions.ToArray(), null, null)!);
        }
    }
}