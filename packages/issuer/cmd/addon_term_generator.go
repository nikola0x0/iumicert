package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"iumicert/crypto/testdata"
	"iumicert/crypto/verkle"
)

var addonTermGeneratorCmd = &cobra.Command{
	Use:   "generate-addon-term",
	Short: "Generate a single additional term for demo/testing",
	Long: `Generate a single new academic term with course completions for existing students.
This is useful for creating additional terms for demo purposes without regenerating all data.`,
	Run: func(cmd *cobra.Command, args []string) {
		termID, _ := cmd.Flags().GetString("term")
		numStudents, _ := cmd.Flags().GetInt("students")
		outputPath, _ := cmd.Flags().GetString("output")

		fmt.Printf("📊 Generating Add-on Term: %s\n", termID)
		fmt.Println("=" + string(make([]byte, 50)))

		if err := runAddonTermGeneration(termID, numStudents, outputPath); err != nil {
			log.Fatalf("❌ Add-on term generation failed: %v", err)
		}

		fmt.Printf("\n🎉 Successfully generated term: %s\n", termID)
		fmt.Printf("📁 Output saved to: %s\n", outputPath)
	},
}

func runAddonTermGeneration(termID string, numStudents int, outputPath string) error {
	// Validate number of students
	if numStudents < 1 {
		return fmt.Errorf("number of students must be at least 1, got %d", numStudents)
	}

	// Fixed course range: 3-6 courses per student
	const minCourses = 3
	const maxCourses = 6

	fmt.Println("\n🔧 Step 1: Initializing test data generator...")
	generator := testdata.NewTestDataGenerator()

	fmt.Printf("\n📚 Step 2: Generating data for term: %s\n", termID)
	fmt.Printf("  • Students: %d\n", numStudents)
	fmt.Printf("  • Courses per student: %d-%d (random)\n", minCourses, maxCourses)

	// Generate course completions
	allCompletions := make([]verkle.CourseCompletion, 0)

	for i := 0; i < numStudents; i++ {
		// Variable course count per student (3-6)
		coursesForStudent := minCourses + (i % (maxCourses - minCourses + 1))

		studentCompletions, err := generator.GenerateTermData(termID, 1, coursesForStudent)
		if err != nil {
			return fmt.Errorf("failed to generate data for student %d: %w", i, err)
		}

		// Update student IDs to match IU Vietnam format
		for j := range studentCompletions {
			studentCompletions[j].StudentID = fmt.Sprintf("ITITIU%05d", i+1)
		}

		allCompletions = append(allCompletions, studentCompletions...)
	}

	fmt.Printf("  ✅ Generated %d course completions for %d students\n", len(allCompletions), numStudents)

	// Step 3: Organize data by student
	fmt.Println("\n💾 Step 3: Organizing term data...")

	studentCourses := make(map[string][]map[string]interface{})

	for _, completion := range allCompletions {
		studentID := completion.StudentID

		courseData := map[string]interface{}{
			"course_id":   completion.CourseID,
			"course_name": completion.CourseName,
			"credits":     completion.Credits,
			"grade":       completion.Grade,
		}

		studentCourses[studentID] = append(studentCourses[studentID], courseData)
	}

	// Create the term data structure expected by the backend
	termData := map[string]interface{}{
		"term_id":  termID,
		"students": studentCourses,
		"metadata": map[string]interface{}{
			"generated_at":    time.Now().Format(time.RFC3339),
			"total_students":  len(studentCourses),
			"total_courses":   len(allCompletions),
			"generated_by":    "generate-addon-term",
		},
	}

	// Step 4: Save to file
	fmt.Println("\n💾 Step 4: Saving term data...")

	// Create output directory if needed
	outputDir := filepath.Dir(outputPath)
	if outputDir != "." && outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Save the file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(termData); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	// Step 5: Display statistics
	fmt.Println("\n📊 Generation Statistics:")
	fmt.Printf("  📖 Term ID: %s\n", termID)
	fmt.Printf("  👥 Students: %d\n", len(studentCourses))
	fmt.Printf("  📚 Total Courses: %d\n", len(allCompletions))
	fmt.Printf("  📁 Output File: %s\n", outputPath)

	// Display per-student breakdown
	fmt.Println("\n  Per-Student Course Count:")
	for studentID, courses := range studentCourses {
		fmt.Printf("    %s: %d courses\n", studentID, len(courses))
	}

	return nil
}

func init() {
	addonTermGeneratorCmd.Flags().String("term", "Semester_1_2025", "Term ID (e.g., Semester_1_2025, Summer_2025)")
	addonTermGeneratorCmd.Flags().Int("students", 5, "Number of students")
	addonTermGeneratorCmd.Flags().String("output", "addon_term.json", "Output file path")

	rootCmd.AddCommand(addonTermGeneratorCmd)
}
