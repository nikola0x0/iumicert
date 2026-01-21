#!/bin/bash
echo "🔨 Building IU-MiCert..."
go build -o micert ./cmd || { echo "❌ Build failed"; exit 1; }
echo "🚀 Starting IU-MiCert Issuer Development Server..."
./micert serve --port 8080 --cors