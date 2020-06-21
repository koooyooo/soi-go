#!/bin/sh

cd $(dirname $0)

go run ../cmd/soi.go "$@"
