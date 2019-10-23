#!/bin/bash

VERSION=$1

go test
go build scheduler.go
tar -czvf oscap-exporter-${VERSION}.tar.gz scheduler