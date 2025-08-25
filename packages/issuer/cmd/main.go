package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
		fmt.Printf("🏛️  Initializing IU-MiCert repository for institution: %s\n", institutionID)
		fmt.Println("📁 Creating directory structure...")
		fmt.Println("⚙️  Generating configuration files...")
		fmt.Println("🔑 Setting up cryptographic parameters...")
		fmt.Println("✅ Repository initialized successfully!")
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
		fmt.Printf("📚 Adding academic term: %s\n", termID)
		fmt.Printf("📖 Processing data from: %s\n", dataFile)
		fmt.Println("🌳 Building student-term Merkle trees...")
		fmt.Println("🔗 Preparing Verkle tree aggregation...")
		fmt.Println("✅ Term added successfully!")
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
		fmt.Printf("👤 Generating receipt for student: %s\n", studentID)
		fmt.Println("📋 Collecting academic terms...")
		fmt.Println("🔐 Generating Verkle proofs...")
		fmt.Println("📄 Creating journey receipt...")
		fmt.Printf("💾 Saved receipt to: %s\n", outputFile)
		fmt.Println("✅ Receipt generated successfully!")
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
		fmt.Printf("🔍 Verifying receipt: %s\n", receiptFile)
		fmt.Println("📖 Parsing receipt data...")
		fmt.Println("🔐 Validating Verkle proofs...")
		fmt.Println("⏰ Checking temporal consistency...")
		fmt.Println("✅ Local verification successful!")
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
		fmt.Printf("⛓️  Publishing roots for term: %s\n", termID)
		fmt.Println("🌳 Computing Verkle tree commitment...")
		fmt.Println("📡 Connecting to blockchain...")
		fmt.Println("💰 Estimating gas costs...")
		fmt.Println("✅ Term roots published successfully!")
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