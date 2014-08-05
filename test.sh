#!/bin/sh -e

. ./build

go test -i ./pipeline
go test -v ./pipeline -race

go test -i ./plumber
go test -v ./plumber -race

go test -i ./stage
go test -v ./stage -race

#go test -i ./tests/functional
#ETCD_BIN_PATH=$(pwd)/bin/plumber go test -v ./tests/functional -race

fmtRes=`gofmt -l $GOFMTPATH`
if [ "$fmtRes" != "" ]; then
	echo "Failed to pass golang format checking."
	echo "Please gofmt modified go files, or run './build --fmt'."
	exit 1
fi
