#!/bin/bash -x

if [ "$DEBUG_MODE" == 1 ]
then
  # shellcheck disable=SC2160
  while [ true ]
  do
    sleep 5
  done
fi

echo "Starting Trade Server..."
bin/tradeserver

# Stop docker from exiting
tail -f /dev/null