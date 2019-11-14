#!/bin/bash

VERSION=$1

go test ./... -v -cover -count=1
go build scheduler.go
tar -czvf oscap-exporter-${VERSION}.tar.gz scheduler