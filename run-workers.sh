#!/bin/bash

# Get count from command line
COUNT=$1

if [ -z "$COUNT" ]; then
  echo "Usage: $0 <count>"
  exit 1
fi

for i in $(seq 0 $((COUNT-1))); do
  export PORT=$((4000+i))
  go run nanoray/worker &
  sleep 0.2
done

sleep infinity

# On exit, kill all workers
trap 'pkill -f worker' EXIT