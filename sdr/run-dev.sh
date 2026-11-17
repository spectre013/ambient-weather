#!/bin/sh
# Run a command in the dev container (Go toolchain + rtl-sdr tools) so nothing
# has to be installed on the host.  ./run-dev.sh go test ./...
cd "$(dirname "$0")" || exit 1
exec docker run --rm -i -v "$PWD":/src -w /src -u "$(id -u):$(id -g)" \
    -e HOME=/tmp -e GOFLAGS=-buildvcs=false -v sdr-dev-cache:/cache -e GOCACHE=/cache/go-build -e GOMODCACHE=/cache/mod \
    $DOCKER_ARGS sdr-dev "$@"
