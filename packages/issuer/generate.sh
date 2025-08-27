#!/bin/bash

echo "🚀 IU-MiCert Generation Script (Single Verkle Architecture)"
echo "============================================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${PURPLE}🌳 Single Verkle Tree System${NC}"
echo -e "${PURPLE}✨ 32-byte proofs • Course-level selective disclosure${NC}"

# Step 1: Generate student academic journeys
echo -e "\n${YELLOW}👥 Step 1: Generating student academic journeys...${NC}"
cd cmd
go run . generate-data
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Student journeys generated (5 students, 6 terms)${NC}"
else
    echo -e "${RED}❌ Failed to generate student data${NC}"
    exit 1
fi
cd ..

# Step 2: Process each term with single Verkle approach
echo -e "\n${YELLOW}🌳 Step 2: Processing terms with single Verkle trees...${NC}"
terms=("Semester_1_2023" "Semester_2_2023" "Summer_2023" "Semester_1_2024" "Semester_2_2024" "Summer_2024")

processed_terms=0
for term in "${terms[@]}"; do
    echo -e "  📚 Processing term: ${CYAN}${term}${NC}"
    
    # Convert data for this term
    cd cmd
    go run . convert-data $term
    if [ $? -eq 0 ]; then
        echo -e "    ✓ Data converted to Verkle format"
    else
        echo -e "    ${RED}❌ Failed to convert data for ${term}${NC}"
        continue
    fi
    
    # Create single Verkle tree for term
    go run . add-term $term ../data/verkle_terms/${term}_completions.json
    if [ $? -eq 0 ]; then
        echo -e "    ${GREEN}✅ Single Verkle tree created for ${term}${NC}"
        ((processed_terms++))
    else
        echo -e "    ${RED}❌ Failed to create Verkle tree for ${term}${NC}"
    fi
    cd ..
done

echo -e "${GREEN}✅ Processed ${processed_terms}/${#terms[@]} terms successfully${NC}"

# Step 3: Generate student receipts
echo -e "\n${YELLOW}🎓 Step 3: Generating student verification receipts...${NC}"
students=("ITITIU00001" "ITITIU00002" "ITITIU00003" "ITITIU00004" "ITITIU00005")

generated_receipts=0
for student in "${students[@]}"; do
    echo -e "  👤 Generating receipt for ${CYAN}${student}${NC}..."
    cd cmd
    go run . generate-receipt $student ../publish_ready/receipts/${student}_journey.json
    if [ $? -eq 0 ]; then
        echo -e "    ${GREEN}✅ Receipt generated with Verkle proofs${NC}"
        ((generated_receipts++))
    else
        echo -e "    ${RED}❌ Failed to generate receipt for ${student}${NC}"
    fi
    cd ..
done

echo -e "${GREEN}✅ Generated ${generated_receipts}/${#students[@]} receipts successfully${NC}"

# Step 4: Show comprehensive summary
echo -e "\n${BLUE}📊 Single Verkle System Summary${NC}"
echo -e "${BLUE}================================${NC}"

# Count all generated data
student_count=$(find data/student_journeys/students/ -name "*.json" 2>/dev/null | wc -l | tr -d ' ')
term_count=$(find data/student_journeys/terms/ -name "*.json" 2>/dev/null | wc -l | tr -d ' ')
verkle_count=$(find data/verkle_terms/ -name "*.json" 2>/dev/null | wc -l | tr -d ' ')
root_count=$(find publish_ready/roots/ -name "*.json" 2>/dev/null | wc -l | tr -d ' ')
receipt_count=$(find publish_ready/receipts/ -name "*.json" 2>/dev/null | wc -l | tr -d ' ')

echo -e "👥 Students: ${GREEN}${student_count}${NC} (with complete academic journeys)"
echo -e "📚 Terms: ${GREEN}${term_count}${NC} (with course completion data)"  
echo -e "🌳 Single Verkle Trees: ${GREEN}${verkle_count}${NC} (pure Verkle, no Merkle)"
echo -e "🔗 Blockchain Roots: ${GREEN}${root_count}${NC} (ready for publishing)"
echo -e "🎓 Verification Receipts: ${GREEN}${receipt_count}${NC} (with 32-byte proofs)"

# Show students with their IDs
if [ -d "data/student_journeys/students/" ]; then
    echo -e "\n${BLUE}📋 Generated Students:${NC}"
    for student_file in data/student_journeys/students/journey_*.json; do
        if [ -f "$student_file" ]; then
            student_id=$(basename "$student_file" .json | sed 's/journey_//')
            courses=$(grep -o '"course_id"' "$student_file" 2>/dev/null | wc -l | tr -d ' ')
            echo -e "  • ${CYAN}${student_id}${NC} - ${courses} total courses across all terms"
        fi
    done
fi

# Show Verkle tree roots
if [ -d "publish_ready/roots/" ]; then
    echo -e "\n${BLUE}🔗 Verkle Tree Roots (Ready for Blockchain):${NC}"
    for root_file in publish_ready/roots/root_*.json; do
        if [ -f "$root_file" ]; then
            term=$(basename "$root_file" .json | sed 's/root_//')
            root_hash=$(grep '"verkle_root"' "$root_file" 2>/dev/null | cut -d'"' -f4 | cut -c1-16)
            echo -e "  • ${CYAN}${term}${NC}: ${root_hash}..."
        fi
    done
fi

# Show architecture benefits
echo -e "\n${PURPLE}🌟 Single Verkle Architecture Benefits:${NC}"
echo -e "  ✨ ${GREEN}32-byte constant proof size${NC} (vs variable Merkle proofs)"
echo -e "  ✨ ${GREEN}Course-level selective disclosure${NC} (reveal specific courses)"
echo -e "  ✨ ${GREEN}Better privacy${NC} (no student data exposure in proofs)"
echo -e "  ✨ ${GREEN}Simplified verification${NC} (single Verkle root per term)"

# Show available commands
echo -e "\n${BLUE}🚀 Ready-to-Use Commands:${NC}"
echo -e "  🔍 Verify receipt: ${YELLOW}cd cmd && go run . verify-local ../publish_ready/receipts/ITITIU00001_journey.json${NC}"
echo -e "  ⛓️  Publish to blockchain: ${YELLOW}cd cmd && go run . publish-roots Semester_1_2023${NC}"
echo -e "  🌐 Start API server: ${YELLOW}cd cmd && go run . serve${NC}"
echo -e "  🖥️  Start web UI: ${YELLOW}cd web/iumicert-issuer && npm run dev${NC}"

# Show next steps
echo -e "\n${CYAN}📋 Next Steps for Thesis Demo:${NC}"
echo -e "  1. ${PURPLE}Test selective disclosure${NC}: Generate receipts with specific courses"
echo -e "  2. ${PURPLE}Blockchain integration${NC}: Publish term roots to Sepolia testnet"
echo -e "  3. ${PURPLE}Frontend demo${NC}: Show course-level proof verification"
echo -e "  4. ${PURPLE}Performance comparison${NC}: Compare with old hybrid approach"

echo -e "\n${GREEN}🎉 Single Verkle system ready!${NC}"
echo -e "${GREEN}Generated ${processed_terms} Verkle trees with ${student_count} students and ${receipt_count} verification receipts.${NC}"