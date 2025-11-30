package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// BlockchainIntegration provides methods for publishing term roots to IUMiCertRegistry
type BlockchainIntegration struct {
	client           *BlockchainClient
	contractAddress  common.Address
	registryContract *IUMiCertRegistry
}

// PublishResult contains the result of publishing a term root
type PublishResult struct {
	TransactionHash string    `json:"transaction_hash"`
	BlockNumber     uint64    `json:"block_number"`
	GasUsed         uint64    `json:"gas_used"`
	Status          string    `json:"status"`
	PublishedAt     time.Time `json:"published_at"`
}

// NewBlockchainIntegration creates a new blockchain integration instance
func NewBlockchainIntegration(network, privateKeyHex, contractAddressHex string) (*BlockchainIntegration, error) {
	// Create blockchain client
	client, err := NewBlockchainClient(network, privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to create blockchain client: %w", err)
	}

	// Parse contract address
	if !strings.HasPrefix(contractAddressHex, "0x") {
		contractAddressHex = "0x" + contractAddressHex
	}
	contractAddress := common.HexToAddress(contractAddressHex)

	// Create contract instance
	registryContract, err := NewIUMiCertRegistry(contractAddress, client.GetClient())
	if err != nil {
		return nil, fmt.Errorf("failed to create registry contract instance: %w", err)
	}

	return &BlockchainIntegration{
		client:           client,
		contractAddress:  contractAddress,
		registryContract: registryContract,
	}, nil
}

// PublishTermRoot publishes a term root to the blockchain
func (bi *BlockchainIntegration) PublishTermRoot(ctx context.Context, verkleRootHex, termID string, totalStudents *big.Int) (*PublishResult, error) {
	// Parse verkle root
	if !strings.HasPrefix(verkleRootHex, "0x") {
		verkleRootHex = "0x" + verkleRootHex
	}
	
	verkleRootBytes := common.FromHex(verkleRootHex)
	if len(verkleRootBytes) != 32 {
		return nil, fmt.Errorf("verkle root must be 32 bytes, got %d", len(verkleRootBytes))
	}

	var verkleRoot [32]byte
	copy(verkleRoot[:], verkleRootBytes)

	// Get transaction options
	auth, err := bi.client.GetTransactOpts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction options: %w", err)
	}

	// Call publishTermRoot function
	tx, err := bi.registryContract.PublishTermRoot(auth, verkleRoot, termID, totalStudents)
	if err != nil {
		return nil, fmt.Errorf("failed to publish term root: %w", err)
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(ctx, bi.client.GetClient(), tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction mining: %w", err)
	}

	// Check if transaction was successful
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("transaction failed with status %d", receipt.Status)
	}

	result := &PublishResult{
		TransactionHash: tx.Hash().Hex(),
		BlockNumber:     receipt.BlockNumber.Uint64(),
		GasUsed:         receipt.GasUsed,
		Status:          "success",
		PublishedAt:     time.Now(),
	}

	return result, nil
}

// VerifyReceiptAnchor verifies if a receipt's blockchain anchor is valid
func (bi *BlockchainIntegration) VerifyReceiptAnchor(ctx context.Context, blockchainAnchorHex string) (bool, string, *big.Int, error) {
	// Parse blockchain anchor
	if !strings.HasPrefix(blockchainAnchorHex, "0x") {
		blockchainAnchorHex = "0x" + blockchainAnchorHex
	}
	
	anchorBytes := common.FromHex(blockchainAnchorHex)
	if len(anchorBytes) != 32 {
		return false, "", nil, fmt.Errorf("blockchain anchor must be 32 bytes, got %d", len(anchorBytes))
	}

	var anchor [32]byte
	copy(anchor[:], anchorBytes)

	// Call verifyReceiptAnchor function
	callOpts := bi.client.GetCallOpts(ctx)
	result, err := bi.registryContract.VerifyReceiptAnchor(callOpts, anchor)
	if err != nil {
		return false, "", nil, fmt.Errorf("failed to verify receipt anchor: %w", err)
	}

	return result.IsValid, result.TermId, result.PublishedAt, nil
}

