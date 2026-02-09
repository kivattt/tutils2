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
CGO_ENABLED=0 go build -ldflags="$ldflags" ./cmd/cat && mv ./cat ./bin/cat
CGO_ENABLED=0 go build -ldflags="$ldflags" ./cmd/hex && mv ./hex ./bin/hex
CGO_ENABLED=0 go build -ldflags="$ldflags" ./cmd/bytes && mv ./bytes ./bin/bytes
CGO_ENABLED=0 go build -ldflags="$ldflags" ./cmd/dirstats && mv ./dirstats ./bin/dirstats
