#!/usr/bin/env bash
set -e

# first arg is `-f` or `--some-option`
if [ "${1#-}" != "$1" ]; then
	set -- "app" "$@"
fi

exec "$@"