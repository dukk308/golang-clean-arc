#!/bin/bash

# Script to stop Docker containers

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🛑 Stopping Docker containers..."
cd "$SCRIPT_DIR"
docker-compose down

echo "✅ Docker containers stopped successfully!"
echo ""
echo "💡 To remove volumes as well: docker-compose down -v"
echo "💡 To view logs: docker-compose logs -f"
