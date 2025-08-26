package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var batchReceiptCmd = &cobra.Command{
	Use:   "generate-all-receipts",
	Short: "Generate receipts for all students in batch",
	Long:  `Generate academic journey receipts for all students found in the system`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🎓 Batch Receipt Generation")
		fmt.Println("=" + string(make([]byte, 50)))
		
		outputDir, _ := cmd.Flags().GetString("output-dir")
		selective, _ := cmd.Flags().GetBool("selective")
		specificTerms, _ := cmd.Flags().GetStringSlice("terms")
		specificCourses, _ := cmd.Flags().GetStringSlice("courses")
		
		if err := runBatchReceiptGeneration(outputDir, selective, specificTerms, specificCourses); err != nil {
			log.Fatalf("❌ Batch receipt generation failed: %v", err)
		}
		
		fmt.Println("\n🎉 Batch receipt generation completed successfully!")
	},
}

func runBatchReceiptGeneration(outputDir string, selective bool, terms, courses []string) error {
	// Step 1: Discover all students
	fmt.Println("\n🔍 Step 1: Discovering available students...")
	students, err := discoverAllStudents()
	if err != nil {
		return fmt.Errorf("failed to discover students: %w", err)
	}
	
	if len(students) == 0 {
		return fmt.Errorf("no students found - please run 'generate-data' first")
	}
	
	fmt.Printf("📚 Found %d students to process\n", len(students))
	
	// Step 2: Create output directory
	if outputDir == "" {
		outputDir = "receipts"
	}
	
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	fmt.Printf("📁 Output directory: %s\n", outputDir)
	
	// Step 3: Generate receipts for each student
	fmt.Printf("\n🎓 Step 2: Generating receipts for %d students...\n", len(students))
	
	successCount := 0
	failureCount := 0
	
	for i, studentID := range students {
		fmt.Printf("  [%d/%d] 📋 Generating receipt for %s...\n", i+1, len(students), studentID)
		
		// Create filename
		suffix := "complete"
		if selective && len(terms) > 0 {
			suffix = "terms"
		}
		if selective && len(courses) > 0 {
			suffix = "courses"  
		}
		
		filename := fmt.Sprintf("%s_%s_journey.json", studentID, suffix)
		outputFile := filepath.Join(outputDir, filename)
		
		// Generate receipt
		err := generateStudentReceipt(studentID, outputFile, terms, courses, selective)
		if err != nil {
			fmt.Printf("    ❌ Failed to generate receipt for %s: %v\n", studentID, err)
			failureCount++
			continue
		}
		
		fmt.Printf("    ✅ Receipt generated: %s\n", filename)
		successCount++
	}
	
	// Step 4: Show summary
	fmt.Printf("\n📊 Batch Generation Summary:\n")
	fmt.Printf("  ✅ Successfully generated: %d receipts\n", successCount)
	if failureCount > 0 {
		fmt.Printf("  ❌ Failed to generate: %d receipts\n", failureCount)
	}
	fmt.Printf("  📁 Receipts saved in: %s/\n", outputDir)
	
	// Step 5: List generated files
	fmt.Printf("\n📄 Generated Receipt Files:\n")
	files, err := filepath.Glob(filepath.Join(outputDir, "*.json"))
	if err == nil {
		for _, file := range files {
			basename := filepath.Base(file)
			fmt.Printf("  • %s\n", basename)
		}
	}
	
	return nil
}

func discoverAllStudents() ([]string, error) {
	// Look for students in the system summary
	summaryPath := "data/generated_student_data/system_summary.json"
	
	data, err := os.ReadFile(summaryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read system summary: %w", err)
	}
	
	var summary struct {
		Students []string `json:"students"`
	}
	
	if err := json.Unmarshal(data, &summary); err != nil {
		return nil, fmt.Errorf("failed to parse system summary: %w", err)
	}
	
	return summary.Students, nil
}

func init() {
	batchReceiptCmd.Flags().StringP("output-dir", "o", "receipts", "Output directory for generated receipts")
	batchReceiptCmd.Flags().BoolP("selective", "s", false, "Use selective disclosure mode")
	batchReceiptCmd.Flags().StringSliceP("terms", "t", []string{}, "Specific terms to include (comma-separated)")
	batchReceiptCmd.Flags().StringSliceP("courses", "c", []string{}, "Specific courses to include (comma-separated)")
	
	rootCmd.AddCommand(batchReceiptCmd)
}