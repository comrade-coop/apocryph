using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Agent;
using Apocryph.FunctionApp.AgentZero.Publications;
using Apocryph.FunctionApp.Command;
using Apocryph.FunctionApp.Ipfs;
using Apocryph.FunctionApp.Model;
﻿using Ipfs;
﻿using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Microsoft.Extensions.Logging;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class AgentZetoStepOrderVerifier
    {
        private class State
        {
            public Hash CurrentStep { get; set; } = new Hash {Bytes = new byte[] {}};
            public Dictionary<String, IHashed<ValidatorSet>> CurrentValidatorSets { get; set; }
        }

        [FunctionName(nameof(AgentZetoStepOrderVerifier))]
        public static async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [PerperStream("stepsStream")] IAsyncEnumerable<IHashed<IAgentStep>> stepsStream,
            [Perper("genesisValidatorSet")] ValidatorSet genesisValidatorSet,
            [Perper("ipfsGateway")] string ipfsGateway,
            [PerperStream("outputStream")] IAsyncCollector<Dictionary<String, IHashed<ValidatorSet>>> outputStream,
            ILogger logger)
        {
            var agentId = "0";
            var ipfs = new IpfsClient(ipfsGateway);

            var genesisJToken = IpfsJsonSettings.JTokenFromObject(genesisValidatorSet);
            var genesisCid = await ipfs.Dag.PutAsync(genesisJToken, cancel: CancellationToken.None);
            var genesisHash = new Hash {Bytes = genesisCid.ToArray()};

            var state = await context.FetchStateAsync<State>() ?? new State()
            {
                CurrentValidatorSets = new Dictionary<String, IHashed<ValidatorSet>>()
                {
                    [agentId] = Hashed.Create(genesisValidatorSet, genesisHash)
                }
            };

            await outputStream.AddAsync(state.CurrentValidatorSets);

            await stepsStream.ForEachAsync(async step =>
            {
                try
                {
                    if (step.Value.Previous != state.CurrentStep)
                    {
                        return;
                    }

                    if (step.Value.PreviousValidatorSet != state.CurrentValidatorSets[agentId])
                    {
                        return;
                    }

                    if (step is AgentOutput output)
                    {
                        foreach (var command in output.Commands)
                        {
                            if (command is PublicationCommand publication && publication.Payload is ValidatorSetPublication change)
                            {
                                var validatorSet = new ValidatorSet();
                                foreach (var (key, weight) in change.Weights)
                                {
                                    var validatorKey = IpfsJsonSettings.BytesToObject<ValidatorKey>(Encoding.UTF8.GetBytes(key));
                                    validatorSet.Weights[validatorKey] = (int) weight;
                                }
                                var jToken = IpfsJsonSettings.JTokenFromObject(validatorSet);
                                var cid = await ipfs.Dag.PutAsync(jToken, cancel: CancellationToken.None);
                                var hash = new Hash {Bytes = cid.ToArray()};

                                state.CurrentValidatorSets[change.AgentId] = Hashed.Create(validatorSet, hash);
                            }
                        }
                    }

                    state.CurrentStep = step.Hash;
                    await outputStream.AddAsync(state.CurrentValidatorSets);
                    await context.UpdateStateAsync(state);
                }
                catch (Exception e)
                {
                    logger.LogError(e.ToString());
                }
            }, CancellationToken.None);
        }
    }
}