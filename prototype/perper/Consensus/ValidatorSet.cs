using System.Collections.Generic;

namespace Apocryph.Consensus
{
    public class ValidatorSet : HashSet<Validator>, ISet<Validator>
    {
        public Validator Proposer;
    }
}
