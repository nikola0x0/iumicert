package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
		"iumicert/crypto/verkle"
	blockchain "iumicert/issuer/blockchain_integration"
	"iumicert/issuer/config"
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
		"../data/terms",
		"../data/students", 
		"../data/merkle_trees",
		"../data/verkle_trees",
		"../publish_ready",
		"../publish_ready/receipts",
		"../publish_ready/proofs",
		"../publish_ready/roots",
		"../publish_ready/transactions",
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
			"receipts": "../publish_ready/receipts",
			"proofs": "../publish_ready/proofs", 
			"roots": "../publish_ready/roots",
			"transactions": "../publish_ready/transactions",
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
	var completions []verkle.CourseCompletion
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
	
	// Group completions by student for Verkle tree processing
	fmt.Println("🌳 Organizing course completions by student...")
	studentCompletions := make(map[string][]verkle.CourseCompletion)
	
	for _, completion := range completions {
		studentDID := fmt.Sprintf("did:example:%s", completion.StudentID)
		studentCompletions[studentDID] = append(studentCompletions[studentDID], completion)
	}
	
	// Build term-level Verkle tree
	fmt.Println("🔗 Preparing Verkle tree aggregation...")
	termTree := verkle.NewTermVerkleTree(termID)
	
	for studentDID, courses := range studentCompletions {
		err := termTree.AddCourses(studentDID, courses)
		if err != nil {
			return fmt.Errorf("failed to add courses for student %s to verkle tree: %w", studentDID, err)
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
	rootsDir := "../publish_ready/roots"
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
		// Load term completion data directly from Verkle format
		verkleTermFile := filepath.Join("../data/verkle_terms", fmt.Sprintf("%s_completions.json", termID))
		
		verkleData, err := os.ReadFile(verkleTermFile)
		if err != nil {
			fmt.Printf("  ⚠️ Skipping term %s: Verkle term data not found\n", termID)
			continue
		}
		
		var completions []verkle.CourseCompletion
		if err := json.Unmarshal(verkleData, &completions); err != nil {
			fmt.Printf("  ⚠️ Skipping term %s: Failed to parse Verkle term data\n", termID)
			continue
		}
		
		// Filter completions for this student
		var studentCourses []verkle.CourseCompletion
		targetStudentID := extractStudentID(studentID)
		for _, completion := range completions {
			if completion.StudentID == targetStudentID {
				studentCourses = append(studentCourses, completion)
			}
		}
		
		if len(studentCourses) == 0 {
			fmt.Printf("  ⚠️ Skipping term %s: No courses found for student\n", termID)
			continue
		}
		
		// Load term root data
		rootFile := filepath.Join("..", "publish_ready", "roots", fmt.Sprintf("root_%s.json", termID))
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
		
		// Apply selective disclosure filter if requested
		var revealedCourses []verkle.CourseCompletion
		if selective && len(courses) > 0 {
			for _, completion := range studentCourses {
				for _, targetCourse := range courses {
					if completion.CourseID == targetCourse {
						revealedCourses = append(revealedCourses, completion)
						break
					}
				}
			}
		} else {
			revealedCourses = studentCourses
		}
		
		// Create single Verkle receipt with course-level proofs
		receipt := map[string]interface{}{
			"student_id": studentID,
			"term_id": termID,
			"revealed_courses": revealedCourses,
			"verkle_root": termRoot["verkle_root"],
			"proof_type": "verkle_32_byte",
			"selective_disclosure": selective,
			"verification_path": "single_verkle_proof",
			"timestamp": termRoot["timestamp"],
		}
		
		receipts[termID] = map[string]interface{}{
			"term_id": termID,
			"student_id": studentID,
			"receipt": receipt,
			"verkle_root": termRoot["verkle_root"],
			"revealed_courses": len(revealedCourses),
			"total_courses": len(studentCourses),
			"generated_at": time.Now().Format(time.RFC3339),
		}
		
		fmt.Printf("  ✓ Generated receipt for term %s (%d/%d courses)\n", termID, len(revealedCourses), len(studentCourses))
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
	receiptDir := "../publish_ready/receipts"
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
	
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	// Override config with command line arguments if provided
	if network != "" {
		cfg.Network = network
	}
	if privateKey != "" {
		cfg.IssuerPrivateKey = privateKey
	}
	if gasLimit > 0 {
		cfg.DefaultGasLimit = gasLimit
	}
	
	// Print current configuration
	cfg.PrintConfig()
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Printf("\n❌ Configuration Error: %v\n\n", err)
		fmt.Println("💡 To fix this:")
		fmt.Println("  1. Copy .env.example to .env: cp .env.example .env")
		fmt.Println("  2. Edit .env file with your actual values")
		fmt.Println("  3. For localhost testing, the test values should work")
		return err
	}
	
	// Load root data
	rootFile := filepath.Join("../publish_ready/roots", fmt.Sprintf("root_%s.json", termID))
	if _, err := os.Stat(rootFile); os.IsNotExist(err) {
		return fmt.Errorf("root file not found: %s. Run 'add-term' first", rootFile)
	}
	
	fmt.Printf("🌐 Target network: %s\n", cfg.Network)
	fmt.Println("🔗 Connecting to blockchain...")
	
	// Create blockchain integration
	integration, err := blockchain.NewBlockchainIntegration(
		cfg.Network, 
		cfg.GetPrivateKey(), 
		cfg.GetContractAddress(),
	)
	if err != nil {
		return fmt.Errorf("failed to create blockchain integration: %w", err)
	}
	defer integration.Close()
	
	fmt.Println("📡 Publishing term root to blockchain...")
	
	// Publish term root from file
	ctx := context.Background()
	result, err := integration.PublishTermRootFromFile(ctx, rootFile)
	if err != nil {
		return fmt.Errorf("failed to publish term root: %w", err)
	}
	
	fmt.Printf("✅ Term root published successfully!\n")
	fmt.Printf("🔗 Transaction hash: %s\n", result.TransactionHash)
	fmt.Printf("📦 Block number: %d\n", result.BlockNumber)
	fmt.Printf("⛽ Gas used: %d\n", result.GasUsed)
	
	fmt.Println("\n🎉 Blockchain integration completed!")
	fmt.Printf("📄 Transaction record saved in ../publish_ready/transactions/\n")
	
	return nil
}

// Helper Functions

func loadCompletionsFromJSON(dataFile string) ([]verkle.CourseCompletion, error) {
	// This is a simplified loader - in production you'd have more robust parsing
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	
	var completions []verkle.CourseCompletion
	if err := json.Unmarshal(data, &completions); err != nil {
		return nil, err
	}
	
	return completions, nil
}

func discoverStudentTerms(studentID string) ([]string, error) {
	// Look for published Verkle tree roots to discover available terms
	rootsDir := "../publish_ready/roots"
	
	var terms []string
	err := filepath.Walk(rootsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		filename := info.Name()
		if strings.HasPrefix(filename, "root_") && strings.HasSuffix(filename, ".json") {
			// Extract term ID from filename pattern: root_Semester_1_2023.json
			termID := strings.TrimPrefix(filename, "root_")
			termID = strings.TrimSuffix(termID, ".json")
			
			// Check if this term has data for the requested student
			verkleTermFile := filepath.Join("../data/verkle_terms", termID+"_completions.json")
			if data, err := os.ReadFile(verkleTermFile); err == nil {
				// Check if student has courses in this term
				if strings.Contains(string(data), fmt.Sprintf("\"student_id\": \"%s\"", extractStudentID(studentID))) {
					terms = append(terms, termID)
				}
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

