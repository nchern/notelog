#!/bin/sh
set -ue

coverprofile="artifacts/coverage.out"

coverprofile_raw="${coverprofile}.raw"

pkg="./pkg/..."

covermode=${COVER_MODE:-"set"}

go test -coverpkg="$pkg" -coverprofile="$coverprofile_raw" -covermode="$covermode" ./...
cat "$coverprofile_raw" \
    | grep -v ".pb." \
    | grep -v ".xo." \
    > $coverprofile

if [ "${1:-""}" = html ]; then
    go tool cover -html="$coverprofile"
else
    go tool cover -func="$coverprofile"
fi
