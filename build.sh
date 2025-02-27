#!/bin/bash

# Set the output directories
OUTPUT_DIR="build"
SERVER_DIR="http_server"
CLIENT_DIR="http_client"

# Create the output directories if they don't exist
mkdir -p $OUTPUT_DIR/macos_arm64
mkdir -p $OUTPUT_DIR/linux_amd64

# Compile http_server for macOS ARM
echo "Compiling http_server for macOS ARM..."
GOOS=darwin GOARCH=arm64 go build -o $OUTPUT_DIR/macos_arm64/http_server $SERVER_DIR/main.go

# Compile http_server for Linux amd64
echo "Compiling http_server for Linux amd64..."
GOOS=linux GOARCH=amd64 go build -o $OUTPUT_DIR/linux_amd64/http_server $SERVER_DIR/main.go

# Compile http_client for macOS ARM
echo "Compiling http_client for macOS ARM..."
GOOS=darwin GOARCH=arm64 go build -o $OUTPUT_DIR/macos_arm64/http_client $CLIENT_DIR/main.go

# Compile http_client for Linux amd64
echo "Compiling http_client for Linux amd64..."
GOOS=linux GOARCH=amd64 go build -o $OUTPUT_DIR/linux_amd64/http_client $CLIENT_DIR/main.go

echo "Compilation finished."