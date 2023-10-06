#!/usr/bin/env bash

set -e


if [ "$1" = "teardown" ]; then
   killall substrate-contracts-node
   exit 0
fi

cd "$(dirname "$0")"

which go >/dev/null
which cargo >/dev/null
which substrate-contracts-node >/dev/null # cargo install contracts-node
which cargo-contract >/dev/null # cargo install cargo-contract
which jq >/dev/null

set -v

TMP="$(mktemp substrate-contracts-node.XXXX --tmpdir)"; echo $TMP; nohup substrate-contracts-node >"$TMP" &
{ while ! echo -n > /dev/tcp/localhost/9944; do sleep 0.1; done; } 2>/dev/null # https://stackoverflow.com/a/44484835

CONTRACT="$(cd ../../../contracts/payment && cargo contract instantiate --suri //Alice -x --args 10000 --skip-confirm --output-json | jq .contract -r)"; echo $CONTRACT

set -x

go run ../../../cmd/tpodserver/ contract check "$CONTRACT" --config config.yaml --rpc ws://localhost:9944/

go run ../../../cmd/tpodserver/ contract claim "$CONTRACT" 123 --config config.yaml --rpc ws://localhost:9944/

go run ../../../cmd/tpodserver/ contract check "$CONTRACT" --config config.yaml --rpc ws://localhost:9944/

