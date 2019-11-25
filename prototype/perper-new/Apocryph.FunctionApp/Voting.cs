using System.Threading;
using System.Threading.Tasks;
using Apocryph.FunctionApp.Model;
using Microsoft.Azure.WebJobs;
using Perper.WebJobs.Extensions.Bindings;
using Perper.WebJobs.Extensions.Model;
using Perper.WebJobs.Extensions.Triggers;

namespace Apocryph.FunctionApp
{
    public class Voting
    {
        [FunctionName("Voting")]
        public static async Task Run([PerperStreamTrigger] IPerperStreamContext context,
            [PerperStream("voteStream")] IPerperStream<VoteMessage> voteStream,
            [PerperStream] IAsyncCollector<VoteMessage> outputStream)
        {
            await voteStream.Listen(async vote => await outputStream.AddAsync(vote),
                CancellationToken.None);
        }
    }
}