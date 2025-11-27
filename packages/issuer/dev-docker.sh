#!/bin/bash

# Development script using Docker with hot-reload
# This mounts your source code so changes take effect immediately

set -e

echo "🚀 Starting IU-MiCert Issuer in Docker Development Mode..."
echo ""
echo "📝 Features:"
echo "  - Source code mounted (changes take effect on restart)"
echo "  - PostgreSQL database"
echo "  - Adminer on http://localhost:8081"
echo "  - API on http://localhost:8080"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found, using defaults"
    echo "   Copy .env.example to .env and configure if needed"
fi

# Start services
docker compose -f docker-compose.dev.yml --profile admin up -d

echo ""
echo "✅ Services started!"
echo ""
echo "📊 View logs:"
echo "   docker compose -f docker-compose.dev.yml logs -f issuer-backend"
echo ""
echo "🔄 After code changes:"
echo "   docker compose -f docker-compose.dev.yml restart issuer-backend"
echo ""
echo "🛑 Stop services:"
echo "   docker compose -f docker-compose.dev.yml down"
echo ""
echo "🌐 Access points:"
echo "   API:     http://localhost:8080"
echo "   Adminer: http://localhost:8081"
echo ""
