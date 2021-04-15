using System;
using System.Collections.Generic;
using Apocryph.Ipfs.Impl;
using Ipfs.Http;
using Microsoft.Azure.WebJobs;
using Microsoft.Azure.WebJobs.Description;
using Microsoft.Azure.WebJobs.Host.Bindings;
using Microsoft.Azure.WebJobs.Host.Config;
using Microsoft.Azure.WebJobs.Hosting;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Options;

[assembly: WebJobsStartup(typeof(Apocryph.Ipfs.Config.IpfsWebJobsStartup))]

namespace Apocryph.Ipfs.Config
{
    [Extension("Ipfs")]
    public class IpfsExtensionConfigProvider : IExtensionConfigProvider
    {
        public void Initialize(ExtensionConfigContext context)
        {
        }
    }

    public class IpfsWebJobsStartup : IWebJobsStartup
    {
        public void Configure(IWebJobsBuilder builder)
        {
            builder.AddExtension<IpfsExtensionConfigProvider>().BindOptions<IpfsConfig>();

            builder.Services.AddSingleton(typeof(IHashResolver), typeof(IpfsHashResolver));
            builder.Services.AddSingleton(typeof(IPeerConnector), typeof(IpfsPeerConnector));
            builder.Services.AddSingleton<IBindingProvider>(services => new ServiceBindingProvider(new HashSet<Type>
            {
                typeof(IHashResolver),
                typeof(IPeerConnector)
            }, services));

            builder.Services.AddSingleton<IpfsClient>(services =>
            {
                var config = services.GetRequiredService<IOptions<IpfsConfig>>().Value;

                return new IpfsClient(config.IpfsApiEndpoint);
            });
        }
    }
}