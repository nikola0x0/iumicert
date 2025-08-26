package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var verificationGuideCmd = &cobra.Command{
	Use:   "verification-guide <receipt-file>",
	Short: "Show step-by-step verification guide for third parties",
	Long:  `Display comprehensive verification instructions showing how to validate blockchain anchors and cryptographic proofs`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		receiptFile := args[0]
		
		if err := showVerificationGuide(receiptFile); err != nil {
			log.Fatalf("❌ Failed to show verification guide: %v", err)
		}
	},
}

func showVerificationGuide(receiptFile string) error {
	// Read receipt file
	data, err := os.ReadFile(receiptFile)
	if err != nil {
		return fmt.Errorf("failed to read receipt file: %w", err)
	}
	
	var receipt map[string]interface{}
	if err := json.Unmarshal(data, &receipt); err != nil {
		return fmt.Errorf("failed to parse receipt: %w", err)
	}
	
	// Display header
	fmt.Println("🔍 BLOCKCHAIN VERIFICATION GUIDE")
	fmt.Println("===============================================")
	
	studentID := receipt["student_id"].(string)
	fmt.Printf("👤 Verifying: %s\n", studentID)
	fmt.Printf("📄 Receipt File: %s\n", receiptFile)
	fmt.Println()
	
	// Step 1: Local verification
	fmt.Println("📋 STEP 1: RUN LOCAL VERIFICATION")
	fmt.Println("-----------------------------------------------")
	fmt.Printf("go run cmd/*.go verify-local %s\n", receiptFile)
	fmt.Println()
	fmt.Println("✅ This will verify:")
	fmt.Println("   • Cryptographic integrity of all proofs")
	fmt.Println("   • Temporal consistency of course completions")
	fmt.Println("   • Merkle tree structure and student data")
	fmt.Println("   • Term root hashes match blockchain anchors")
	fmt.Println()
	
	// Step 2: Blockchain verification details
	fmt.Println("⛓️  STEP 2: VERIFY BLOCKCHAIN ANCHORS")
	fmt.Println("-----------------------------------------------")
	
	termReceipts := receipt["term_receipts"].(map[string]interface{})
	termCount := 0
	
	for termID, termData := range termReceipts {
		termCount++
		termReceipt := termData.(map[string]interface{})["receipt"].(map[string]interface{})
		blockchainAnchor := termReceipt["blockchain_anchor"].(string)
		timestamp := termReceipt["timestamp"].(string)
		
		fmt.Printf("[%d] 📚 Term: %s\n", termCount, termID)
		fmt.Printf("    🔗 Verkle Root: %s\n", blockchainAnchor)
		fmt.Printf("    📅 Published: %s\n", timestamp)
		fmt.Printf("    🌐 Network: Sepolia Testnet\n")
		fmt.Println()
		
		fmt.Println("    🔍 Manual Blockchain Verification:")
		fmt.Printf("    • Visit: https://sepolia.etherscan.io/\n")
		fmt.Printf("    • Search for transaction containing root: %s...\n", blockchainAnchor[:16])
		fmt.Printf("    • Verify publisher is authorized university wallet\n")
		fmt.Printf("    • Confirm timestamp matches: %s\n", timestamp[:10])
		fmt.Println()
	}
	
	// Step 3: Academic integrity verification
	fmt.Println("🎓 STEP 3: ACADEMIC INTEGRITY CHECKS")
	fmt.Println("-----------------------------------------------")
	fmt.Println("✅ Verify the following academic standards:")
	fmt.Println("   • Course completion dates are chronologically consistent")
	fmt.Println("   • Prerequisites are satisfied (earlier courses before advanced)")
	fmt.Println("   • Credit hours match institutional standards")
	fmt.Println("   • GPA calculations are mathematically correct")
	fmt.Println("   • Issuing institution is accredited and authorized")
	fmt.Println()
	
	// Step 4: Cryptographic verification details
	fmt.Println("🔐 STEP 4: CRYPTOGRAPHIC PROOF VERIFICATION")
	fmt.Println("-----------------------------------------------")
	fmt.Println("The receipt contains these cryptographic components:")
	fmt.Println()
	
	for termID, termData := range termReceipts {
		termReceipt := termData.(map[string]interface{})["receipt"].(map[string]interface{})
		revealedCourses := termReceipt["revealed_courses"].([]interface{})
		
		fmt.Printf("📚 %s:\n", termID)
		fmt.Printf("   🌳 Merkle Tree: Proves %d courses belong to student\n", len(revealedCourses))
		fmt.Printf("   🔗 Verkle Proof: Links student data to blockchain root\n")
		fmt.Printf("   🕐 Timestamps: Proves when courses were completed\n")
		fmt.Printf("   📋 Course Hashes: Ensures course data hasn't been tampered\n")
		fmt.Println()
	}
	
	// Step 5: Privacy verification
	if receiptType, ok := receipt["receipt_type"].(map[string]interface{}); ok {
		selective := receiptType["selective_disclosure"].(bool)
		
		fmt.Println("🔒 STEP 5: PRIVACY VERIFICATION")
		fmt.Println("-----------------------------------------------")
		if selective {
			fmt.Println("🔒 This is a SELECTIVE DISCLOSURE receipt")
			fmt.Println("   • Student chose to reveal only specific terms/courses")
			fmt.Println("   • Hidden data is cryptographically protected")
			fmt.Println("   • Blockchain still proves complete academic integrity")
			fmt.Println("   • Verifier sees only authorized information")
		} else {
			fmt.Println("📖 This is a COMPLETE JOURNEY receipt")
			fmt.Println("   • Student chose to reveal full academic history")
			fmt.Println("   • All terms and courses are visible")
			fmt.Println("   • Complete timeline of academic progression")
		}
		fmt.Println()
	}
	
	// Summary
	fmt.Println("📊 VERIFICATION SUMMARY")
	fmt.Println("===============================================")
	fmt.Printf("🎓 Student: %s\n", studentID)
	fmt.Printf("📚 Terms to verify: %d\n", len(termReceipts))
	fmt.Println("⛓️  Blockchain: Sepolia Testnet")
	fmt.Println("🔐 Cryptography: Merkle + Verkle Trees")
	fmt.Println("✅ Status: Ready for third-party verification")
	fmt.Println()
	
	fmt.Println("🚨 IMPORTANT FOR VERIFIERS:")
	fmt.Println("• This receipt is mathematically tamper-proof")
	fmt.Println("• Any modification will break cryptographic verification")
	fmt.Println("• Blockchain anchors provide immutable audit trail")
	fmt.Println("• Institution digital signatures ensure authenticity")
	
	return nil
}

func init() {
	rootCmd.AddCommand(verificationGuideCmd)
}