using System;
using System.Collections.Generic;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using System.Security.Cryptography;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace TestHarness.FunctionApp
{
    public class ChainAgentPong
    {
        [FunctionName(nameof(ChainAgentPong))]
        [return: Perper("$return")]
        public Task<(byte[]?, (string, object[])[], Dictionary<Guid, string[]>, Dictionary<Guid, string>)> Run([PerperWorkerTrigger] PerperWorkerContext context,
            [Perper("input")] (byte[]?, (string, byte[]), Guid?) input, CancellationToken cancellationToken)
        {
            var (state, (_, _), ownReference) = input;
            var otherReference = JsonSerializer.Deserialize<ChainAgentState>(state).OtherReference;
            var actions = new List<(string, object[])>
            {
                ("Invoke", new object[]
                {
                    otherReference,
                    (typeof(string).FullName!, JsonSerializer.SerializeToUtf8Bytes("Pong"))
                })
            };
            using var sha1 = new SHA1CryptoServiceProvider();
            var carrier = Convert.ToBase64String(sha1.ComputeHash(state));
            return Task.FromResult((state, actions.ToArray(), new Dictionary<Guid, string[]>(), new Dictionary<Guid, string>() { { otherReference, carrier }, { ownReference!.Value, carrier } }));
        }
    }
}