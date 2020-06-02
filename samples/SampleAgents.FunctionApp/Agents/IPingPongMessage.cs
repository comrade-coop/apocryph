using System;
using Apocryph.Core.Agent;

namespace SampleAgents.FunctionApp.Agents
{
    public interface IPingPongMessage
    {
        [ReferenceAttachment]
        Guid? AgentOne { get; set; }

        [ReferenceAttachment]
        Guid? AgentTwo { get; set; }

        string? Content { get; set; }
    }
}