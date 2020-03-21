#!/bin/bash
set -e
echo "Building for linux amd64 statically"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o az-offerlist

