#!/usr/bin/bash

set -e
set -o pipefail

function run_command
{
  NAME=$1; PWD=$2; CMD=$3
  PADDED_NAME=`printf '%-20s' "$NAME" | head -c 20`

  script -qfec "cd $PWD; $CMD" "logs/$NAME.log" | sed -ue "s|^|$PADDED_NAME |"
}

function get_agent_name # path
{
  sed -nEe 's|.*"PERPER_AGENT_NAME":.*?"([^"]+)".*|\1|p' $1/local.settings.json
}

function run_agent # path
{
  PWD=$1
  AGENTNAME=`get_agent_name $PWD`

  run_command $AGENTNAME $PWD 'func start -p 0'
}

function cleanup {
  kill $(jobs -p) 2>/dev/null
}

trap cleanup EXIT

cd "$(dirname "$0")"
mkdir -p logs

export PERPER_ROOT_AGENT=`get_agent_name .`

AGENTS=(
  .
  ../../src/Apocryph.Executor.FunctionApp/
  ../../src/Apocryph.KoTH.FunctionApp/
#   ../../src/Apocryph.KoTH.SimpleMiner.FunctionApp/
  ../../src/Apocryph.Routing.FunctionApp/
  ../../src/Apocryph.Consensus.Dummy.FunctionApp/
)

fuser -n tcp 5001 &> /dev/null || { echo "IPFS is not running on port 5001" && exit 1; }
{ fuser -n tcp 10800 &> /dev/null && fuser -n tcp 40400 &> /dev/null; } || { echo "Perper Fabric is not running on ports 10800, 40400" && exit 1; }

for agent in "${AGENTS[@]}"; do
  run_agent $agent &
  sleep 1
done

wait
