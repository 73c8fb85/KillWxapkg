#!/bin/bash
set -e

cd "$(dirname "$0")/.."

export CGO_ENABLED=1

if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "Building macOS dylib..."
    GOOS=darwin GOARCH=amd64 go build -buildmode=c-shared -o libwxapkg_killer.dylib ./ffi
    echo "Build successful: libwxapkg_killer.dylib"
elif [[ "$OSTYPE" == "linux"* ]]; then
    echo "Building Linux so..."
    GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o libwxapkg_killer.so ./ffi
    echo "Build successful: libwxapkg_killer.so"
fi

echo "Copy the library to your Flutter project"
