#!/usr/bin/env bash

set -e

go get -u github.com/alecthomas/gometalinter
gometalinter --install

go get -t -v ./...
