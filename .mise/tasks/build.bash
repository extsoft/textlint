#!/usr/bin/env bash
#MISE description="builds the project"
set -o errexit -o nounset

: "${DIST_DIR:=dist}"
go mod tidy
VERSION=$(git describe --tags --always --abbrev=7 2>/dev/null || echo "dev")
go build -ldflags="-s -w -X prosefmt/cmd/prosefmt.version=${VERSION}" -o "${DIST_DIR}/prosefmt" .
