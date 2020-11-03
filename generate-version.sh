#!/bin/sh
set -ue

VERSION_FILE=$1

COMMIT_HASH=$(git rev-parse --short HEAD)

cat > $VERSION_FILE <<EOF
package cli

func init() { version="$COMMIT_HASH" }

EOF
