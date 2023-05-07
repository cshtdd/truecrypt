#!/bin/bash

SCRIPT_DIR="$(dirname -- "$0")"
OS_NAME="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH_NAME="$(uname -m)"
$SCRIPT_DIR/bin/truecrypt_${OS_NAME}_${ARCH_NAME} -encrypt
