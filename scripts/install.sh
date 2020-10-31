#!/bin/sh

cd "$(dirname $0)" || exit

go install ../cmd/soi2/soi2.go
