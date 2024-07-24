#!/bin/bash
clc -s -e ufile_test.go
cat Version.dat
go mod tidy
go fmt .
echo no staticcheck due to rangefunc use
# staticcheck .
go vet .
echo no golangci-line due to rangefunc use
# golangci-lint run
git st
