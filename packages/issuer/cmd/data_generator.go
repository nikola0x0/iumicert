package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"iumicert/crypto/verkle"
	"iumicert/crypto/testdata"
)

var dataGeneratorCmd = &cobra.Command{
	Use:   "generate-data",
	Short: "Generate comprehensive academic dataset for testing",
	Long:  `Generate realistic academic records with multi-term student journeys for system testing and validation`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📊 Starting Student Data Generation")
		fmt.Println("=" + string(make([]byte, 50)))
		
		if err := runStudentDataGeneration(); err != nil {
			log.Fatalf("❌ Student data generation failed: %v", err)
		}
		
		fmt.Println("\n🎉 Student data generation completed successfully!")
	},
}

func runStudentDataGeneration() error {
	// Step 1: Initialize generator
	fmt.Println("\n🔧 Step 1: Initializing test data generator...")
	generator := testdata.NewTestDataGenerator()
	
	// Configuration for IU Vietnam realistic data generation - 6 specific terms only
	terms := []string{"Semester_1_2023", "Semester_1_2024", "Semester_2_2023", "Semester_2_2024", "Summer_2023", "Summer_2024"}
	numStudents := 5  // Generate 5 students for demo
	minCoursesPerTerm := 3
	maxCoursesPerTerm := 6
	
	// Storage for organized data
	allStudentData := make(map[string]map[string]interface{})
	termSummaries := make(map[string]interface{})
	
	// Step 2: Generate data for each term
	fmt.Printf("\n📚 Step 2: Generating data for %d terms...\n", len(terms))
	
	for _, termID := range terms {
		fmt.Printf("\n  📖 Generating data for term: %s\n", termID)
		
		// Generate variable course completions per student (3-6 courses each)
		allCompletions := make([]verkle.CourseCompletion, 0)
		
		for i := 0; i < numStudents; i++ {
			// Random course count between min and max
			coursesForStudent := minCoursesPerTerm + (i % (maxCoursesPerTerm - minCoursesPerTerm + 1))
			
			// Generate completions for this student
			studentCompletions, err := generator.GenerateTermData(termID, 1, coursesForStudent)
			if err != nil {
				return fmt.Errorf("failed to generate data for student %d in term %s: %w", i, termID, err)
			}
			
			// Update student IDs to be unique across the batch (IU Vietnam format)
			for j := range studentCompletions {
				studentCompletions[j].StudentID = fmt.Sprintf("ITITIU%05d", i+1)
			}
			
			allCompletions = append(allCompletions, studentCompletions...)
		}
		
		fmt.Printf("    ✅ Generated %d course completions for %d students\n", len(allCompletions), numStudents)
		completions := allCompletions
		
		// Organize by student
		studentCompletions := make(map[string][]interface{})
		for _, completion := range completions {
			studentDID := fmt.Sprintf("did:example:%s", completion.StudentID)
			
			// Initialize student data if not exists
			if allStudentData[studentDID] == nil {
				allStudentData[studentDID] = make(map[string]interface{})
				allStudentData[studentDID]["student_id"] = studentDID
				allStudentData[studentDID]["terms"] = make(map[string]interface{})
			}
			
			// Convert completion to map for JSON serialization
			completionData := map[string]interface{}{
				"course_id":        completion.CourseID,
				"course_name":      completion.CourseName,
				"credits":          completion.Credits,
				"grade":            completion.Grade,
				"instructor":       completion.Instructor,
				"started_at":       completion.StartedAt,
				"completed_at":     completion.CompletedAt,
				"assessed_at":      completion.AssessedAt,
				"issued_at":        completion.IssuedAt,
				"attempt_no":       completion.AttemptNo,
				"issuer_id":        completion.IssuerID,
				"term_id":          completion.TermID,
			}
			
			studentCompletions[studentDID] = append(studentCompletions[studentDID], completionData)
		}
		
		// Add term data to each student's record
		for studentDID, courses := range studentCompletions {
			terms := allStudentData[studentDID]["terms"].(map[string]interface{})
			terms[termID] = map[string]interface{}{
				"term_id":      termID,
				"courses":      courses,
				"total_credits": len(courses) * 3, // Assuming 3 credits per course
				"gpa":          calculateGPA(courses),
			}
		}
		
		// Create term summary
		termSummaries[termID] = map[string]interface{}{
			"term_id":         termID,
			"total_students":  len(studentCompletions),
			"total_courses":   len(completions),
			"courses_offered": getUniqueCourses(completions),
		}
		
		fmt.Printf("    ✅ Organized data for %d students\n", len(studentCompletions))
	}
	
	// Step 3: Save organized data
	fmt.Println("\n💾 Step 3: Saving organized data...")
	
	// Create output directory
	outputDir := "../data/student_journeys"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Save individual student journeys
	studentsDir := filepath.Join(outputDir, "students")
	if err := os.MkdirAll(studentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create students directory: %w", err)
	}
	
	for studentDID, studentData := range allStudentData {
		filename := fmt.Sprintf("journey_%s.json", extractStudentID(studentDID))
		filepath := filepath.Join(studentsDir, filename)
		
		if err := saveStudentJSONFile(filepath, studentData); err != nil {
			return fmt.Errorf("failed to save student data for %s: %w", studentDID, err)
		}
		
		fmt.Printf("  ✅ Saved %s\n", filename)
	}
	
	// Save term summaries
	termsDir := filepath.Join(outputDir, "terms")
	if err := os.MkdirAll(termsDir, 0755); err != nil {
		return fmt.Errorf("failed to create terms directory: %w", err)
	}
	
	for termID, termData := range termSummaries {
		filename := fmt.Sprintf("summary_%s.json", termID)
		filepath := filepath.Join(termsDir, filename)
		
		if err := saveStudentJSONFile(filepath, termData); err != nil {
			return fmt.Errorf("failed to save term summary for %s: %w", termID, err)
		}
		
		fmt.Printf("  ✅ Saved %s\n", filename)
	}
	
	// Save complete system summary
	systemSummary := map[string]interface{}{
		"generation_timestamp": time.Now().Format(time.RFC3339),
		"total_students":      len(allStudentData),
		"total_terms":         len(terms),
		"terms_generated":     terms,
		"students":           getStudentIDs(allStudentData),
		"term_summaries":     termSummaries,
	}
	
	summaryPath := filepath.Join(outputDir, "system_summary.json")
	if err := saveStudentJSONFile(summaryPath, systemSummary); err != nil {
		return fmt.Errorf("failed to save system summary: %w", err)
	}
	
	fmt.Printf("  ✅ Saved system_summary.json\n")
	
	// Step 4: Display generation statistics
	fmt.Println("\n📊 Step 4: Generation Statistics")
	fmt.Printf("  📚 Total Terms: %d\n", len(terms))
	fmt.Printf("  👥 Total Students: %d\n", len(allStudentData))
	fmt.Printf("  📖 Total Course Completions: %d\n", getTotalCompletions(allStudentData))
	fmt.Printf("  📁 Output Directory: %s\n", outputDir)
	
	return nil
}

