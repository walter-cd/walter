#!/bin/sh -e

. ./build

go test -i ./config
go test -v ./config -race

go test -i ./pipelines
go test -v ./pipelines -race

go test -i ./plumber
go test -v ./plumber -race

go test -i ./stages
go test -v ./stages -race

#go test -i ./tests/functional
#ETCD_BIN_PATH=$(pwd)/bin/plumber go test -v ./tests/functional -race

fmtRes=`gofmt -l $GOFMTPATH`
if [ "$fmtRes" != "" ]; then
	echo "Failed to pass golang format checking."
	echo "Please gofmt modified go files, or run './build --fmt'."
	exit 1
fi
