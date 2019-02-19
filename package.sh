#!/usr/bin/env bash

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    TARCMD="tar"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    command -v gtar >/dev/null 2>&1 || { echo >&2 "I require 'gtar' but it's not installed. Install with 'brew install gnu-tar'."; exit 1; }
    TARCMD="gtar"
else
    TARCMD="tar"
fi

env GOOS=linux GOARCH=amd64 go build

BUILD=`git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD`

${TARCMD} -czvf gdriver-go.${BUILD}.tar.gz --transform 's,^,gdriver-go/,' gdriver-go html etc/settings.json
