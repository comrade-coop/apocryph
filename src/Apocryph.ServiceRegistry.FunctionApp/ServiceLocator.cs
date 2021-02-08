using Perper.WebJobs.Extensions.Model;

namespace Apocryph.ServiceRegistry
{
    [PerperData]
    public struct ServiceLocator
    {
        public ServiceLocator(string type, string id)
        {
            Type = type;
            Id = id;
        }

        public string Type { get; private set; }
        public string Id { get; private set; }
    }
}