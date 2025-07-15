#!/bin/bash

# Build script for Waterlogger with build timestamp
# This script automatically injects the build date and time into the binary

BUILD_TIME=$(date '+%H:%M:%S')
BUILD_DATE=$(date '+%Y-%m-%d')

echo "Building Waterlogger..."
echo "Build Date: $BUILD_DATE"
echo "Build Time: $BUILD_TIME"

go build -ldflags "-X main.BuildTime=$BUILD_TIME -X main.BuildDate=$BUILD_DATE" -o waterlogger ./cmd/waterlogger

if [ $? -eq 0 ]; then
    echo "Build completed successfully!"
    echo "Binary: ./waterlogger"
else
    echo "Build failed!"
    exit 1
fi