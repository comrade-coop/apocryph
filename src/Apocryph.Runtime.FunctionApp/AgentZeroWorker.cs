using System;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;
using Apocryph.Core.Consensus;

namespace Apocryph.Runtime.FunctionApp
{
    public class AgentZeroWorker
    {
        [FunctionName(nameof(AgentZeroWorker))]
        [return: Perper("$return")]
        public Task<(byte[]?, (string, object[])[], Dictionary<Guid, string[]>, Dictionary<Guid, string>)> Run([PerperWorkerTrigger] PerperWorkerContext context,
            [Perper("input")] (byte[]?, (string, byte[]), Guid?) input, CancellationToken cancellationToken)
        {
            var (serializedState, (messageType, serializedMessage), sender) = input;

            var state = JsonSerializer.Deserialize<AgentZeroState>(serializedState);

            var message = JsonSerializer.Deserialize(serializedMessage, Type.GetType(messageType));

            state = AgentZero.Run(state, message, sender);

            return Task.FromResult(((byte[]?)JsonSerializer.SerializeToUtf8Bytes(state), new (string, object[])[0], new Dictionary<Guid, string[]>(), new Dictionary<Guid, string>())!);
        }
    }
}