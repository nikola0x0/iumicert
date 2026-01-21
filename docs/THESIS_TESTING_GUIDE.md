# IU-MiCert Thesis Testing Guide

## Complete Workflow Using Dashboards + Screenshot Capture

This guide walks through testing the complete system using the web dashboards while capturing measurements and screenshots for Chapter 6.

---

## Prerequisites

1. **Generate base data first:**
   ```bash
   cd /Users/nikola/Developer/iumicert/packages/issuer
   ./reset.sh && ./generate.sh
   ```
   This creates students, terms, and Verkle trees.

2. **Start the API server:**
   ```bash
   cd /Users/nikola/Developer/iumicert/packages/issuer
   ./dev.sh
   # Server runs on http://localhost:8080
   ```

3. **Have your wallet ready:**
   - Metamask installed
   - Connected to Sepolia testnet
   - Issuer private key imported
   - Some Sepolia ETH for gas

4. **Access the dashboards:**
   - Issuer Dashboard: https://iumicert-issuer.vercel.app
   - Student/Verifier Portal: https://iu-micert.vercel.app

---

## Part 1: Credential Issuance Workflow (Screenshots 6.1, 6.2, 6.7, 6.8, 6.9)

### Step 1.1: Main Dashboard
📸 **Screenshot 6.1** - Issuer Dashboard Main Page

1. Open https://iumicert-issuer.vercel.app
2. You should see list of terms with their publication status
3. **Capture screenshot** showing:
   - Term list (Semester_1_2023, Semester_2_2023, etc.)
   - Publication status (unpublished/published)
   - Term root hashes (if available)
   - Action buttons

**Measurements to record:**
- Number of terms displayed: _____
- Terms ready for publishing: _____

### Step 1.2: Demo Data Generation (Optional)
If you need to regenerate data from the dashboard:
1. Click "Demo Data" tab
2. Click generate button
3. Wait for completion

### Step 1.3: Publish Term Root
📸 **Screenshot 6.2** - Publish Term Root Interface
📸 **Screenshot 6.8** - Blockchain Publication

1. Select "Semester_1_2023" from the term list
2. Click "Publish to Blockchain"
3. **Capture Screenshot 6.2** showing:
   - Term identifier
   - Computed Verkle root (32-byte hash)
   - Gas estimation
   - Sepolia testnet indicator

4. Metamask will pop up:
   - **Capture Screenshot 6.8** showing Metamask confirmation
   - Note the gas estimate shown
   - Confirm transaction

5. Wait for confirmation (15-30 seconds)
6. Dashboard should update with transaction hash

**Measurements to record:**
- Transaction hash: _____
- Gas used (from Etherscan): _____
- Block number: _____
- Block timestamp: _____
- Confirmation time: _____ seconds

7. Visit Etherscan: `https://sepolia.etherscan.io/tx/[YOUR_TX_HASH]`
   - Record exact gas used
   - Record transaction cost in ETH

### Step 1.4: Receipt Generation
📸 **Screenshot 6.9** - Receipt Generation Progress

After publishing, the system auto-generates receipts. If using CLI:

```bash
cd /Users/nikola/Developer/iumicert/packages/issuer/cmd
go run . generate-receipt ITITIU00001
```

**Capture screenshot** if dashboard shows progress, otherwise show terminal output.

**Measurements to record:**
- Time to generate receipt: _____ ms
- Receipt file size: `ls -lh ../publish_ready/receipts/ITITIU00001_journey.json`
- Number of courses in receipt: `grep -o '"course_id"' ../publish_ready/receipts/ITITIU00001_journey.json | wc -l`
- Number of terms in receipt: _____

---

## Part 2: Verification Workflow (Screenshots 6.4, 6.5, 6.6, 6.10-6.14)

### Step 2.1: Receipt Upload
📸 **Screenshot 6.4** & **Screenshot 6.10** - Receipt Upload

1. Open https://iu-micert.vercel.app
2. You should see upload interface
3. **Capture screenshot** showing:
   - Drag-and-drop zone
   - Upload button
   - Instructions

### Step 2.2: Upload Receipt
📸 **Screenshot 6.12** - Receipt Upload in Action

1. Upload the receipt file: `ITITIU00001_journey.json`
2. **Capture screenshot** during file validation
3. System should show receipt preview

**Measurements to record:**
- File validation time: _____ ms
- Number of terms detected: _____
- Number of courses detected: _____

### Step 2.3: Verification Process
📸 **Screenshot 6.13** - Verification in Progress
📸 **Screenshot 6.5** - Verification Success
📸 **Screenshot 6.14** - Final Result

1. Click "Verify" button
2. **Capture Screenshot 6.13** showing verification progress:
   - Blockchain queries
   - Proof verifications
   - Progress indicators

3. When complete, **Capture Screenshot 6.5** showing:
   - Green checkmark / success indicator
   - List of verified courses
   - Course details (ID, name, grade, credits)
   - Blockchain confirmation

4. Click to expand details
5. **Capture Screenshot 6.6** and **6.14** showing:
   - Per-course verification times
   - Proof sizes
   - Blockchain query results
   - Term version information

**Measurements to record:**
- Total verification time: _____ ms
- Blockchain query time: _____ ms
- Per-course verification times: _____ ms each
- Number of courses verified: _____
- All courses passed: Yes/No

---

## Part 3: Selective Disclosure Test (Screenshot 6.11)

### Step 3.1: Create Filtered Receipt
📸 **Screenshot 6.11** - Receipt Filtering

