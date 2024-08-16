#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

handler()
{
  kill -s SIGINT "$PID"
}

exec /usr/local/bin/honor-downloader "$@" &
PID=$!
trap 'handler $PID' SIGTERM
wait "$PID"
