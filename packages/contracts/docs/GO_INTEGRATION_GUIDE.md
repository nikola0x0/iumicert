# Go Integration Guide for Updated IUMiCertRegistry

## Overview
The updated contract adds versioning support. Your Go code needs to call new functions.

## New Smart Contract Functions to Integrate

### 1. supersedeTerm (for revocation)
```solidity
function supersedeTerm(
    string memory _termId,
    bytes32 _newVerkleRoot,
    uint256 _newTotalStudents,
    string memory _reason
) external onlyOwner
```

### 2. checkRootStatus (enhanced verification)
```solidity
function checkRootStatus(bytes32 _verkleRoot)
    external view
    returns (
        uint8 status,           // 0=Invalid, 1=Current, 2=Outdated, 3=Superseded
        string memory termId,
        uint256 version,
        bytes32 latestRoot,
        string memory message
    )
```

### 3. getLatestRoot (get current version)
```solidity
function getLatestRoot(string memory _termId)
    external view
    returns (
        bytes32 rootHash,
        uint256 version,
        uint256 totalStudents,
        uint256 publishedAt
    )
```

### 4. getTermHistory (version tracking)
```solidity
function getTermHistory(string memory _termId)
    external view
    returns (
        uint256[] memory versions,
        bytes32[] memory roots
    )
```

## Required Go Code Updates

### Step 1: Generate New Go Bindings

```bash
cd packages/issuer

# Generate ABI from compiled contract
forge inspect IUMiCertRegistry abi > blockchain_integration/IUMiCertRegistry.abi

# Generate Go bindings using abigen
abigen --abi=blockchain_integration/IUMiCertRegistry.abi \
       --pkg=blockchain \
       --type=IUMiCertRegistry \
       --out=blockchain_integration/iumicert_registry.go
```

### Step 2: Update blockchain_integration/client.go

Add these new functions:

```go
// SupersedeTerm publishes a new version of a term root (for revocation)
func (c *Client) SupersedeTerm(
    termID string,
    newRoot [32]byte,
    totalStudents *big.Int,
    reason string,
) (*types.Transaction, error) {
    auth, err := c.getAuth()
    if err != nil {
        return nil, err
    }

    tx, err := c.registry.SupersedeTerm(auth, termID, newRoot, totalStudents, reason)
    if err != nil {
        return nil, fmt.Errorf("failed to supersede term: %w", err)
    }

    return tx, nil
}

// CheckRootStatus checks the validity status of a root
func (c *Client) CheckRootStatus(root [32]byte) (
    status uint8,
    termID string,
    version *big.Int,
    latestRoot [32]byte,
    message string,
    err error,
) {
    callOpts := &bind.CallOpts{Context: context.Background()}
    return c.registry.CheckRootStatus(callOpts, root)
}

// GetLatestRoot returns the latest version for a term
func (c *Client) GetLatestRoot(termID string) (
    rootHash [32]byte,
    version *big.Int,
    totalStudents *big.Int,
    publishedAt *big.Int,
    err error,
) {
    callOpts := &bind.CallOpts{Context: context.Background()}
    return c.registry.GetLatestRoot(callOpts, termID)
}

// GetTermHistory returns all versions for a term
func (c *Client) GetTermHistory(termID string) (
    versions []*big.Int,
    roots [][32]byte,
    err error,
) {
    callOpts := &bind.CallOpts{Context: context.Background()}
    return c.registry.GetTermHistory(callOpts, termID)
}
```

### Step 3: Update cmd/verify.go

Enhance verification to check root status:

```go
func verifyReceipt(receiptPath string) error {
    // ... existing code to load receipt ...

    for termID, termReceipt := range receipt.TermReceipts {
        termRoot := termReceipt.VerkleRoot
        
        // NEW: Check root status
        status, _, version, latestRoot, message, err := blockchainClient.CheckRootStatus(termRoot)
        if err != nil {
            log.Warnf("Could not check root status for %s: %v", termID, err)
            // Continue with verification anyway
        } else {
            switch status {
            case 0: // Invalid
                return fmt.Errorf("❌ Root not found in registry for term %s", termID)
            case 1: // Valid & Current
                log.Infof("✅ Using current version (v%s) for %s", version, termID)
            case 2: // Valid but Outdated
                log.Warnf("⚠️  Using old version (v%s) for %s. Latest: v%s", 
                    version, termID, latestRoot)
                log.Warnf("    Receipt is valid but student should download updated version")
            case 3: // Superseded
                return fmt.Errorf("❌ Root superseded for term %s\nReason: %s\n"+
                    "Student must download updated receipt", termID, message)
            }
        }

        // ... continue with existing Verkle proof verification ...
    }

    return nil
}
```

### Step 4: Add New CLI Commands

#### cmd/revoke.go (NEW FILE)
```go
package main

import (
    "crypto/sha256"
    "fmt"
    "github.com/spf13/cobra"
)

var revokeCmd = &cobra.Command{
    Use:   "revoke-and-supersede [termID] [reason]",
    Short: "Revoke credentials and publish new term root",
    Long: `Rebuild Verkle tree excluding revoked credentials and publish new version.
    
Example:
  ./micert revoke-and-supersede Semester_1_2023 "Revoked: Student 00001 fraudulent IT089IU"`,
    Args: cobra.ExactArgs(2),
    Run: runRevoke,
}

func runRevoke(cmd *cobra.Command, args []string) {
    termID := args[0]
    reason := args[1]

    fmt.Printf("Processing revocations for %s...\n", termID)

    // 1. Load original term data
    termData, err := loadTermData(termID)
    if err != nil {
        log.Fatal(err)
    }

    // 2. Load revocation requests (from data/revocations/termID.json)
    revocations, err := loadRevocationRequests(termID)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d revocation requests\n", len(revocations))

    // 3. Filter out revoked credentials
    filteredData := filterRevokedCredentials(termData, revocations)

    // 4. Rebuild Verkle tree
    newTree, newRoot, err := buildVerkleTree(filteredData)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("New root: %x\n", newRoot)

    // 5. Publish to blockchain
    tx, err := blockchainClient.SupersedeTerm(
        termID,
        newRoot,
        big.NewInt(int64(len(filteredData.Students))),
        reason,
    )
    if err != nil {
        log.Fatalf("Failed to publish: %v", err)
    }

    fmt.Printf("✅ Published! Transaction: %s\n", tx.Hash().Hex())

    // 6. Regenerate receipts for all students
    fmt.Println("Regenerating student receipts...")
    for _, studentID := range filteredData.Students {
        err := generateReceipt(studentID)
        if err != nil {
            log.Warnf("Failed to generate receipt for %s: %v", studentID, err)
        }
    }

    fmt.Println("✅ Revocation complete!")
}

func init() {
    rootCmd.AddCommand(revokeCmd)
}
```

## Testing Checklist

After deploying and updating Go code:

1. ✅ Publish initial term root: `./micert publish-roots Semester_1_2023`
2. ✅ Verify root status: Should show "Current version (v1)"
3. ✅ Create revocation request
4. ✅ Run supersede command
5. ✅ Verify root status: Old root should show "Superseded", new root "Current (v2)"
6. ✅ Old receipt verification should fail or show warning
7. ✅ New receipt verification should pass

## Contract ABI Location

After deployment, get the ABI:
```bash
forge inspect src/IUMiCertRegistry.sol:IUMiCertRegistry abi > IUMiCertRegistry.abi
```

This ABI is needed for:
- Go bindings generation (abigen)
- Frontend integration (ethers.js)
- Direct contract calls (cast)
