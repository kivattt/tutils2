#!/usr/bin/env bash
bin=./bin

if [ ! -d $bin ]; then
	mkdir $bin
fi

ldflags="-s -w"

CGO_ENABLED=0 go build -ldflags="$ldflags" ./cmd/pwd && mv ./pwd ./bin/pwd
CGO_ENABLED=0 go build -ldflags="$ldflags" ./cmd/ls && mv ./ls ./bin/ls
CGO_ENABLED=0 go build -ldflags="$ldflags" ./cmd/ascii && mv ./ascii ./bin/ascii
CGO_ENABLED=0 go build -ldflags="$ldflags" ./cmd/xxd && mv ./xxd ./bin/xxd
CGO_ENABLED=0 go build -ldflags="$ldflags" ./cmd/countchars && mv ./countchars ./bin/countchars
