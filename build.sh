#! /usr/bin/env bash

cd $(dirname "$0")

mkdir -p ./bin

export GO111MODULE=on

go build -ldflags "-s -w" -mod=readonly -mod=vendor -o ./bin/slog ./main/main.go
