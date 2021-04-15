#!/usr/bin/bash

# Utils

set -e
set -o pipefail

function cleanup
{
  kill $(jobs -p) 2>/dev/null
}

trap cleanup EXIT

cd "$(dirname "$0")"

# CLI

INSTANCE_COUNT=2
START_IPFS=''
START_FABRIC=''
NO_BUILD=''
NO_RUN=''
DRY_RUN=''
DUMMY_CONSENSUS=''
_EXIT=''
for i in "$@"
do
case $i in
    -n=*|--instances=*)
    INSTANCE_COUNT="${i#*=}"
    ;;
    -B|--no-build|-b|--only-build)
    NO_BUILD=YES
    ;;
    -R|--no-run)
    NO_RUN=YES
    ;;
    -i|--ipfs)
    START_IPFS=YES
    ;;
    -f|--fabric)
    START_FABRIC=YES
    ;;
    -d|--dry-run)
    DRY_RUN=YES
    ;;
    --dummy-consensus)
    DUMMY_CONSENSUS=YES
    ;;
    -h|--help)
    echo "Usage: $0 [-n=2,--instances=2] [-B,--no-build] [-R,--no-run,-b,--only-build] [-i,--ipfs] [-f,--fabric] [-d,--dry-run] [-h,--help] [--dummy-consensus]"
    _EXIT=1
    ;;
    *)
    echo unrecognised option $i
    _EXIT=2
    ;;
esac
done

if [ $_EXIT ]; then set +e; exit $_EXIT; fi

# Configuration

START_AGENT=../samples/SampleAgents.FunctionApp/
export SAMPLE_AGENTS_CONSENSUS='Snowball'
if [ $DUMMY_CONSENSUS ]; then SAMPLE_AGENTS_CONSENSUS='Dummy'; fi
AGENTS=(
  ../src/Apocryph.Executor.FunctionApp/
  ../src/Apocryph.KoTH.FunctionApp/
  ../src/Apocryph.KoTH.SimpleMiner.FunctionApp/
  ../src/Apocryph.Routing.FunctionApp/
  ../src/Apocryph.Consensus.$SAMPLE_AGENTS_CONSENSUS.FunctionApp/
  $START_AGENT
)

declare -A INSTANCES
for instance in $(seq "$INSTANCE_COUNT"); do
  INSTANCES[$instance]="$((10800 + $instance - 1)) $((40400 + $instance - 1)) $((5001 + $instance - 1))"
done

# Command execution and logging

mkdir -p logs

function run_command # name pwd command
{
  NAME=$1; PWD=$2; CMD=$3
  PADDED_NAME=`printf '%-20s' "$NAME" | head -c 20`

  if [ ! $DRY_RUN ]; then
    script -qafec "cd $PWD; $CMD" "logs/$NAME.log" | sed -ue "s|^|$PADDED_NAME |"
  else
    echo "$PADDED_NAME cd $PWD; $CMD"
  fi
}

# IPFS

declare -A IPFS_PEERS

function configure_ipfs # instance ipfsport
{
  INSTANCE=$1; IPFS_PORT=$2

  SWARM_PORT=$(($IPFS_PORT + 5000))

  export IPFS_PATH=${TMPDIR:-/tmp}/test-ipfs-instance-$IPFS_PORT
  rm -rf $IPFS_PATH

  run_command "$INSTANCE-IPFS" . "ipfs init"
  run_command "$INSTANCE-IPFS" . "ipfs config profile apply test > /dev/null"
  run_command "$INSTANCE-IPFS" . "ipfs config --json Experimental.Libp2pStreamMounting true"
  run_command "$INSTANCE-IPFS" . "ipfs config Addresses.API /ip4/127.0.0.1/tcp/$IPFS_PORT"
  run_command "$INSTANCE-IPFS" . "ipfs config Pubsub.Router floodsub"
  run_command "$INSTANCE-IPFS" . "ipfs config --json Addresses.Swarm [\\\"/ip4/127.0.0.1/tcp/$SWARM_PORT\\\"]"

  if [ ! $DRY_RUN ]; then
    IPFS_PEERS[$INSTANCE]="/ip4/127.0.0.1/tcp/$SWARM_PORT/p2p/$(ipfs config Identity.PeerID)"
  fi
}

function run_ipfs # instance ipfsport
{
  INSTANCE=$1; IPFS_PORT=$2

  export IPFS_PATH=${TMPDIR:-/tmp}/test-ipfs-instance-$IPFS_PORT

  run_command "$INSTANCE-IPFS" . "ipfs daemon --enable-pubsub-experiment"
}

