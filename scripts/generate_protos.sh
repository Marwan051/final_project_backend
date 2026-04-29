#!/bin/bash

# Exit on error
set -e

# Change to the root of the project
cd "$(dirname "$0")/.."

echo "Generating proto stubs..."

PROTO_DIR="internal/service"

# Check if protoc and necessary plugins are installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc could not be found. Please install the protobuf compiler."
    exit 1
fi

if ! command -v protoc-gen-go &> /dev/null; then
    echo "Error: protoc-gen-go could not be found."
    echo "Please run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Error: protoc-gen-go-grpc could not be found."
    echo "Please run: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

# Find all proto files and compile them
find "$PROTO_DIR" -name "*.proto" -type f -print0 | while IFS= read -r -d '' proto_file; do
    echo "Compiling $proto_file..."
    protoc \
        --proto_path=. \
        --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        "$proto_file"
done

echo "Proto stubs generated successfully!"
