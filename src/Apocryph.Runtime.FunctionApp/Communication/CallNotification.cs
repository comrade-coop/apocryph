using System.Collections.Generic;
using Apocryph.Agent;
using Apocryph.Runtime.FunctionApp.Ipfs;
using Ipfs;

namespace Apocryph.Runtime.FunctionApp.Communication
{
    public class CallNotification
    {
        public string From { get; set; }
        // TODO: Use Merkle proofs to show that the command is part of the Block
        public AgentCommand Command { get; set; }
        public Cid Block { get; set; }
        public Cid ValidatorSet { get; set; }
        public List<ISigned<Cid>> Signatures { get; set; }
    }
}