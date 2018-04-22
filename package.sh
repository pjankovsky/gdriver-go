#!/usr/bin/env bash

BUILD=`git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD`

tar -czvf gdriver-go.${BUILD}.tar.gz --transform 's,^,gdriver-go/,' gdriver-go html etc/settings.json
