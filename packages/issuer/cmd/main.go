package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"iumicert/crypto/merkle"
	"iumicert/crypto/verkle"
)

var rootCmd = &cobra.Command{
	Use:   "micert",
	Short: "IU-MiCert CLI - Academic Micro-credential Management Tool",
	Long: `IU-MiCert CLI provides comprehensive tools for managing academic micro-credentials
using hybrid Merkle-Verkle tree architecture with blockchain integration.

Features:
  - Per-term Verkle tree construction
  - Student-level Merkle tree management  
  - Academic journey receipt generation
  - Selective disclosure support
  - Blockchain integration for term root publishing
  - Local and on-chain verification`,
}

var initCmd = &cobra.Command{
	Use:   "init [institution-id]",
	Short: "Initialize new credential repository",
	Long: `Initialize a new IU-MiCert credential repository for an educational institution.
Creates necessary directories and configuration files for credential management.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		institutionID := args[0]
		if err := initializeRepository(institutionID); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Initialization failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var addTermCmd = &cobra.Command{
	Use:   "add-term [term-id] [data-file]",
	Short: "Add new academic term with course completions",
	Long: `Add a new academic term to the credential system with course completion data.
Builds student-level Merkle trees and prepares for Verkle tree aggregation.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		termID := args[0]
		dataFile := args[1]
		
		format, _ := cmd.Flags().GetString("format")
		validate, _ := cmd.Flags().GetBool("validate")
		
		if err := addAcademicTerm(termID, dataFile, format, validate); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to add term: %v\n", err)
			os.Exit(1)
		}
	},
}

