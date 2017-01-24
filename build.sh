#!/bin/sh

$GOPATH/bin/go-bindata data www/...

go build
