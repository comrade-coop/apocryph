using System.Collections.Generic;
using Apocryph.FunctionApp.Command;

namespace Apocryph.FunctionApp.Model
{
    public class CallNotification
    {
        public string From { get; set; }
        // TODO: Use Merkle proofs to show that the command is part of the Step
        public SendMessageCommand Command { get; set; }
        public Hash Step { get; set; }
        public Hash ValidatorSet { get; set; }
        public List<ISigned<Commit>> Commits { get; set; }
    }
}