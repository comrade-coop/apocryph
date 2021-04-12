using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using Apocryph.Ipfs;

namespace Apocryph.Consensus.Snowball.FunctionApp
{
    public class SnowballState
    {
        public Hash<Block>? CurrentValue { get; protected set; } = null;
        public Hash<Block>? LastValue { get; protected set; } = null;
        public Dictionary<Hash<Block>, int> Memory { get; protected set; } = new Dictionary<Hash<Block>, int>();
        public int Confidence { get; protected set; } = 0;

        public void ProcessQuery(Hash<Block> requestSuggestion)
        {
            if (CurrentValue == null)
            {
                CurrentValue = requestSuggestion;
                Confidence = 0;
            }
        }

        public IEnumerable<Peer> SamplePeers(SnowballParameters parameters, Peer[] peers)
        {
            return peers.OrderBy(_ => RandomNumberGenerator.GetInt32(peers.Length)).Take(parameters.K);
        }

        public bool ProcessResponses(SnowballParameters parameters, IEnumerable<Hash<Block>> responses)
        {
            if (Confidence > parameters.Beta) return true;

            var threshold = parameters.Alpha * parameters.K;

            var responseValues = responses.GroupBy(response => response)
                .Where(group => group.Count() > threshold)
                .Select(group => group.Key);

            foreach (var responseValue in responseValues)
            {
                Memory[responseValue] = Memory.TryGetValue(responseValue, out var answerCount) ? answerCount + 1 : 1;

                if (CurrentValue == null)
                {
                    CurrentValue = responseValue;
                }

                if (Memory[responseValue] > Memory[CurrentValue])
                {
                    CurrentValue = responseValue;
                }

                if (responseValue != LastValue)
                {
                    Confidence = 0;
                    LastValue = responseValue;
                }
                else
                {
                    if (Confidence++ > parameters.Beta) return true;
                }
            }

            return false;
        }
    }
}