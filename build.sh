#!/bin/bash

#
# Build binaries for multiple architectures
#
# Inspiration:
# https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04

BINARY="bin/truecrypt_"

SETTINGS=("darwin_arm64" "darwin_amd64")

for suffix in "${SETTINGS[@]}"
do
    suffix_pieces=(${suffix/_/ }) # only God understands this line
	  GOOS=${suffix_pieces[0]}
	  GOARCH=${suffix_pieces[1]}
    BINARY_PATH="${BINARY}${suffix}"

    echo "Testing Os: ${GOOS} Arch: ${GOARCH}"
    env GOOS="${GOOS}" GOARCH="${GOARCH}" go clean -testcache || exit 1
    env GOOS="${GOOS}" GOARCH="${GOARCH}" go test ./... || exit 1

    echo "Building for Os: ${GOOS} Arch: ${GOARCH} Output: ${BINARY_PATH}"
    env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o "${BINARY_PATH}" || exit 1
done
