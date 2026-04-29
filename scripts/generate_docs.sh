#!/bin/bash

# Navigate to the backend directory (assuming the script is run from anywhere inside backend)
cd "$(dirname "$0")/.."

echo "Generating Swagger documentation..."
swag init -g cmd/routing_app_backend/main.go --parseDependency --parseInternal

echo "Swagger documentation generation complete."