// GetLatestRootForTerm gets the latest version info for a term
func (bi *BlockchainIntegration) GetLatestRootForTerm(ctx context.Context, termID string) (*LatestRootInfo, error) {
	callOpts := bi.client.GetCallOpts(ctx)
	result, err := bi.registryContract.GetLatestRoot(callOpts, termID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest root: %w", err)
	}

	return &LatestRootInfo{
		RootHash:      result.RootHash,
		Version:       result.Version,
		TotalStudents: result.TotalStudents,
		PublishedAt:   result.PublishedAt,
	}, nil
}

// CheckRootStatus checks the validity status of a verkle root
func (bi *BlockchainIntegration) CheckRootStatus(ctx context.Context, verkleRootHex string) (*RootStatus, error) {
	// Parse verkle root
	if !strings.HasPrefix(verkleRootHex, "0x") {
		verkleRootHex = "0x" + verkleRootHex
	}

	verkleRootBytes := common.FromHex(verkleRootHex)
	if len(verkleRootBytes) != 32 {
		return nil, fmt.Errorf("verkle root must be 32 bytes, got %d", len(verkleRootBytes))
	}

	var verkleRoot [32]byte
	copy(verkleRoot[:], verkleRootBytes)

	callOpts := bi.client.GetCallOpts(ctx)
	result, err := bi.registryContract.CheckRootStatus(callOpts, verkleRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to check root status: %w", err)
	}

	return &RootStatus{
		Status:     result.Status,
		TermID:     result.TermId,
		Version:    result.Version,
		LatestRoot: result.LatestRoot,
		Message:    result.Message,
	}, nil
}

// GetPublishedRootsCount gets the total number of published roots
func (bi *BlockchainIntegration) GetPublishedRootsCount(ctx context.Context) (*big.Int, error) {
	callOpts := bi.client.GetCallOpts(ctx)
	count, err := bi.registryContract.GetPublishedRootsCount(callOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to get published roots count: %w", err)
	}
	return count, nil
}

// GetPublishedRoot gets a published root by index
func (bi *BlockchainIntegration) GetPublishedRoot(ctx context.Context, index *big.Int) ([32]byte, error) {
	callOpts := bi.client.GetCallOpts(ctx)
	root, err := bi.registryContract.GetPublishedRoot(callOpts, index)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to get published root at index %s: %w", index.String(), err)
	}
	return root, nil
}

// Close closes the blockchain client connection
func (bi *BlockchainIntegration) Close() {
	if bi.client != nil {
		bi.client.Close()
	}
}

// LatestRootInfo contains information about the latest version of a term root
type LatestRootInfo struct {
	RootHash      [32]byte `json:"root_hash"`
	Version       *big.Int `json:"version"`
	TotalStudents *big.Int `json:"total_students"`
	PublishedAt   *big.Int `json:"published_at"`
}

// RootStatus contains status information for a verkle root
type RootStatus struct {
	Status     uint8    `json:"status"`      // 0=Invalid, 1=Current, 2=Outdated, 3=Superseded
	TermID     string   `json:"term_id"`
	Version    *big.Int `json:"version"`
	LatestRoot [32]byte `json:"latest_root"`
	Message    string   `json:"message"`
}

// PublishTermRootFromFile publishes a term root from a JSON file
func (bi *BlockchainIntegration) PublishTermRootFromFile(ctx context.Context, rootFilePath string) (*PublishResult, error) {
	// Read root data from file
	data, err := os.ReadFile(rootFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read root file: %w", err)
	}

	var rootData map[string]interface{}
	if err := json.Unmarshal(data, &rootData); err != nil {
		return nil, fmt.Errorf("failed to parse root data: %w", err)
	}

	// Extract required fields
	verkleRootHex, ok := rootData["verkle_root"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid verkle_root in root data")
	}

	termID, ok := rootData["term_id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid term_id in root data")
	}

	totalStudentsFloat, ok := rootData["total_students"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid total_students in root data")
	}
	totalStudents := big.NewInt(int64(totalStudentsFloat))

	// Publish to blockchain
	result, err := bi.PublishTermRoot(ctx, verkleRootHex, termID, totalStudents)
	if err != nil {
		return nil, fmt.Errorf("failed to publish term root: %w", err)
	}

	// Save transaction record
	if err := bi.saveTransactionRecord(result, rootFilePath); err != nil {
		// Log the error but don't fail the operation
		fmt.Printf("Warning: failed to save transaction record: %v\n", err)
	}

	return result, nil
}

