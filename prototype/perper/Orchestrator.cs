using Microsoft.Azure.WebJobs;
using Mock.Perper;

public static class Orchestrator
{
    [FunctionName("Orchestrator")]
    public static void Run(PerperStreamContext context)
    {
        var ipfsInput = context.CallStreamFunction("IPFS_Input");

        var consensusRouter = context.CallMultiStreamFunction("ConsensusRouter");
        var consensusVoting = context.CallMultiStreamFunction("ConsensusVoting", consensusRouter["vote"]);
        var consensusPre = context.CallMultiStreamFunction("ConsensusPre", consensusVoting, consensusRouter["vote"]);

        var executionPre = context.CallMultiStreamFunction("ExecutionPre", consensusPre["execution"]);

        var contractsInput = context.CallStreamFunction("Messaging_Input", executionPre["contract"]);
        var contractsOutput = context.CallStreamFunction("Messaging_Output", contractsInput);

        var executionPost = context.CallMultiStreamFunction("ExecutionPost", executionPre["execution"], contractsOutput);

        var consensusPost = context.CallStreamFunction("ConsensusPost", consensusPre["consensus"], executionPost["output"]);

        context.CallActivityFunction("IPFS_Output", executionPost["output"]);
    }
}
