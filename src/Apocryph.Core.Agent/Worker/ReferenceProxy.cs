using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;

namespace Apocryph.Core.Agent.Worker
{
    public class ReferenceProxy : DispatchProxy
    {
        private object _impl;
        private Dictionary<Guid, object> _attachedReferences;

        public void Init(object impl, Dictionary<Guid, object> attachedReferences)
        {
            _impl = impl;
            _attachedReferences = attachedReferences;
        }

        protected override object? Invoke(MethodInfo targetMethod, object[] args)
        {
            if (targetMethod.GetCustomAttributes<ReferenceAttachmentAttribute>().Any())
            {
                _attachedReferences[(Guid)args.Single()] = _impl;
            }

            return targetMethod.Invoke(_impl, args);
        }
    }
}