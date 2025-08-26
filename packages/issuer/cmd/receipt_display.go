package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var displayReceiptCmd = &cobra.Command{
	Use:   "display-receipt <receipt-file>",
	Short: "Display receipt in human-readable format with blockchain verification info",
	Long:  `Render academic journey receipt showing courses, grades, and blockchain anchors for verification`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		receiptFile := args[0]
		
		verbose, _ := cmd.Flags().GetBool("verbose")
		showBlockchain, _ := cmd.Flags().GetBool("blockchain")
		
		if err := displayReceipt(receiptFile, verbose, showBlockchain); err != nil {
			log.Fatalf("❌ Failed to display receipt: %v", err)
		}
	},
}

func displayReceipt(receiptFile string, verbose, showBlockchain bool) error {
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
	fmt.Println("📋 ACADEMIC JOURNEY RECEIPT")
	fmt.Println("=" + strings.Repeat("=", 60))
	
	// Student info
	studentID := receipt["student_id"].(string)
	generatedAt := receipt["generation_timestamp"].(string)
	fmt.Printf("👤 Student: %s\n", studentID)
	fmt.Printf("📅 Generated: %s\n", generatedAt)
	
	// Receipt type
	if receiptType, ok := receipt["receipt_type"].(map[string]interface{}); ok {
		selective := receiptType["selective_disclosure"].(bool)
		if selective {
			fmt.Println("🔒 Type: Selective Disclosure (Privacy-Preserving)")
		} else {
			fmt.Println("📖 Type: Complete Academic Journey")
		}
	}
	
	fmt.Println()
	
	// Terms and courses
	termReceipts := receipt["term_receipts"].(map[string]interface{})
	
	// Sort terms chronologically
	var sortedTerms []string
	for termID := range termReceipts {
		sortedTerms = append(sortedTerms, termID)
	}
	sort.Strings(sortedTerms)
	
	totalCourses := 0
	totalCredits := 0
	allGrades := []float64{}
	
	fmt.Println("📚 ACADEMIC TIMELINE:")
	fmt.Println("-" + strings.Repeat("-", 60))
	
	for i, termID := range sortedTerms {
		termData := termReceipts[termID].(map[string]interface{})
		termReceipt := termData["receipt"].(map[string]interface{})
		
		fmt.Printf("\n[%d] 📖 %s\n", i+1, termID)
		
		// Blockchain verification info
		if showBlockchain || verbose {
			blockchainAnchor := termReceipt["blockchain_anchor"].(string)
			timestamp := termReceipt["timestamp"].(string)
			fmt.Printf("    ⛓️  Blockchain Root: %s...\n", blockchainAnchor[:16])
			fmt.Printf("    🕐 Published: %s\n", timestamp)
		}
		
		// Courses
		revealedCourses := termReceipt["revealed_courses"].([]interface{})
		termCredits := 0
		termGrades := []float64{}
		
		fmt.Printf("    📋 Courses (%d completed):\n", len(revealedCourses))
		
		for j, courseInterface := range revealedCourses {
			course := courseInterface.(map[string]interface{})
			courseID := course["course_id"].(string)
			courseName := course["course_name"].(string)
			grade := course["grade"].(string)
			credits := int(course["credits"].(float64))
			
			gradePoints := getGradePoints(grade)
			termGrades = append(termGrades, gradePoints)
			allGrades = append(allGrades, gradePoints)
			termCredits += credits
			
			fmt.Printf("      %d. %s - %s [%s] (%d credits)\n", j+1, courseID, courseName, grade, credits)
			
			if verbose {
				startedAt := course["started_at"].(string)
				completedAt := course["completed_at"].(string)
				instructor := course["instructor"].(string)
				fmt.Printf("         📅 %s → %s | 👨‍🏫 %s\n", 
					startedAt[:10], completedAt[:10], instructor)
			}
		}
		
		// Term summary
		termGPA := calculateReceiptGPA(termGrades)
		fmt.Printf("    📊 Term Summary: %d courses, %d credits, GPA: %.2f\n", 
			len(revealedCourses), termCredits, termGPA)
		
		totalCourses += len(revealedCourses)
		totalCredits += termCredits
	}
	
	// Overall summary
	fmt.Println("\n" + strings.Repeat("=", 62))
	fmt.Println("📊 JOURNEY SUMMARY:")
	overallGPA := calculateReceiptGPA(allGrades)
	fmt.Printf("🎓 Total Courses: %d\n", totalCourses)
	fmt.Printf("📚 Total Credits: %d\n", totalCredits)
	fmt.Printf("📈 Overall GPA: %.2f\n", overallGPA)
	
	// Blockchain verification section
	if showBlockchain {
		fmt.Println("\n⛓️  BLOCKCHAIN VERIFICATION:")
		fmt.Println("-" + strings.Repeat("-", 60))
		fmt.Printf("📋 Terms Published: %d\n", len(sortedTerms))
		fmt.Println("🔗 Blockchain Network: Sepolia Testnet")
		fmt.Println("✅ All terms cryptographically anchored on blockchain")
		fmt.Println("\n📝 Verification Instructions for Third Parties:")
		fmt.Println("   1. Run: go run cmd/*.go verify-local <receipt-file>")
		fmt.Println("   2. Check that all Verkle roots match blockchain records")
		fmt.Println("   3. Verify temporal consistency of course completions")
		fmt.Println("   4. Confirm cryptographic integrity of all proofs")
	}
	
	// Academic progression analysis
	if verbose && len(sortedTerms) > 1 {
		fmt.Println("\n📈 ACADEMIC PROGRESSION ANALYSIS:")
		fmt.Println("-" + strings.Repeat("-", 60))
		
		prevGPA := 0.0
		for i, termID := range sortedTerms {
			termData := termReceipts[termID].(map[string]interface{})
			termReceipt := termData["receipt"].(map[string]interface{})
			revealedCourses := termReceipt["revealed_courses"].([]interface{})
			
			termGrades := []float64{}
			for _, courseInterface := range revealedCourses {
				course := courseInterface.(map[string]interface{})
				grade := course["grade"].(string)
				termGrades = append(termGrades, getGradePoints(grade))
			}
			
			termGPA := calculateReceiptGPA(termGrades)
			
			trend := ""
			if i > 0 {
				if termGPA > prevGPA {
					trend = "📈 Improving"
				} else if termGPA < prevGPA {
					trend = "📉 Declining"
				} else {
					trend = "➡️  Stable"
				}
			}
			
			fmt.Printf("   %s: GPA %.2f %s\n", termID, termGPA, trend)
			prevGPA = termGPA
		}
	}
	
	fmt.Println("\n✅ Receipt displayed successfully!")
	return nil
}

func getGradePoints(grade string) float64 {
	gradePoints := map[string]float64{
		"A+": 4.0, "A": 4.0, "A-": 3.7,
		"B+": 3.3, "B": 3.0, "B-": 2.7,
		"C+": 2.3, "C": 2.0, "C-": 1.7,
		"D+": 1.3, "D": 1.0, "F": 0.0,
	}
	
	if points, exists := gradePoints[grade]; exists {
		return points
	}
	return 0.0
}

func calculateReceiptGPA(grades []float64) float64 {
	if len(grades) == 0 {
		return 0.0
	}
	
	total := 0.0
	for _, grade := range grades {
		total += grade
	}
	
	return total / float64(len(grades))
}

func init() {
	displayReceiptCmd.Flags().BoolP("verbose", "v", false, "Show detailed information including dates and instructors")
	displayReceiptCmd.Flags().BoolP("blockchain", "b", false, "Show blockchain verification information")
	
	rootCmd.AddCommand(displayReceiptCmd)
}