// saveTransactionRecord saves the transaction record to the publish_ready/transactions directory
func (bi *BlockchainIntegration) saveTransactionRecord(result *PublishResult, rootFilePath string) error {
	// Create transaction record
	txRecord := map[string]interface{}{
		"transaction_hash": result.TransactionHash,
		"block_number":     result.BlockNumber,
		"gas_used":         result.GasUsed,
		"status":           result.Status,
		"published_at":     result.PublishedAt,
		"root_file_path":   rootFilePath,
		"contract_address": bi.contractAddress.Hex(),
	}

	// Ensure transactions directory exists
	txDir := "publish_ready/transactions"
	if err := os.MkdirAll(txDir, 0755); err != nil {
		return fmt.Errorf("failed to create transactions directory: %w", err)
	}

	// Create filename based on transaction hash
	filename := fmt.Sprintf("tx_%s_%s.json", 
		result.TransactionHash[2:10], // First 8 chars after 0x
		result.PublishedAt.Format("20060102_150405"))
	
	filepath := filepath.Join(txDir, filename)

	// Marshal and save
	txData, err := json.MarshalIndent(txRecord, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal transaction record: %w", err)
	}

	if err := os.WriteFile(filepath, txData, 0644); err != nil {
		return fmt.Errorf("failed to write transaction record: %w", err)
	}

	return nil
}

// SupersedeTerm publishes a new version of a term root (for revocation)
func (bi *BlockchainIntegration) SupersedeTerm(ctx context.Context, termID string, newVerkleRootHex string, totalStudents *big.Int, reason string) (*PublishResult, error) {
	// Parse new verkle root
	if !strings.HasPrefix(newVerkleRootHex, "0x") {
		newVerkleRootHex = "0x" + newVerkleRootHex
	}

	verkleRootBytes := common.FromHex(newVerkleRootHex)
	if len(verkleRootBytes) != 32 {
		return nil, fmt.Errorf("verkle root must be 32 bytes, got %d", len(verkleRootBytes))
	}

	var newVerkleRoot [32]byte
	copy(newVerkleRoot[:], verkleRootBytes)

	// Get transaction options
	auth, err := bi.client.GetTransactOpts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction options: %w", err)
	}

	// Call supersedeTerm function
	tx, err := bi.registryContract.SupersedeTerm(auth, termID, newVerkleRoot, totalStudents, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to supersede term: %w", err)
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(ctx, bi.client.GetClient(), tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction mining: %w", err)
	}

	// Check if transaction was successful
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("transaction failed with status %d", receipt.Status)
	}

	result := &PublishResult{
		TransactionHash: tx.Hash().Hex(),
		BlockNumber:     receipt.BlockNumber.Uint64(),
		GasUsed:         receipt.GasUsed,
		Status:          "success",
		PublishedAt:     time.Now(),
	}

	return result, nil
}

// GetTermHistory gets complete version history for a term
func (bi *BlockchainIntegration) GetTermHistory(ctx context.Context, termID string) ([]uint, [][32]byte, error) {
	callOpts := bi.client.GetCallOpts(ctx)
	result, err := bi.registryContract.GetTermHistory(callOpts, termID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get term history: %w", err)
	}

	// Convert *big.Int versions to uint
	versions := make([]uint, len(result.Versions))
	for i, v := range result.Versions {
		versions[i] = uint(v.Uint64())
	}

	return versions, result.Roots, nil
}

// LoadPrivateKeyFromEnv loads private key from environment variable
func LoadPrivateKeyFromEnv() (*ecdsa.PrivateKey, error) {
	privateKeyHex := os.Getenv("ISSUER_PRIVATE_KEY")
	if privateKeyHex == "" {
		return nil, fmt.Errorf("ISSUER_PRIVATE_KEY environment variable not set")
	}

	// Remove 0x prefix if present
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	return crypto.HexToECDSA(privateKeyHex)
}

// GeneratePrivateKeyHex generates a new private key for testing
func GeneratePrivateKeyHex() (string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate private key: %w", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	return hex.EncodeToString(privateKeyBytes), nil
}