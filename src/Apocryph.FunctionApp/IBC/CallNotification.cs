using System.Collections.Generic;
using Apocryph.FunctionApp.Model;
using Ipfs;

namespace Apocryph.FunctionApp.IBC
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