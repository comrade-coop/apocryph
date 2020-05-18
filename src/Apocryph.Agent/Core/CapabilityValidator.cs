using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;

namespace Apocryph.Agent.Core
{
    public class CapabilityValidator
    {
        private readonly IDictionary<Guid, (string, string[])> _capabilities;
        private readonly HashSet<string> _carriers;

        public CapabilityValidator(IDictionary<Guid, (string, string[])> capabilities)
        {
            _capabilities = capabilities;
            _carriers = new HashSet<string>();
        }

        public void AddCapabilities(IDictionary<Guid, (string, string[])> capabilities)
        {
            foreach (var (key, value) in capabilities)
            {
                _capabilities.Add(key, value);
            }
        }

        public void RegisterCarrier(byte[] carrierPayload)
        {
            using var sha1 = new SHA1CryptoServiceProvider();
            _carriers.Add(Convert.ToBase64String(sha1.ComputeHash(carrierPayload)));
        }

        public bool ValidateMessageAndRegisterAsCarrier(Guid reference, (string, byte[]) message)
        {
            if (!_capabilities.TryGetValue(reference, out var capability)) return false;
            var (_, messageTypes) = capability;
            var (messageType, messagePayload) = message;
            var result = messageTypes.Contains(messageType);
            if (result)
            {
                RegisterCarrier(messagePayload);
            }

            return result;
        }

        public bool ValidateAttachedReference((Guid, string) attachedReference)
        {
            var (reference, carrier) = attachedReference;
            return _capabilities.ContainsKey(reference) && _carriers.Contains(carrier);
        }
    }
}