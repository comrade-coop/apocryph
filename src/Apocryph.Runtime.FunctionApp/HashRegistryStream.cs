using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Apocryph.Core.Consensus.Blocks;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Config;
using Perper.WebJobs.Extensions.Model;

namespace Apocryph.Runtime.FunctionApp
{
    public class HashRegistryStream
    {
        [FunctionName(nameof(HashRegistryStream))]
        public async Task Run([PerperStreamTrigger] PerperStreamContext context,
            [Perper("filter")] Type filter,
            [Perper("input")] IAsyncEnumerable<object> input,
            [Perper("output")] IAsyncCollector<HashRegistryEntry> output)
        {
            await foreach (var value in input)
            {
                if (filter.IsAssignableFrom(value.GetType()))
                {
                    var hash = Hash.From(value);

                    await output.AddAsync(new HashRegistryEntry(hash.ToString(), value));
                }
            }
        }

        public static T? GetObjectByHash<T>(IQueryable<HashRegistryEntry> queryable, Hash hash)
            where T: class
        {
            var s = hash.ToString();
            return (T?)queryable.Where(x => x.Hash == s).FirstOrDefault()?.Value;
        }
    }
}