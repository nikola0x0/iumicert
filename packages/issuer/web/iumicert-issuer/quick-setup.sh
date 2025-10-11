#!/bin/bash

echo "🔄 IU-MiCert Quick Setup Script (5 Students)"
echo "============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Clear existing data
echo -e "\n${YELLOW}📁 Step 1: Clearing existing data...${NC}"
find data/ blockchain_ready/ -type f -name "*.json" -delete 2>/dev/null || true
rm -f *.json 2>/dev/null || true
echo -e "${GREEN}✅ Data cleared successfully${NC}"

# Step 2: Re-initialize the system
echo -e "\n${YELLOW}🏛️  Step 2: Re-initializing system...${NC}"
go run cmd/*.go init "IU-VNUHCM"
echo -e "${GREEN}✅ System initialized${NC}"

# Step 3: Generate student data (5 students)
echo -e "\n${YELLOW}👥 Step 3: Generating data for 5 students...${NC}"
go run cmd/*.go generate-data
echo -e "${GREEN}✅ Student data generated${NC}"

# Step 4: Batch process all terms
echo -e "\n${YELLOW}🌳 Step 4: Processing all terms (building Merkle/Verkle trees)...${NC}"
go run cmd/*.go batch-process
echo -e "${GREEN}✅ All terms processed${NC}"

# Step 5: Generate receipts for all 5 students
echo -e "\n${YELLOW}🎓 Step 5: Generating receipts for all students...${NC}"
mkdir -p receipts
students=("ITITIU00001" "ITITIU00002" "ITITIU00003" "ITITIU00004" "ITITIU00005")

for student in "${students[@]}"; do
    echo -e "  📋 Generating receipt for ${student}..."
    go run cmd/*.go generate-receipt $student receipts/${student}_journey.json
    echo -e "    ${GREEN}✅ ${student} receipt generated${NC}"
done

# Step 6: Show summary
echo -e "\n${BLUE}📊 Setup Summary:${NC}"
student_count=$(find data/generated_student_data/students/ -name "*.json" 2>/dev/null | wc -l)
receipt_count=$(find receipts/ -name "*.json" 2>/dev/null | wc -l)
echo -e "👥 Students generated: ${GREEN}${student_count}${NC}"
echo -e "🎓 Receipts generated: ${GREEN}${receipt_count}${NC}"

echo -e "\n${GREEN}🎉 Quick setup completed successfully!${NC}"
echo -e "${GREEN}You now have 5 students with complete receipts ready for blockchain publishing.${NC}"
