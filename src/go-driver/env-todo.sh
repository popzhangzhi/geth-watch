#!/bin/sh
#abandon
set -e

#if [ ! -f "build/env.sh" ]; then
#   echo "$0 must be run from the root of the repository."
#  exit 2
#fi

# Create fake Go workspace if it doesn't exist yet.

root="$PWD"
workspace="$root/src"


# Set up the environment to use the workspace.
GOPATH="$root/../go-ethereum/build/_workspace:$root"
export GOPATH
GOBIN="$root/bin"
export GOBIN

# Run the command inside the workspace.
cd "$workspace"
PWD="$workspace"

#输出当前运行go的env
#go env

# Launch the arguments with the configured environment.
exec "$@"

#直接运行测试文件
go run test.go

#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
