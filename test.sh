#!/usr/bin/env bash

go build -o ./heimdall bifrost/main.go
go test
