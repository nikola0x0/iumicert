package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"iumicert/crypto/merkle"
)

var convertDataCmd = &cobra.Command{
	Use:   "convert-data [term-id]",
	Short: "Convert generated student data to Merkle tree format",
	Long:  `Convert the generated student journey data into the format expected by the add-term command`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		termID := args[0]
		
		fmt.Printf("🔄 Converting student data for term: %s\n", termID)
		
		if err := convertStudentDataToMerkleFormat(termID); err != nil {
			fmt.Printf("❌ Conversion failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✅ Conversion completed successfully!")
	},
}

type StudentJourney struct {
	StudentID string                 `json:"student_id"`
	Terms     map[string]TermData    `json:"terms"`
}

type TermData struct {
	TermID       string        `json:"term_id"`
	Courses      []CourseData  `json:"courses"`
	TotalCredits int           `json:"total_credits"`
	GPA          float64       `json:"gpa"`
}

type CourseData struct {
	CourseID    string    `json:"course_id"`
	CourseName  string    `json:"course_name"`
	Credits     uint8     `json:"credits"`
	Grade       string    `json:"grade"`
	Instructor  string    `json:"instructor"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
	AssessedAt  time.Time `json:"assessed_at"`
	IssuedAt    time.Time `json:"issued_at"`
	AttemptNo   int       `json:"attempt_no"`
	IssuerID    string    `json:"issuer_id"`
	TermID      string    `json:"term_id"`
}

func convertStudentDataToMerkleFormat(termID string) error {
	// Read all student journey files
	studentsDir := "data/generated_student_data/students"
	files, err := filepath.Glob(filepath.Join(studentsDir, "journey_*.json"))
	if err != nil {
		return fmt.Errorf("failed to find student files: %w", err)
	}

	var allCompletions []merkle.CourseCompletion

	fmt.Printf("📚 Processing %d student files...\n", len(files))

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue // Skip problematic files
		}

		var journey StudentJourney
		if err := json.Unmarshal(data, &journey); err != nil {
			continue // Skip problematic files
		}

		// Check if this student has data for the requested term
		if termData, exists := journey.Terms[termID]; exists {
			// Convert each course to CourseCompletion format
			for _, course := range termData.Courses {
				completion := merkle.CourseCompletion{
					IssuerID:    course.IssuerID,
					StudentID:   extractStudentIDFromDID(journey.StudentID),
					TermID:      termID,
					CourseID:    course.CourseID,
					CourseName:  course.CourseName,
					AttemptNo:   uint8(course.AttemptNo),
					StartedAt:   course.StartedAt,
					CompletedAt: course.CompletedAt,
					AssessedAt:  course.AssessedAt,
					IssuedAt:    course.IssuedAt,
					Grade:       course.Grade,
					Credits:     course.Credits,
					Instructor:  course.Instructor,
				}
				
				allCompletions = append(allCompletions, completion)
			}
		}
	}

	if len(allCompletions) == 0 {
		return fmt.Errorf("no completions found for term %s", termID)
	}

	// Create output directory
	outputDir := "data/converted_terms"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save converted data
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_completions.json", termID))
	outputData, err := json.MarshalIndent(allCompletions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal completions: %w", err)
	}

	if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("📊 Converted %d course completions for term %s\n", len(allCompletions), termID)
	fmt.Printf("💾 Saved to: %s\n", outputFile)

	return nil
}

func extractStudentIDFromDID(did string) string {
	// Extract student ID from DID format "did:example:ITITIU00001"
	if len(did) >= 12 {
		return did[12:] // Skip "did:example:"
	}
	return did
}

func init() {
	rootCmd.AddCommand(convertDataCmd)
}