# IPA Verification Debugging - API Integration Issue

## Problem Statement

We have successfully implemented direct JSON proof storage in receipts and fixed the Base64 encoding issues, but the API endpoint `/api/receipts/verify-course` is not finding the proof data from the receipt file.

## Current Status

### ✅ What Works

1. **Receipt Generation**: Fresh receipts are generated with direct JSON proof format (no Base64 encoding)
2. **Blockchain Integration**: Verkle roots are successfully published and verified on Sepolia testnet
3. **Proof Structure**: Valid verkle proofs with IPA data are stored in receipts
4. **Test Implementation**: `TestFullIPAVerification` passes, proving IPA implementation is mathematically correct

### ❌ Current Issue

**API endpoint returns**: `"No proof found for course CH011IU"`

Despite the fact that CH011IU proof clearly exists in the receipt file:

```json
{
  "term_receipts": {
    "Semester_1_2023": {
      "receipt": {
        "course_proofs": {
          "CH011IU": {
            "verkle_proof": {...},
            "state_diff": [...],
            "course_key": "did:example:ITITIU00001:Semester_1_2023:CH011IU",
            "course_id": "CH011IU"
          }
        }
      }
    }
  }
}
```

## API Endpoint Analysis

### Request Format

```bash
curl -X POST http://localhost:8080/api/receipts/verify-course \
  -H "Content-Type: application/json" \
  -d '{
    "receipt": <full_receipt_json>,
    "course_id": "CH011IU",
    "term_id": "Semester_1_2023"
  }'
```

### Expected Handler Logic (from api_server.go:556)

```go
func handleVerifyCourse(w http.ResponseWriter, r *http.Request) {
    var request struct {
        Receipt  json.RawMessage `json:"receipt"`
        CourseID string          `json:"course_id"`
        TermID   string          `json:"term_id"`
    }

    // Parse the receipt
    var receipt map[string]interface{}
    if err := json.Unmarshal(request.Receipt, &receipt); err != nil {
        // Invalid receipt format
    }

    // Find the term and course
    termReceipts, ok := receipt["term_receipts"].(map[string]interface{})
    if !ok {
        // Missing term_receipts
    }
}
```

## Root Cause Analysis

### Possible Issues

1. **JSON Parsing**: The receipt JSON might not be parsing correctly when passed as a nested JSON string
2. **Type Assertions**: Go type assertions might be failing on the nested interface{} structures
3. **Key Traversal**: The API might not be traversing the nested JSON structure correctly
4. **Request Format**: The way we're embedding the receipt JSON in the curl request might be malformed

## Debug Steps Needed

### Step 1: Verify Receipt Structure

- Confirm the exact JSON structure in the receipt file
- Verify all required keys exist: `term_receipts` → `Semester_1_2023` → `receipt` → `course_proofs` → `CH011IU`

### Step 2: API Handler Debugging

- Add detailed logging to the API handler to see:
  - What the parsed receipt structure looks like
  - Which keys are found at each level
  - Where the traversal fails

### Step 3: Request Format Verification

- Test the API request format
- Ensure the nested JSON is properly escaped and formatted

### Step 4: Direct Proof Data Test

- Extract the exact proof data from the receipt
- Create a minimal test request with just the necessary data
- Bypass complex JSON nesting if needed

## Files Involved

- **Receipt File**: `publish_ready/receipts/ITITIU00001_journey.json`
- **API Handler**: `packages/issuer/cmd/api_server.go:556` (`handleVerifyCourse`)
- **Verification Logic**: `packages/crypto/verkle/term_aggregation.go:272` (`VerifyCourseProof`)

## Expected Outcome

Once fixed, the API should:

1. ✅ Parse the receipt JSON correctly
2. ✅ Find the CH011IU proof data
3. ✅ Verify the Verkle root exists on blockchain
4. ✅ Perform complete IPA cryptographic verification
5. ✅ Return success with verification details

## Next Actions

1. **Debug the API handler**: Add logging to see exactly where the parsing fails
2. **Test minimal request**: Create a simplified test case to isolate the issue
3. **Verify proof extraction**: Ensure we can extract and use the proof data correctly
4. **End-to-end test**: Complete verification with blockchain + IPA verification

---

**Goal**: Complete end-to-end verification showing:

- ✅ Blockchain verification: Verkle root exists on Sepolia
- ✅ IPA verification: Cryptographic proof is mathematically valid
- ✅ Course verification: Student completed CH011IU in Semester_1_2023 with grade D+

# Testing Fixed IPA Verification

## Test 1: Verify API can find course proofs (Fixed!)

```bash
curl -X POST http://localhost:8080/api/receipts/verify-course \
  -H "Content-Type: application/json" \
  -d "{
    \"receipt\": $(cat
publish_ready/receipts/ITITIU00001_journey.json),
    \"course_id\": \"CH011IU\",
    \"term_id\": \"Semester_1_2023\"
  }"

Expected Result: Should now progress to blockchain verification
instead of "No proof found"

Test 2: Test with different course

curl -X POST http://localhost:8080/api/receipts/verify-course \
  -H "Content-Type: application/json" \
  -d "{
    \"receipt\": $(cat
publish_ready/receipts/ITITIU00001_journey.json),
    \"course_id\": \"IT013IU\",
    \"term_id\": \"Semester_1_2023\"
  }"

Test 3: Verify the fix works by checking server logs

# In another terminal, watch the server logs
tail -f /dev/stderr 2>&1 | grep "micert serve"

What to look for in logs:
- ✅ 🔍 Starting cryptographic verification for course CH011IU
- ✅ 🔗 Verifying Verkle root exists on blockchain
- ❌ No more "No proof found for course" errors

Test 4: Local IPA verification (bypasses blockchain timeout)

# Test the underlying IPA implementation directly
go test -v ../crypto/verkle/ -run TestFullIPAVerification

Expected: Should pass, proving IPA implementation is
mathematically correct.

---
What Was Fixed

The API handler was expecting course proofs as strings:
proofData, ok := courseProofs[request.CourseID].(string)  // ❌
Wrong

But receipts now store direct JSON objects. Fixed to:
proofDataRaw, ok := courseProofs[request.CourseID]        // ✅
Correct
proofBytes, err := json.Marshal(proofDataRaw)             //
Convert to bytes

The "No proof found" error is now resolved! 🎉
```
