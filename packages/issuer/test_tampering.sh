#!/bin/bash

# Test: What happens if we modify course data but keep proofs?

echo "=== Testing Receipt Tampering Detection ==="

# 1. Get a real receipt
RECEIPT="publish_ready/receipts/ITITIU00001_journey.json"

if [ ! -f "$RECEIPT" ]; then
  echo "❌ Receipt not found. Run ./generate.sh first"
  exit 1
fi

# 2. Create a tampered copy
TAMPERED="/tmp/tampered_receipt.json"
cp "$RECEIPT" "$TAMPERED"

# 3. Modify grade in revealed_courses using jq
echo "Original grade:"
jq -r '.term_receipts.Semester_1_2023.receipt.revealed_courses[0].grade' "$RECEIPT"

# Change grade from whatever it is to "FAKE_GRADE"
jq '.term_receipts.Semester_1_2023.receipt.revealed_courses[0].grade = "FAKE_GRADE"' "$RECEIPT" > "$TAMPERED"

echo "Tampered grade:"
jq -r '.term_receipts.Semester_1_2023.receipt.revealed_courses[0].grade' "$TAMPERED"

echo ""
echo "Proofs remain unchanged (same as original)"
echo ""

# 4. Try to verify the tampered receipt
echo "=== Attempting to verify tampered receipt ==="
echo "Expected: VERIFICATION SHOULD FAIL (value hash mismatch)"
echo ""

./micert verify-local "$TAMPERED"

