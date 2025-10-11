#!/bin/bash
echo "🚀 Starting IU-MiCert Issuer Development Server..."
go run cmd/*.go serve --port 8080 --cors