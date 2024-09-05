#!/usr/bin/env bash
bin=./bin

if [ ! -d $bin ]; then
	mkdir $bin
fi

ldflags="-s -w"

CGO_ENABLED=0 go build -ldflags="$ldflags" ./cmd/ls
mv ./ls ./bin/ls # Go is an awful choice for multiple binaries in a single project
