using System;
using System.Linq;
using System.Security.Cryptography;
using Ipfs;

namespace Apocryph.Runtime.FunctionApp.ValidatorSelection
{
    public struct ValidatorSlotClaim
    {
        public ValidatorKey Key { get; set; } // Might be removed if wrapped in ISigned<..>
        public Cid AgentId { get; set; }
    }
}