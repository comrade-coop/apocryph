using System.Collections.Generic;
using Perper.WebJobs.Extensions.Model;
using IInputStream = System.String; // From CreateBlankStream

namespace Apocryph.ServiceRegistry
{
    [PerperData]
    public class Service
    {
        public Service(Dictionary<string, IInputStream> inputs, Dictionary<string, IStream> outputs)
        {
            Inputs = inputs;
            Outputs = outputs;
        }

        public Dictionary<string, IInputStream> Inputs { get; set; }
        public Dictionary<string, IStream> Outputs { get; set; }
    }
}