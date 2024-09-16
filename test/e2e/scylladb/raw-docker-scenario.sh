#!/usr/bin/sh

set -e
set -v

check_if_ready () {
  docker exec "$1" bash -c "while ! cat /proc/net/tcp | grep $(printf '%x\n' 9042) >/dev/null; do sleep 1; done"
}

if which cqlsh > /dev/null; then
  echo "Using cqlsh via docker! (consider installing it with 'pip install cqlsh')"
  docker start cqlsh-container || docker run --rm -d --name cqlsh-container --network host --entrypoint sleep scylladb/scylla 10000
  # https://stackoverflow.com/a/44739847
  alias cqlsh="docker exec cqlsh-container cqlsh"
fi

docker compose up -d

check_if_ready scyllat1
check_if_ready scyllat2

# cqlsh 172.123.0.2 -e "CREATE KEYSPACE IF NOT EXISTS test WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': 1} AND tablets = { 'enabled': false }; CREATE TABLE IF NOT EXISTS test.test1 (id uuid PRIMARY KEY, count counter);"
cqlsh 172.123.0.2 -e "CREATE KEYSPACE IF NOT EXISTS test WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 2} AND tablets = { 'enabled': false }; CREATE TABLE IF NOT EXISTS test.test1 (id uuid PRIMARY KEY, count counter);"

cqlsh 172.123.0.3 -e "UPDATE test.test1 SET count = count + 4 where id = 6c89c24d-0690-4dac-8f6e-71e98ad2b165;"

cqlsh 172.123.0.2 -e "SELECT * FROM test.test1"
cqlsh 172.123.0.2 -e "UPDATE test.test1 SET count = count + 5 where id = 6c89c24d-0690-4dac-8f6e-71e98ad2b165;"
cqlsh 172.123.0.2 -e "UPDATE test.test1 SET count = count + 5 where id = 1664d2a0-c9c6-4318-8373-a57cc462608b;"

cqlsh 172.123.0.3 -e "SELECT * FROM test.test1"

docker stop -t 0 scyllat2

cqlsh 172.123.0.2 -e "CONSISTENCY ALL; SELECT * FROM test.test1 USING TIMEOUT 200ms" || true # Doesn't work, that's expected 
cqlsh 172.123.0.2 -e "SELECT * FROM test.test1" # Still works!
cqlsh 172.123.0.2 -e "UPDATE test.test1 SET count = count + 6 where id = 1664d2a0-c9c6-4318-8373-a57cc462608b;" # Still works!

docker start scyllat2

check_if_ready scyllat2

sleep 1

docker stop -t 0 scyllat1

cqlsh 172.123.0.3 -e "CONSISTENCY ALL; SELECT * FROM test.test1 USING TIMEOUT 200ms" || true # Doesn't work, that's expected 
cqlsh 172.123.0.3 -e "SELECT * FROM test.test1" # Still works! (But counts might be off because we most likely interrupted sync)

docker compose down
