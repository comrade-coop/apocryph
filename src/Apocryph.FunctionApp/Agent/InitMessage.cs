using System;
using System.Collections.Generic;
using Apocryph.FunctionApp.Command;

namespace Apocryph.FunctionApp.Agent
{
    public struct InitMessage
    {
        public override bool Equals(object? obj)
        {
            if (obj is InitMessage other)
            {
                return true;
            }

            return false;
        }

        public override int GetHashCode()
        {
            return 0;
        }
    }
}