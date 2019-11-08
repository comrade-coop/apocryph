using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Akka.Actor;

namespace Apocryph.Prototype
{
    public class VotingRound
    {
        private readonly ISet<CoreService.Vote> _votes = new HashSet<CoreService.Vote>();
        private readonly Dictionary<string, int> _validatorStakes = new Dictionary<string, int>();
        private bool _finished = false;

        public VotingRound(Dictionary<string, int> validatorStakes)
        {
            _validatorStakes = validatorStakes;
        }

        public void AddVote(CoreService.Vote vote)
        {
            _votes.Add(vote);
        }

        public bool HasEnoughVotes()
        {
            if (_finished) return false;

            var necessary = _validatorStakes.Values.Sum() * 2 / 3;
            var present = _votes.Select(x => _validatorStakes[x.Signer]).Sum();

            return present >= necessary;
        }

        public CoreService.Chunk GetOptimalChunk()
        {
            var necessary = _validatorStakes.Values.Sum() * 2 / 3;
            var optimalChunk = _votes.Select(x => x.LastChunk).Aggregate((x, n) => (x != null ? x.Height : 0) > (n != null ? n.Height : 0) ? x : n);
            var present = 0;

            while (optimalChunk != null) {
                present += _votes.Where(x => x.LastChunk == optimalChunk)
                    .Select(x => _validatorStakes[x.Signer]).Sum();
                if (present >= necessary) {
                    break;
                }
                optimalChunk = optimalChunk.Previous;
            }

            return optimalChunk;
        }

        public ISet<CoreService.Vote> GetVotes()
        {
            return _votes;
        }

        public void Finish()
        {
            _finished = true;
        }
    }
}
