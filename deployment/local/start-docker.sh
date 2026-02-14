#!/bin/bash

# Script to copy .env file and start Docker containers

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
LOCAL_DIR="$SCRIPT_DIR"

echo "📋 Starting Docker setup..."
echo "Root directory: $ROOT_DIR"
echo "Local directory: $LOCAL_DIR"

# Copy .env file from root to local directory
if [ -f "$ROOT_DIR/.env" ]; then
    echo "📋 Copying .env file from root to deployment/local..."
    cp "$ROOT_DIR/.env" "$LOCAL_DIR/.env"
    echo "✅ .env file copied successfully"
else
    echo "⚠️  .env file not found in root directory"
    echo "📋 Using existing .env in deployment/local (if available)"
fi

# Start Docker containers
echo "🐳 Starting Docker containers..."
cd "$LOCAL_DIR"
docker-compose --env-file .env up -d

echo "✅ Docker containers started successfully!"
echo ""
echo "📊 Service URLs:"
echo "  - Grafana: http://localhost:3000"
echo "  - Prometheus: http://localhost:9090"
echo "  - Tempo: http://localhost:3200"
echo "  - RabbitMQ Management: http://localhost:15672"
echo "  - PostgreSQL: localhost:5432"
echo "  - Valkey: localhost:6379"
echo ""
echo "💡 To view logs: docker-compose logs -f"
echo "💡 To stop containers: docker-compose down"
