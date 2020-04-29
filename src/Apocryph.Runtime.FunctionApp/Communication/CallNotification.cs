using System.Collections.Generic;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Ipfs;

namespace Apocryph.Runtime.FunctionApp.Commuication
{
    public class CallNotification
    {
        public string From { get; set; }
        // TODO: Use Merkle proofs to show that the command is part of the Step
        public SendMessageCommand Command { get; set; }
        public Cid Step { get; set; }
        public Cid ValidatorSet { get; set; }
        public List<ISigned<Commit>> Commits { get; set; }
    }
}