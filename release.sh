#!/bin/sh

#
# Build release versions for all supported platforms to make it easy to
# upload to GitHub or your own company storage (Google Drive, Dropbox, S3).
#

# Build for supported operating systems & architectures
GOOS=darwin GOARCH=amd64 go build -o tech-talk-mac
GOOS=linux GOARCH=386 go build -o tech-talk-linux-386
GOOS=linux GOARCH=amd64 go build -o tech-talk-linux-amd64
GOOS=windows GOARCH=386 go build