function connect_ipfs # instance ipfsport
{
  INSTANCE=$1; IPFS_PORT=$2

  export IPFS_PATH=${TMPDIR:-/tmp}/test-ipfs-instance-$IPFS_PORT

  for other_instance in "${!IPFS_PEERS[@]}"; do
    if [ $INSTANCE != $other_instance ]; then
      run_command "$INSTANCE-IPFS" . "ipfs swarm connect ${IPFS_PEERS[$other_instance]}"
    fi
  done
}

# Fabric

function run_fabric # instance igniteport grpcport
{
  INSTANCE=$1; IGNITE_PORT=$2; GRPC_PORT=$3

  run_command "$INSTANCE-FABRIC" . "docker run --rm -p $IGNITE_PORT:$IGNITE_PORT -p $GRPC_PORT:$GRPC_PORT -it obecto/perper-fabric:0.6.0-alpha7 --ignite-port $IGNITE_PORT --grpc-port $GRPC_PORT --no-discovery"
}

# Services

function start_services # instance igniteport grpcport ipfsport
{
  INSTANCE=$1; IGNITE_PORT=$2; GRPC_PORT=$3; IPFS_PORT=$4
  if ! fuser -n tcp $IGNITE_PORT &> /dev/null || fuser -n tcp $GRPC_PORT &> /dev/null; then
    if [ $START_FABRIC ]; then
      run_fabric $INSTANCE $IGNITE_PORT $GRPC_PORT &
    else
      echo "WARNING: Perper Fabric is not running on ports $IGNITE_PORT, $GRPC_PORT";
    fi
  fi
  if ! fuser -n tcp $IPFS_PORT &> /dev/null; then
    if [ $START_IPFS ]; then
      configure_ipfs $INSTANCE $IPFS_PORT # NOTE: Should run in foreground, modifies global variable
      run_ipfs $INSTANCE $IPFS_PORT &
    else
      echo "WARNING: IPFS is not running on port $IPFS_PORT";
    fi
  fi
}

# Agents

function get_agent_name # path
{
  sed -nEe 's|.*"PERPER_AGENT_NAME":.*?"([^"]+)".*|\1|p' $1/local.settings.json
}

function build_agent # instance path igniteport grpcport ipfsport
{
  INSTANCE=$1; PWD=$2; IGNITE_PORT=$3; GRPC_PORT=$4; IPFS_PORT=$5
  AGENTNAME=`get_agent_name $PWD`

  run_command "$INSTANCE-$AGENTNAME" $PWD "dotnet build -p:'FabricIgnitePort=$IGNITE_PORT;FabricGrpcPort=$GRPC_PORT;IpfsPort=$IPFS_PORT' --output bin/output$INSTANCE"
}

function run_agent # instance path
{
  INSTANCE=$1; PWD=$2; IGNITE_PORT=$3; GRPC_PORT=$4; IPFS_PORT=$5
  AGENTNAME=`get_agent_name $PWD`

  run_command "$INSTANCE-$AGENTNAME" $PWD "cd bin/output$INSTANCE; func start --no-build -p 0"
}

# Execution

if [ ! $NO_BUILD ]; then
  for agent in "${AGENTS[@]}"; do
    for instance in "${!INSTANCES[@]}"; do
      build_agent $instance $agent ${INSTANCES[$instance]} &
      if [ ! $DRY_RUN ]; then
        sleep 2 # Otherwise msbuild gets into a race
      fi
    done
  done

  wait; if [ $DRY_RUN ]; then echo; fi
fi

if [ $START_IPFS ] || [ $START_FABRIC ]; then
  for instance in "${!INSTANCES[@]}"; do
    start_services $instance ${INSTANCES[$instance]} # NOTE: Should run in foreground, modifies global variable
  done

  if [ ! $DRY_RUN ]; then
    sleep 7 # Give fabric/ipfs some time to start up
    for instance in "${!INSTANCES[@]}"; do
      ports=(${INSTANCES[$instance]})
      connect_ipfs $instance ${ports[2]} &
    done
  fi
fi

if [ ! $NO_RUN ]; then
  export PERPER_ROOT_AGENT=`get_agent_name $START_AGENT`
  for agent in "${AGENTS[@]}"; do
    for instance in "${!INSTANCES[@]}"; do
      run_agent $instance $agent ${INSTANCES[$instance]} &
    done
  done

  wait; if [ $DRY_RUN ]; then echo; fi
fi
