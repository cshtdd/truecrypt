#!/bin/bash

SCRIPT_DIR="$(dirname -- "$0")"

pushd $SCRIPT_DIR
./truecrypt clean
popd
