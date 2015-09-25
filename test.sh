#!/bin/sh -e

# walter: a deployment pipeline template
# Copyright (C) 2014 Recruit Technologies Co., Ltd. and contributors
# (see CONTRIBUTORS.md)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

. ./build

go test -i ./config || exit 1
go test -v ./config -race || exit 1

go test -i ./pipelines || exit 1
go test -v ./pipelines -race || exit 1

go test -i ./walter || exit 1
go test -v ./walter -race || exit 1

go test -i ./stages || exit 1
go test -v ./stages -race || exit 1

go test -i ./engine || exit 1
go test -v ./engine -race || exit 1

go test -i ./messengers || exit 1
go test -v ./messengers -race || exit 1

go test -i ./services || exit 1
go test -v ./services -race || exit 1

fmtRes=`gofmt -l $GOFMTPATH`
if [ "$fmtRes" != "" ]; then
    echo "Failed to pass golang format checking."
    echo "Please gofmt modified go files, or run './build --fmt'."
    exit 1
fi

if [ "--lint" = "$1" ]; then
    for WALTER_SOURCE in $GOFMTPATH
    do
	lintRes=`golint $WALTER_SOURCE`
	if [ "$lintRes" != "" ]; then
	    echo "Failed to pass source code linting."
	    echo $lintRes
	    echo "Please fix the errors."
	    exit 1
	fi
    done
fi
