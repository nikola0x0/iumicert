# Essential Commands for Development

## Backend Development

### CLI Commands
```bash
# Complete reset and regenerate data
./reset.sh && ./generate.sh

# Start API server
./dev.sh  # or ./micert serve --port 8080 --cors

# Manual data pipeline
./micert generate-data        # Create test data
./micert batch-process        # Process all terms
./micert generate-all-receipts # Create receipts

# Blockchain operations (requires ISSUER_PRIVATE_KEY)
./micert publish-roots        # Publish to blockchain

# Verification
./micert verify-local <receipt.json>
./micert display-receipt <receipt.json>
```

### Go Development
```bash
# Build the CLI
go build -o micert cmd/*.go

# Run with verbose logging  
./micert --verbose <command>

# Get help for any command
./micert <command> --help
```

## Frontend Development

### Web App
```bash
cd web/iumicert-issuer

# Install dependencies
npm install

# Development server
npm run dev  # Runs on http://localhost:3000

# Build for production
npm run build

# Start production server
npm start

# Lint code
npm run lint
```

## System Management

### Environment Setup
```bash
# Configure blockchain access
cp .env.example .env
# Edit .env with ISSUER_PRIVATE_KEY

# Check system health
curl http://localhost:8080/api/health
```

### Data Management
```bash
# Check data status
ls data/student_journeys/students/
ls publish_ready/receipts/

# Clean reset
./reset.sh

# Full pipeline regeneration
./generate.sh
```

## Testing & Verification
```bash
# Test full system
./micert version
./micert verification-guide

# API endpoints test
curl http://localhost:8080/api/terms
curl http://localhost:8080/api/students

# Frontend development
cd web/iumicert-issuer && npm run dev
```