func calculateGPA(courses []interface{}) float64 {
	if len(courses) == 0 {
		return 0.0
	}
	
	gradePoints := map[string]float64{
		"A+": 4.0, "A": 4.0, "A-": 3.7,
		"B+": 3.3, "B": 3.0, "B-": 2.7,
		"C+": 2.3, "C": 2.0, "C-": 1.7,
		"D+": 1.3, "D": 1.0, "F": 0.0,
	}
	
	total := 0.0
	for _, courseInterface := range courses {
		course := courseInterface.(map[string]interface{})
		grade := course["grade"].(string)
		if points, exists := gradePoints[grade]; exists {
			total += points
		}
	}
	
	return total / float64(len(courses))
}

func getUniqueCourses(completions []verkle.CourseCompletion) []string {
	courseSet := make(map[string]bool)
	for _, completion := range completions {
		courseSet[completion.CourseID] = true
	}
	
	courses := make([]string, 0, len(courseSet))
	for course := range courseSet {
		courses = append(courses, course)
	}
	return courses
}

func extractStudentID(studentDID string) string {
	// Extract student ID from DID format "did:example:STU001"
	if len(studentDID) >= 12 {
		return studentDID[12:] // Skip "did:example:"
	}
	return studentDID
}

func getStudentIDs(allStudentData map[string]map[string]interface{}) []string {
	students := make([]string, 0, len(allStudentData))
	for studentDID := range allStudentData {
		students = append(students, extractStudentID(studentDID))
	}
	return students
}

func getTotalCompletions(allStudentData map[string]map[string]interface{}) int {
	total := 0
	for _, studentData := range allStudentData {
		terms := studentData["terms"].(map[string]interface{})
		for _, termInterface := range terms {
			term := termInterface.(map[string]interface{})
			courses := term["courses"].([]interface{})
			total += len(courses)
		}
	}
	return total
}

func saveStudentJSONFile(filepath string, data interface{}) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func init() {
	rootCmd.AddCommand(dataGeneratorCmd)
}