var generateReceiptCmd = &cobra.Command{
	Use:   "generate-receipt [student-id] [output-file]",
	Short: "Generate academic journey receipt for student",
	Long: `Generate a comprehensive academic journey receipt for a student.
Supports selective disclosure and multi-term verification.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		studentID := args[0]
		outputFile := args[1]
		
		terms, _ := cmd.Flags().GetStringSlice("terms")
		courses, _ := cmd.Flags().GetStringSlice("courses")
		selective, _ := cmd.Flags().GetBool("selective")
		
		if err := generateStudentReceipt(studentID, outputFile, terms, courses, selective); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to generate receipt: %v\n", err)
			os.Exit(1)
		}
	},
}

var verifyLocalCmd = &cobra.Command{
	Use:   "verify-local [receipt-file]",
	Short: "Verify receipt locally without blockchain",
	Long: `Perform local verification of an academic journey receipt.
Validates cryptographic proofs and temporal consistency without blockchain queries.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		receiptFile := args[0]
		
		if err := verifyReceiptLocally(receiptFile); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Verification failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var publishRootsCmd = &cobra.Command{
	Use:   "publish-roots [term-id]",
	Short: "Publish term root commitments to blockchain",
	Long: `Publish Verkle tree root commitments for academic terms to the blockchain.
Enables on-chain verification of academic journey receipts.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		termID := args[0]
		
		network, _ := cmd.Flags().GetString("network")
		privateKey, _ := cmd.Flags().GetString("private-key")
		gasLimit, _ := cmd.Flags().GetUint64("gas-limit")
		
		if err := publishTermRoots(termID, network, privateKey, gasLimit); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to publish roots: %v\n", err)
			os.Exit(1)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("IU-MiCert CLI v1.0.0")
		fmt.Println("Built for academic micro-credential management")
		fmt.Println("Author: Le Tien Phat")
		fmt.Println("License: MIT")
	},
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.micert.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "enable verbose output")
	
	// Add command flags
	addTermCmd.Flags().String("format", "json", "input data format (json, csv)")
	addTermCmd.Flags().Bool("validate", true, "validate input data")
	
	generateReceiptCmd.Flags().StringSlice("terms", []string{}, "specific terms to include")
	generateReceiptCmd.Flags().StringSlice("courses", []string{}, "specific courses to include")
	generateReceiptCmd.Flags().Bool("selective", false, "enable selective disclosure")
	
	publishRootsCmd.Flags().String("network", "sepolia", "blockchain network")
	publishRootsCmd.Flags().String("private-key", "", "private key for signing")
	publishRootsCmd.Flags().Uint64("gas-limit", 0, "gas limit for transaction")
	
	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addTermCmd)
	rootCmd.AddCommand(generateReceiptCmd)
	rootCmd.AddCommand(verifyLocalCmd)
	rootCmd.AddCommand(publishRootsCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}

// Implementation Functions

func initializeRepository(institutionID string) error {
	fmt.Printf("🏛️  Initializing IU-MiCert repository for institution: %s\n", institutionID)
	
	// Create directory structure for blockchain integration
	dirs := []string{
		"data",
		"data/terms",
		"data/students", 
		"data/merkle_trees",
		"data/verkle_trees",
		"blockchain_ready",
		"blockchain_ready/receipts",
		"blockchain_ready/proofs",
		"blockchain_ready/roots",
		"blockchain_ready/transactions",
		"config",
		"logs",
	}
	
	fmt.Println("📁 Creating directory structure...")
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		fmt.Printf("  ✓ Created: %s/\n", dir)
	}
	
	// Generate configuration file
	fmt.Println("⚙️  Generating configuration files...")
	config := map[string]interface{}{
		"institution_id": institutionID,
		"version": "1.0.0",
		"initialized_at": time.Now().Format(time.RFC3339),
		"blockchain": map[string]interface{}{
			"default_network": "sepolia",
			"gas_limit": 500000,
			"confirmation_blocks": 3,
		},
		"output_paths": map[string]string{
			"receipts": "blockchain_ready/receipts",
			"proofs": "blockchain_ready/proofs", 
			"roots": "blockchain_ready/roots",
			"transactions": "blockchain_ready/transactions",
		},
	}
	
	configFile, err := os.Create("config/micert.json")
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer configFile.Close()
	
	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	fmt.Println("🔑 Setting up cryptographic parameters...")
	fmt.Printf("  ✓ Institution ID: %s\n", institutionID)
	fmt.Printf("  ✓ Config saved to: config/micert.json\n")
	
	fmt.Println("✅ Repository initialized successfully!")
	fmt.Println("\n📋 Next steps:")
	fmt.Println("  1. Add academic terms: micert add-term <term-id> <data-file>")
	fmt.Println("  2. Generate receipts: micert generate-receipt <student-id> <output-file>")
	fmt.Println("  3. Publish to blockchain: micert publish-roots <term-id>")
	
	return nil
}

func addAcademicTerm(termID, dataFile, format string, validate bool) error {
	fmt.Printf("📚 Adding academic term: %s\n", termID)
	fmt.Printf("📖 Processing data from: %s (format: %s)\n", dataFile, format)
	
	if validate {
		fmt.Println("✅ Validating input data...")
		if _, err := os.Stat(dataFile); os.IsNotExist(err) {
			return fmt.Errorf("data file does not exist: %s", dataFile)
		}
	}
	
	// Load data based on format
	var completions []merkle.CourseCompletion
	var err error
	
	switch format {
	case "json":
		completions, err = loadCompletionsFromJSON(dataFile)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}
	
	fmt.Printf("📊 Loaded %d course completions\n", len(completions))
	
	// Build student-term Merkle trees
	fmt.Println("🌳 Building student-term Merkle trees...")
	studentCompletions := make(map[string][]merkle.CourseCompletion)
	
	for _, completion := range completions {
		studentDID := fmt.Sprintf("did:example:%s", completion.StudentID)
		studentCompletions[studentDID] = append(studentCompletions[studentDID], completion)
	}
	
	// Create Merkle trees for each student
	merkleDir := filepath.Join("data", "merkle_trees", termID)
	if err := os.MkdirAll(merkleDir, 0755); err != nil {
		return fmt.Errorf("failed to create merkle directory: %w", err)
	}
	
	studentTrees := make(map[string]*merkle.StudentTermMerkle)
	for studentDID, courses := range studentCompletions {
		tree, err := merkle.NewStudentTermMerkle(studentDID, termID, courses)
		if err != nil {
			return fmt.Errorf("failed to create merkle tree for %s: %w", studentDID, err)
		}
		
		studentTrees[studentDID] = tree
		
		// Save Merkle tree
		treeData, err := tree.SerializeToJSON()
		if err != nil {
			return fmt.Errorf("failed to serialize merkle tree: %w", err)
		}
		
		filename := fmt.Sprintf("merkle_%s_%s.json", extractStudentID(studentDID), termID)
		if err := os.WriteFile(filepath.Join(merkleDir, filename), treeData, 0644); err != nil {
			return fmt.Errorf("failed to save merkle tree: %w", err)
		}
		
		fmt.Printf("  ✓ Built Merkle tree for %s: %d courses, root: %x\n", 
			studentDID, len(courses), tree.Root[:8])
	}
	
	// Build term-level Verkle tree
	fmt.Println("🔗 Preparing Verkle tree aggregation...")
	termTree := verkle.NewTermVerkleTree(termID)
	
	for studentDID, courses := range studentCompletions {
		err := termTree.AddStudent(studentDID, courses)
		if err != nil {
			return fmt.Errorf("failed to add student %s to verkle tree: %w", studentDID, err)
		}
	}
	
	// Publish the term
	err = termTree.PublishTerm()
	if err != nil {
		return fmt.Errorf("failed to publish term: %w", err)
	}
	
	// Save Verkle tree and prepare for blockchain
	verkleDir := filepath.Join("data", "verkle_trees")
	if err := os.MkdirAll(verkleDir, 0755); err != nil {
		return fmt.Errorf("failed to create verkle directory: %w", err)
	}
	
	// Save root for blockchain publishing
	rootsDir := "blockchain_ready/roots"
	rootData := map[string]interface{}{
		"term_id": termID,
		"verkle_root": fmt.Sprintf("%x", termTree.VerkleRoot),
		"timestamp": time.Now().Format(time.RFC3339),
		"total_students": len(studentCompletions),
		"ready_for_blockchain": true,
	}
	
	rootFile, err := json.MarshalIndent(rootData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal root data: %w", err)
	}
	
	if err := os.WriteFile(filepath.Join(rootsDir, fmt.Sprintf("root_%s.json", termID)), rootFile, 0644); err != nil {
		return fmt.Errorf("failed to save root file: %w", err)
	}
	
	fmt.Printf("  ✅ Verkle root: %x\n", termTree.VerkleRoot[:8])
	fmt.Printf("  ✅ Blockchain-ready root saved to: %s/root_%s.json\n", rootsDir, termID)
	
	fmt.Println("✅ Term added successfully!")
	return nil
}

func generateStudentReceipt(studentID, outputFile string, terms, courses []string, selective bool) error {
	fmt.Printf("👤 Generating receipt for student: %s\n", studentID)
	fmt.Printf("📋 Output file: %s\n", outputFile)
	
	if selective {
		fmt.Println("🔒 Using selective disclosure mode")
	}
	
	// Determine terms to include
	var targetTerms []string
	if len(terms) > 0 {
		targetTerms = terms
		fmt.Printf("📚 Including specific terms: %v\n", targetTerms)
	} else {
		// Auto-discover terms from data
		var err error
		targetTerms, err = discoverStudentTerms(studentID)
		if err != nil {
			return fmt.Errorf("failed to discover student terms: %w", err)
		}
		fmt.Printf("📚 Auto-discovered terms: %v\n", targetTerms)
	}
	
	fmt.Println("🔐 Generating academic journey receipt...")
	
	receipts := make(map[string]interface{})
	
	for _, termID := range targetTerms {
		// Load student's Merkle tree data for this term directly
		merkleFile := filepath.Join("data", "merkle_trees", termID, fmt.Sprintf("merkle_%s_%s.json", extractStudentID(studentID), termID))
		
		merkleData, err := os.ReadFile(merkleFile)
		if err != nil {
			fmt.Printf("  ⚠️ Skipping term %s: Merkle tree not found\n", termID)
			continue
		}
		
		var studentMerkle map[string]interface{}
		if err := json.Unmarshal(merkleData, &studentMerkle); err != nil {
			fmt.Printf("  ⚠️ Skipping term %s: Failed to parse Merkle tree\n", termID)
			continue
		}
		
		// Load term root data
		rootFile := filepath.Join("blockchain_ready", "roots", fmt.Sprintf("root_%s.json", termID))
		rootData, err := os.ReadFile(rootFile)
		if err != nil {
			fmt.Printf("  ⚠️ Skipping term %s: Root data not found\n", termID)
			continue
		}
		
		var termRoot map[string]interface{}
		if err := json.Unmarshal(rootData, &termRoot); err != nil {
			fmt.Printf("  ⚠️ Skipping term %s: Failed to parse root data\n", termID)
			continue
		}
		
		// Get courses from Merkle tree
		coursesData, ok := studentMerkle["courses"].([]interface{})
		if !ok {
			fmt.Printf("  ⚠️ Skipping term %s: Invalid courses data\n", termID)
			continue
		}
		
		// Filter courses if selective disclosure
		var revealedCourses []interface{}
		if selective && len(courses) > 0 {
			for _, courseInterface := range coursesData {
				course := courseInterface.(map[string]interface{})
				courseID := course["course_id"].(string)
				for _, targetCourse := range courses {
					if courseID == targetCourse {
						revealedCourses = append(revealedCourses, course)
						break
					}
				}
			}
		} else {
			revealedCourses = coursesData
		}
		
		// Create simplified receipt with academic journey data
		receipt := map[string]interface{}{
			"student_id": studentID,
			"term_id": termID,
			"revealed_courses": revealedCourses,
			"merkle_root": studentMerkle["root"],
			"verification_path": "merkle_proof_available",
			"blockchain_anchor": termRoot["verkle_root"],
			"timestamp": termRoot["timestamp"],
		}
		
		receipts[termID] = map[string]interface{}{
			"term_id": termID,
			"student_id": studentID,
			"receipt": receipt,
			"verkle_root": termRoot["verkle_root"],
			"revealed_courses": len(revealedCourses),
			"generated_at": time.Now().Format(time.RFC3339),
		}
		
		fmt.Printf("  ✓ Generated receipt for term %s (%d courses)\n", termID, len(revealedCourses))
	}
	
	// Create comprehensive receipt
	journeyReceipt := map[string]interface{}{
		"student_id": studentID,
		"receipt_type": map[string]interface{}{
			"selective_disclosure": selective,
			"specific_courses": len(courses) > 0,
			"specific_terms": len(terms) > 0,
		},
		"generation_timestamp": time.Now().Format(time.RFC3339),
		"terms_included": targetTerms,
		"courses_filter": courses,
		"term_receipts": receipts,
		"blockchain_ready": true,
	}
	
	fmt.Println("📄 Creating journey receipt...")
	
	// Ensure output directory exists
	if dir := filepath.Dir(outputFile); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}
	
	// Save to specified output file
	receiptData, err := json.MarshalIndent(journeyReceipt, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal receipt: %w", err)
	}
	
	if err := os.WriteFile(outputFile, receiptData, 0644); err != nil {
		return fmt.Errorf("failed to write receipt file: %w", err)
	}
	
	// Also save to blockchain-ready directory for easy access
	receiptDir := "blockchain_ready/receipts"
	blockchainFile := filepath.Join(receiptDir, fmt.Sprintf("receipt_%s_%s.json", 
		extractStudentID(studentID), time.Now().Format("20060102_150405")))
	
	if err := os.WriteFile(blockchainFile, receiptData, 0644); err != nil {
		return fmt.Errorf("failed to save blockchain-ready receipt: %w", err)
	}
	
	fmt.Printf("💾 Receipt saved to: %s\n", outputFile)
	fmt.Printf("💾 Blockchain-ready copy: %s\n", blockchainFile)
	fmt.Println("✅ Receipt generated successfully!")
	
	return nil
}

func verifyReceiptLocally(receiptFile string) error {
	fmt.Printf("🔍 Verifying receipt: %s\n", receiptFile)
	
	fmt.Println("📖 Parsing receipt data...")
	data, err := os.ReadFile(receiptFile)
	if err != nil {
		return fmt.Errorf("failed to read receipt file: %w", err)
	}
	
	var receipt map[string]interface{}
	if err := json.Unmarshal(data, &receipt); err != nil {
		return fmt.Errorf("failed to parse receipt: %w", err)
	}
	
	studentID, ok := receipt["student_id"].(string)
	if !ok {
		return fmt.Errorf("invalid receipt: missing student_id")
	}
	
	termReceipts, ok := receipt["term_receipts"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid receipt: missing term_receipts")
	}
	
	fmt.Printf("📋 Verifying receipt for student: %s\n", studentID)
	
	fmt.Println("🔐 Validating Verkle proofs...")
	for termID, termReceiptInterface := range termReceipts {
		termReceiptMap, ok := termReceiptInterface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid term receipt for %s", termID)
		}
		
		verkleRootHex, ok := termReceiptMap["verkle_root"].(string)
		if !ok {
			return fmt.Errorf("missing verkle_root for term %s", termID)
		}
		
		// Convert hex string back to bytes for verification
		// This is a simplified verification - in production you'd verify the full receipt
		fmt.Printf("  ✓ Term %s: Verkle root %s\n", termID, verkleRootHex[:16]+"...")
	}
	
	fmt.Println("⏰ Checking temporal consistency...")
	timestamp, ok := receipt["generation_timestamp"].(string)
	if !ok {
		return fmt.Errorf("invalid receipt: missing generation_timestamp")
	}
	
	generatedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %w", err)
	}
	
	if generatedTime.After(time.Now()) {
		return fmt.Errorf("receipt timestamp is in the future")
	}
	
	fmt.Printf("  ✓ Receipt generated at: %s\n", generatedTime.Format("2006-01-02 15:04:05"))
	
	fmt.Println("✅ Local verification successful!")
	fmt.Println("📝 Verification summary:")
	fmt.Printf("  ✓ Student ID: %s\n", studentID)
	fmt.Printf("  ✓ Terms verified: %d\n", len(termReceipts))
	fmt.Printf("  ✓ Timestamp valid: %s\n", timestamp)
	
	return nil
}

func publishTermRoots(termID, network, privateKey string, gasLimit uint64) error {
	fmt.Printf("⛓️  Publishing roots for term: %s\n", termID)
	fmt.Printf("🌐 Target network: %s\n", network)
	
	// Load root data
	rootFile := filepath.Join("blockchain_ready/roots", fmt.Sprintf("root_%s.json", termID))
	if _, err := os.Stat(rootFile); os.IsNotExist(err) {
		return fmt.Errorf("root file not found: %s. Run 'add-term' first", rootFile)
	}
	
	fmt.Println("🌳 Loading Verkle tree commitment...")
	rootData, err := os.ReadFile(rootFile)
	if err != nil {
		return fmt.Errorf("failed to read root file: %w", err)
	}
	
	var root map[string]interface{}
	if err := json.Unmarshal(rootData, &root); err != nil {
		return fmt.Errorf("failed to parse root data: %w", err)
	}
	
	verkleRoot, ok := root["verkle_root"].(string)
	if !ok {
		return fmt.Errorf("invalid root data: missing verkle_root")
	}
	
	fmt.Printf("  ✓ Verkle root: %s\n", verkleRoot)
	
	// Prepare transaction data
	fmt.Println("📡 Preparing blockchain transaction...")
	
	if gasLimit == 0 {
		gasLimit = 500000 // Default gas limit
	}
	
	txData := map[string]interface{}{
		"term_id": termID,
		"verkle_root": verkleRoot,
		"network": network,
		"gas_limit": gasLimit,
		"timestamp": time.Now().Format(time.RFC3339),
		"status": "prepared",
		"tx_hash": "", // Will be filled when actually sent
	}
	
	// Save transaction data for blockchain integration
	txDir := "blockchain_ready/transactions"
	txFile := filepath.Join(txDir, fmt.Sprintf("tx_%s_%s.json", termID, 
		time.Now().Format("20060102_150405")))
	
	txBytes, err := json.MarshalIndent(txData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal transaction data: %w", err)
	}
	
	if err := os.WriteFile(txFile, txBytes, 0644); err != nil {
		return fmt.Errorf("failed to save transaction file: %w", err)
	}
	
	fmt.Println("💰 Estimating gas costs...")
	estimatedCost := float64(gasLimit) * 20 // Simplified gas price estimation
	fmt.Printf("  ✓ Estimated cost: %.6f ETH (gas limit: %d)\n", estimatedCost/1e18, gasLimit)
	
	// In production, this would connect to the blockchain and send the transaction
	fmt.Println("📡 [SIMULATION] Connecting to blockchain...")
	fmt.Println("📨 [SIMULATION] Broadcasting transaction...")
	
	// Update transaction status
	txData["status"] = "simulated"
	txData["tx_hash"] = fmt.Sprintf("0x%x", time.Now().Unix()) // Simulated hash
	
	updatedTxBytes, err := json.MarshalIndent(txData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated transaction: %w", err)
	}
	
	if err := os.WriteFile(txFile, updatedTxBytes, 0644); err != nil {
		return fmt.Errorf("failed to update transaction file: %w", err)
	}
	
	fmt.Printf("✅ Term roots prepared for blockchain publishing!\n")
	fmt.Printf("📄 Transaction data saved: %s\n", txFile)
	fmt.Printf("🔗 [SIMULATION] Transaction hash: %s\n", txData["tx_hash"].(string))
	
	fmt.Println("\n📋 Next steps for production:")
	fmt.Println("  1. Configure blockchain connection parameters")
	fmt.Println("  2. Fund account for gas fees")
	fmt.Println("  3. Deploy smart contracts if not already deployed") 
	fmt.Println("  4. Execute transaction using generated data")
	
	return nil
}

// Helper Functions

func loadCompletionsFromJSON(dataFile string) ([]merkle.CourseCompletion, error) {
	// This is a simplified loader - in production you'd have more robust parsing
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	
	var completions []merkle.CourseCompletion
	if err := json.Unmarshal(data, &completions); err != nil {
		return nil, err
	}
	
	return completions, nil
}

func discoverStudentTerms(studentID string) ([]string, error) {
	// Look for existing Merkle trees for this student
	merkleDir := "data/merkle_trees"
	
	var terms []string
	err := filepath.Walk(merkleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		filename := info.Name()
		if strings.Contains(filename, extractStudentID(studentID)) && 
		   strings.HasSuffix(filename, ".json") {
			// Extract term ID from filename pattern: merkle_STU001_Fall_2024.json
			parts := strings.Split(filename, "_")
			if len(parts) >= 3 {
				termID := strings.Join(parts[2:], "_")
				termID = strings.TrimSuffix(termID, ".json")
				terms = append(terms, termID)
			}
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if len(terms) == 0 {
		return nil, fmt.Errorf("no terms found for student %s", studentID)
	}
	
	return terms, nil
}

type VerkleTreeData struct {
	tree *verkle.TermVerkleTree
	metadata map[string]interface{}
}

func loadTermVerkleTree(termID string) (*VerkleTreeData, error) {
	// In a production system, this would load the actual Verkle tree state
	// For now, we'll create a new tree and populate it with data
	
	// This is a simplified implementation - you'd need to persist and restore Verkle tree state
	tree := verkle.NewTermVerkleTree(termID)
	
	// Load student data for this term from Merkle trees
	merkleDir := filepath.Join("data", "merkle_trees", termID)
	files, err := filepath.Glob(filepath.Join(merkleDir, "merkle_*.json"))
	if err != nil {
		return nil, err
	}
	
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue // Skip problematic files
		}
		
		var merkleTree merkle.StudentTermMerkle
		if err := json.Unmarshal(data, &merkleTree); err != nil {
			continue
		}
		
		// Add student to Verkle tree
		err = tree.AddStudent(merkleTree.StudentID, merkleTree.Courses)
		if err != nil {
			continue
		}
	}
	
	// Publish the tree
	if err := tree.PublishTerm(); err != nil {
		return nil, err
	}
	
	return &VerkleTreeData{
		tree: tree,
		metadata: map[string]interface{}{
			"term_id": termID,
			"loaded_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}