1. Open `ITITIU00001_journey.json` in text editor (VS Code recommended)
2. **Capture screenshot** showing the JSON structure
3. Remove some courses (delete the course objects and their proofs)
4. Keep only 3 courses (e.g., IT013IU, IT153IU, IT254IU)
5. Save as `ITITIU00001_filtered.json`

**Measurements to record:**
- Original file size: _____ KB
- Filtered file size: _____ KB
- Original course count: _____
- Filtered course count: 3

### Step 3.2: Verify Filtered Receipt

1. Upload `ITITIU00001_filtered.json` to verifier portal
2. Click Verify
3. **Confirm it still passes verification** ✅

**Measurement:**
- Filtered receipt verification: Success/Fail
- Privacy preserved: _____ courses NOT revealed

---

## Part 4: Revocation Test (Screenshots 6.15-6.18)

### Step 4.1: Submit Revocation Request
📸 **Screenshot 6.15** - Revocation Request

1. Go to Issuer Dashboard → Revocations tab
2. Click "Submit Revocation Request"
3. Fill in:
   - Student ID: ITITIU00003
   - Course ID: IT254IU
   - Term: Semester_1_2024
   - Reason: "Academic misconduct"
4. **Capture screenshot** of the form
5. Submit

### Step 4.2: Approve Revocation
📸 **Screenshot 6.16** - Revocation Approval

1. You should see the pending request in table
2. **Capture screenshot** showing pending request
3. Click "Approve" button

### Step 4.3: Process Revocation
📸 **Screenshot 6.17** - Tree Reconstruction

1. Click "Process Revocations" button
2. System will rebuild tree
3. **Capture screenshot** if progress is shown

**Measurements to record:**
- Tree rebuild time: _____ ms
- Original course count: _____
- New course count: _____
- New root hash: _____

### Step 4.4: Publish New Version
📸 **Screenshot 6.18** - Version Update

1. System publishes v2 root to blockchain
2. **Capture screenshot** showing:
   - v1 root (original)
   - v2 root (after revocation)
   - Transaction confirmation

**Measurements to record:**
- Gas used for v2 publication: _____
- Transaction hash: _____
- Version increment: v1 → v2

### Step 4.5: Test Superseded Receipt

1. Try to verify the old ITITIU00003 receipt (v1)
2. System should show warning: "Newer version available"

**Measurement:**
- Warning displayed: Yes/No

---

## Part 5: Additional Performance Measurements

### Timing Benchmarks (Need CLI for Precision)

```bash
cd /Users/nikola/Developer/iumicert/packages/issuer/cmd

# Test tree insertion timing
go run . test-verify  # This runs comprehensive tests with timing

# Or measure individual operations
time go run . add-term TestTerm_2026 ../test_data.json
```

### Proof Size Inspection

```bash
cd /Users/nikola/Developer/iumicert/packages/issuer

# Get proof size from receipt
# Open receipt and look at "verkle_proof" field size
# Count bytes in the proof structure
```

---

## Part 6: Screenshot Checklist

### Issuer Dashboard Screenshots:
- [ ] 6.1 - Main dashboard page
- [ ] 6.2 - Publish term root interface
- [ ] 6.3 - Revocation management page
- [ ] 6.7 - Tree construction progress
- [ ] 6.8 - Blockchain publication (Metamask)
- [ ] 6.9 - Receipt generation progress
- [ ] 6.15 - Revocation request form
- [ ] 6.16 - Revocation approval interface
- [ ] 6.17 - Tree rebuild progress
- [ ] 6.18 - Version update confirmation

### Student/Verifier Portal Screenshots:
- [ ] 6.4 - Receipt upload interface
- [ ] 6.5 - Verification success result
- [ ] 6.6 - Detailed verification view
- [ ] 6.10 - Student receipt download
- [ ] 6.11 - JSON filtering (text editor)
- [ ] 6.12 - Upload in progress
- [ ] 6.13 - Verification in progress
- [ ] 6.14 - Final result display

---

## Part 7: Data Collection Sheet

Create a file `chapter6_measurements.md` to track everything:

```markdown
# Chapter 6 Measurements

## Proof Sizes
- Single course proof: _____
- Full receipt (20 courses): _____
- Filtered receipt (3 courses): _____

## Timing Benchmarks
- Tree insertion (per course): _____ ms
- Tree commitment (all courses): _____ ms
- Proof generation (per course): _____ ms
- Proof verification (per course): _____ ms
- Full receipt verification (20 courses): _____ ms

## Gas Costs
- Publish term root: _____ gas
- ETH cost: _____ ETH
- USD cost (at $2000/ETH): $_____
- Transaction hash: _____

## Scenario 1: Issuance
- Students: 5
- Total courses: _____
- Tree construction time: _____ ms
- Receipt generation time: _____ s
- Total cost: $_____

## Scenario 2: Verification
- Courses verified: _____
- Verification time: _____ ms
- Privacy preserved: _____ courses hidden

## Scenario 3: Revocation
- Revocation process time: _____ ms
- v2 gas cost: _____ gas
- Impacted students: _____
```

---

## Quick Start Command

**Run this to start everything:**

```bash
# Terminal 1: Generate data
cd /Users/nikola/Developer/iumicert/packages/issuer
./reset.sh && ./generate.sh

# Terminal 2: Start API server
cd /Users/nikola/Developer/iumicert/packages/issuer
./dev.sh

# Then open dashboards in browser and follow steps above
```

---

## Tips for Screenshots

1. Use **full browser window** for cleaner screenshots
2. **Zoom to 100%** for consistency
3. Use **Cmd+Shift+4** on Mac to capture specific regions
4. Name files clearly: `screenshot_6_1_main_dashboard.png`
5. Capture during actual operations for authenticity

---

Ready to start? Let me know when you begin and I can help you record the measurements!
