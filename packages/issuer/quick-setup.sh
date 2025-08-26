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
rm -rf data/merkle_trees/* 2>/dev/null || true
rm -rf data/verkle_trees/* 2>/dev/null || true  
rm -rf data/generated_student_data/* 2>/dev/null || true
rm -rf data/converted_terms/* 2>/dev/null || true
rm -rf blockchain_ready/receipts/* 2>/dev/null || true
rm -rf blockchain_ready/roots/* 2>/dev/null || true
rm -rf blockchain_ready/transactions/* 2>/dev/null || true
rm -f *.json 2>/dev/null || true
echo -e "${GREEN}✅ Data cleared successfully${NC}"

# Step 2: Re-initialize the system
echo -e "\n${YELLOW}🏛️  Step 2: Re-initializing system...${NC}"
go run cmd/*.go init "IU-VNUHCM"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ System initialized${NC}"
else
    echo -e "${RED}❌ Failed to initialize system${NC}"
    exit 1
fi

# Step 3: Generate student data (5 students)
echo -e "\n${YELLOW}👥 Step 3: Generating data for 5 students...${NC}"
go run cmd/*.go generate-data
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Student data generated${NC}"
else
    echo -e "${RED}❌ Failed to generate data${NC}"
    exit 1
fi

# Step 4: Batch process all terms
echo -e "\n${YELLOW}🌳 Step 4: Processing all terms (building Merkle/Verkle trees)...${NC}"
go run cmd/*.go batch-process
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ All terms processed${NC}"
else
    echo -e "${RED}❌ Failed to process terms${NC}"
    exit 1
fi

# Step 5: Generate receipts for all 5 students
echo -e "\n${YELLOW}🎓 Step 5: Generating receipts for all students...${NC}"
mkdir -p receipts
students=("ITITIU00001" "ITITIU00002" "ITITIU00003" "ITITIU00004" "ITITIU00005")

for student in "${students[@]}"; do
    echo -e "  📋 Generating receipt for ${student}..."
    go run cmd/*.go generate-receipt $student receipts/${student}_journey.json
    if [ $? -eq 0 ]; then
        echo -e "    ${GREEN}✅ ${student} receipt generated${NC}"
    else
        echo -e "    ${RED}❌ Failed to generate receipt for ${student}${NC}"
    fi
done

# Step 6: Show summary
echo -e "\n${BLUE}📊 Setup Summary:${NC}"
echo -e "${BLUE}==================${NC}"

# Count generated data
student_count=$(find data/generated_student_data/students/ -name "*.json" 2>/dev/null | wc -l | tr -d ' ')
term_count=$(find data/generated_student_data/terms/ -name "*.json" 2>/dev/null | wc -l | tr -d ' ')
merkle_count=$(find data/merkle_trees/ -name "*.json" 2>/dev/null | wc -l | tr -d ' ')
root_count=$(find blockchain_ready/roots/ -name "*.json" 2>/dev/null | wc -l | tr -d ' ')
receipt_count=$(find receipts/ -name "*.json" 2>/dev/null | wc -l | tr -d ' ')

echo -e "👥 Students generated: ${GREEN}${student_count}${NC}"
echo -e "📚 Terms processed: ${GREEN}${term_count}${NC}"  
echo -e "🌳 Merkle trees: ${GREEN}${merkle_count}${NC}"
echo -e "🔗 Blockchain roots: ${GREEN}${root_count}${NC}"
echo -e "🎓 Receipts generated: ${GREEN}${receipt_count}${NC}"

# Show students list
if [ -d "data/generated_student_data/students/" ]; then
    echo -e "\n${BLUE}📋 Generated Students:${NC}"
    ls data/generated_student_data/students/ | sed 's/journey_//g' | sed 's/.json//g' | sed 's/^/  • /'
fi

# Show receipt files
if [ -d "receipts/" ]; then
    echo -e "\n${BLUE}📄 Generated Receipts:${NC}"
    ls receipts/ | sed 's/^/  • /'
fi

# Show available commands
echo -e "\n${BLUE}🚀 Quick Commands:${NC}"
echo -e "  🔍 Verify receipt: ${YELLOW}go run cmd/*.go verify-local receipts/ITITIU00001_journey.json${NC}"
echo -e "  ⛓️  Publish to blockchain: ${YELLOW}go run cmd/*.go publish-roots Semester_1_2023${NC}"
echo -e "  🌐 Start API server: ${YELLOW}npm run dev${NC}"
echo -e "  🖥️  Start web UI: ${YELLOW}cd web/iumicert-issuer && npm run dev${NC}"

echo -e "\n${GREEN}🎉 Quick setup completed successfully!${NC}"
echo -e "${GREEN}You now have 5 students with complete academic journeys and receipts ready for blockchain publishing.${NC}"