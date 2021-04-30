#!/usr/bin/env bash

set -xe

SDK_ROOT=$1
CODGEN_ROOT=$2

cd "$1"/internal/repotools
go run ./cmd/gomodgen -build "$CODGEN_ROOT"
