#!/bin/bash

echo "🧹 IU-MiCert Reset Script (Single Verkle Architecture)"
echo "======================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

echo -e "${PURPLE}🌳 Resetting Single Verkle Tree System${NC}"
echo -e "${PURPLE}✨ Clean architecture: No Merkle dependencies${NC}"

# Step 1: Clear all existing data
echo -e "\n${YELLOW}📁 Step 1: Clearing all existing data...${NC}"

# Remove all data directories
rm -rf data/student_journeys/* 2>/dev/null || true
rm -rf data/verkle_terms/* 2>/dev/null || true
rm -rf publish_ready/receipts/* 2>/dev/null || true
rm -rf publish_ready/roots/* 2>/dev/null || true
rm -rf publish_ready/proofs/* 2>/dev/null || true
rm -rf publish_ready/transactions/* 2>/dev/null || true

# Remove any leftover files
rm -f *.json 2>/dev/null || true
rm -f cmd/*.json 2>/dev/null || true

# Recreate directory structure
mkdir -p data/student_journeys/students
mkdir -p data/student_journeys/terms  
mkdir -p data/verkle_terms
mkdir -p publish_ready/{receipts,roots,proofs,transactions}

echo -e "${GREEN}✅ All data cleared and directories recreated${NC}"

# Step 2: Show clean structure
echo -e "\n${BLUE}📊 Clean Directory Structure:${NC}"
echo -e "${BLUE}==============================${NC}"
echo -e "${BLUE}📁 data/${NC}"
echo -e "${BLUE}  ├── student_journeys/    # Generated academic journeys${NC}"
echo -e "${BLUE}  └── verkle_terms/        # Terms converted for single Verkle${NC}"
echo -e "${BLUE}📁 publish_ready/${NC}"
echo -e "${BLUE}  ├── receipts/            # Student verification receipts${NC}"
echo -e "${BLUE}  ├── roots/               # Verkle tree roots for blockchain${NC}"
echo -e "${BLUE}  ├── proofs/              # Cryptographic proofs (32-byte Verkle)${NC}"
echo -e "${BLUE}  └── transactions/        # Blockchain transaction records${NC}"

echo -e "\n${GREEN}🎉 Reset completed successfully!${NC}"
echo -e "${GREEN}System ready for single Verkle tree generation.${NC}"
echo -e "\n${YELLOW}Next steps:${NC}"
echo -e "  1. Run ${PURPLE}./generate.sh${NC} to create sample data"
echo -e "  2. Or manually run: ${PURPLE}cd cmd && go run . generate-data${NC}"