package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var batchProcessCmd = &cobra.Command{
	Use:   "batch-process",
	Short: "Process all available terms from generated student data",
	Long:  `Automatically process all available academic terms into the credential system`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🏭 Starting Batch Processing Pipeline")
		fmt.Println("=" + string(make([]byte, 50)))
		
		if err := runBatchProcessing(); err != nil {
			fmt.Printf("❌ Batch processing failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✅ Batch processing completed successfully!")
	},
}

func runBatchProcessing() error {
	// Step 1: Discovery available terms from generated data
	fmt.Println("\n🔍 Step 1: Discovering available terms...")
	
	termSummariesDir := "data/generated_student_data/terms"
	summaryFiles, err := filepath.Glob(filepath.Join(termSummariesDir, "summary_*.json"))
	if err != nil {
		return fmt.Errorf("failed to find term summaries: %w", err)
	}

	if len(summaryFiles) == 0 {
		return fmt.Errorf("no term summaries found. Run 'generate-data' first")
	}

	var availableTerms []string
	for _, file := range summaryFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var summary map[string]interface{}
		if err := json.Unmarshal(data, &summary); err != nil {
			continue
		}

		if termID, ok := summary["term_id"].(string); ok {
			availableTerms = append(availableTerms, termID)
		}
	}

	fmt.Printf("📚 Found %d terms to process: %v\n", len(availableTerms), availableTerms)

	// Step 2: Convert and process each term
	for i, termID := range availableTerms {
		fmt.Printf("\n📖 Step %d: Processing term %s (%d/%d)\n", i+2, termID, i+1, len(availableTerms))
		
		// Convert data
		fmt.Printf("  🔄 Converting student data to Merkle format...\n")
		if err := convertStudentDataToMerkleFormat(termID); err != nil {
			fmt.Printf("  ⚠️  Warning: Failed to convert %s: %v\n", termID, err)
			continue
		}

		// Process term
		fmt.Printf("  🌳 Building Merkle/Verkle trees...\n")
		dataFile := filepath.Join("data/converted_terms", fmt.Sprintf("%s_completions.json", termID))
		
		if err := addAcademicTerm(termID, dataFile, "json", true); err != nil {
			fmt.Printf("  ⚠️  Warning: Failed to process %s: %v\n", termID, err)
			continue
		}

		fmt.Printf("  ✅ Term %s processed successfully\n", termID)
		time.Sleep(100 * time.Millisecond) // Brief pause between terms
	}

	// Step 3: Generate summary report
	fmt.Printf("\n📊 Step %d: Generating processing report...\n", len(availableTerms)+2)
	
	report := map[string]interface{}{
		"processing_timestamp": time.Now().Format(time.RFC3339),
		"terms_processed":     len(availableTerms),
		"terms":              availableTerms,
		"status":             "completed",
		"ready_for_blockchain": true,
	}

	reportPath := "data/batch_processing_report.json"
	reportData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create report: %w", err)
	}

	if err := os.WriteFile(reportPath, reportData, 0644); err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	fmt.Printf("📄 Processing report saved: %s\n", reportPath)

	return nil
}

func init() {
	rootCmd.AddCommand(batchProcessCmd)
}