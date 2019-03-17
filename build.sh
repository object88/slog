#! /usr/bin/env bash

cd $(dirname "$0")

mkdir -p ./bin
mkdir -p ./mocks

export PATH=$PWD/bin:$PATH
if ! [ -x $PWD/bin/mockgen ]; then
  echo "Building mockgen"
  go build -o $PWD/bin/mockgen ./vendor/github.com/golang/mock/mockgen
fi

export GO111MODULE=on

echo "Running generator"
go generate ./...

echo "Building..."
go build -ldflags "-s -w" -mod=readonly -mod=vendor -o ./bin/slog ./main/main.go

go test ./... -count=1

echo "Done."
