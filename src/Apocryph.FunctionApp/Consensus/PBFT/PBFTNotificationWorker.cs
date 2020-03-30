using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Ipfs;
using Apocryph.FunctionApp.Model;
﻿using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.FunctionApp
{
    public static class PBFTNotificationWorker
    {
        [FunctionName(nameof(PBFTNotificationWorker))]
        [return: Perper("$return")]
        public static async Task<object> Run([PerperWorkerTrigger] PerperWorkerContext context,
            [Perper("agentId")] string agentId,
            [Perper("validatorSet")] IHashed<ValidatorSet> validatorSet,
            [Perper("notification")] ISigned<object> notification,
            [Perper("ipfsGateway")] string ipfsGateway)
        {
            var topic = "apocryph-agentNotifications-" + agentId;

            var ipfs = new IpfsClient(ipfsGateway);

            var bytes = IpfsJsonSettings.ObjectToBytes(notification);

            await ipfs.PubSub.PublishAsync(topic, bytes, CancellationToken.None);

            return null;
        }
    }
}