package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"iumicert/crypto/merkle"
	"iumicert/crypto/verkle"
)

var cryptoIntegrationCmd = &cobra.Command{
	Use:   "test-crypto",
	Short: "Validate cryptographic tree integration and verification",
	Long:  `Execute comprehensive testing of Merkle-Verkle tree hybrid architecture with proof generation and validation`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🌳 Starting Merkle/Verkle Integration Test")
		fmt.Println("=" + string(make([]byte, 50)))
		
		dataDir, _ := cmd.Flags().GetString("data-dir")
		if dataDir == "" {
			dataDir = "data/generated_student_data"
		}
		
		if err := runMerkleVerkleIntegration(dataDir); err != nil {
			log.Fatalf("❌ Integration test failed: %v", err)
		}
		
		fmt.Println("\n🎉 Integration test completed successfully!")
	},
}

func runMerkleVerkleIntegration(dataDir string) error {
	// Step 1: Load student data from files
	fmt.Println("\n📊 Step 1: Loading student data...")
	studentData, err := loadStudentDataFromFiles(dataDir)
	if err != nil {
		return fmt.Errorf("failed to load student data: %w", err)
	}
	
	fmt.Printf("✅ Loaded data for %d students\n", len(studentData))
	
	// Step 2: Process each term separately
	fmt.Println("\n🌲 Step 2: Building Merkle trees for each student-term...")
	
	// Group data by term for Verkle tree processing
	termData := make(map[string]map[string][]merkle.CourseCompletion)
	studentTrees := make(map[string]map[string]*merkle.StudentTermMerkle) // [termID][studentID]
	
	for studentID, student := range studentData {
		studentTrees[studentID] = make(map[string]*merkle.StudentTermMerkle)
		
		terms := student["terms"].(map[string]interface{})
		for termID, termInterface := range terms {
			term := termInterface.(map[string]interface{})
			courses := term["courses"].([]interface{})
			
			// Convert course data back to CourseCompletion structs
			var completions []merkle.CourseCompletion
			for _, courseInterface := range courses {
				course := courseInterface.(map[string]interface{})
				completion, err := convertToCourseCompletion(course)
				if err != nil {
					return fmt.Errorf("failed to convert course data: %w", err)
				}
				completions = append(completions, completion)
			}
			
			// Create student-term Merkle tree
			tree, err := merkle.NewStudentTermMerkle(studentID, termID, completions)
			if err != nil {
				return fmt.Errorf("failed to create merkle tree for %s/%s: %w", studentID, termID, err)
			}
			
			studentTrees[studentID][termID] = tree
			
			// Organize for term-level processing
			if termData[termID] == nil {
				termData[termID] = make(map[string][]merkle.CourseCompletion)
			}
			termData[termID][studentID] = completions
			
			fmt.Printf("  ✓ Student %s, Term %s: %d courses, root: %x\n", 
				studentID, termID, len(completions), tree.Root[:8])
		}
	}
	
	// Step 3: Test Merkle proof verification
	fmt.Println("\n🔍 Step 3: Testing Merkle proof verification...")
	
	for studentID, termTrees := range studentTrees {
		for termID, tree := range termTrees {
			// Test proof verification for each course
			for _, courseID := range tree.ListCourses() {
				valid, err := tree.VerifyProofForCourse(courseID, tree.Root)
				if err != nil {
					return fmt.Errorf("failed to verify merkle proof for %s/%s/%s: %w", studentID, termID, courseID, err)
				}
				if !valid {
					return fmt.Errorf("merkle proof verification failed for %s/%s/%s", studentID, termID, courseID)
				}
			}
			fmt.Printf("  ✅ All Merkle proofs verified for %s/%s\n", studentID, termID)
		}
	}
	
	// Step 4: Build Verkle trees for each term
	fmt.Println("\n🌳 Step 4: Building Verkle trees for each term...")
	
	termVerkleTrees := make(map[string]*verkle.TermVerkleTree)
	
	for termID, students := range termData {
		fmt.Printf("\n  📚 Processing term: %s\n", termID)
		
		termTree := verkle.NewTermVerkleTree(termID)
		
		for studentID, completions := range students {
			err := termTree.AddStudent(studentID, completions)
			if err != nil {
				return fmt.Errorf("failed to add student %s to verkle tree: %w", studentID, err)
			}
		}
		
		// Publish the term
		err := termTree.PublishTerm()
		if err != nil {
			return fmt.Errorf("failed to publish term %s: %w", termID, err)
		}
		
		termVerkleTrees[termID] = termTree
		fmt.Printf("    ✅ Term %s published with Verkle root: %x\n", termID, termTree.VerkleRoot[:8])
	}
	
	// Step 5: Test receipt generation and verification
	fmt.Println("\n📄 Step 5: Testing receipt generation and verification...")
	
	for termID, termTree := range termVerkleTrees {
		fmt.Printf("\n  📋 Testing receipts for term: %s\n", termID)
		
		for studentID := range termData[termID] {
			// Generate full receipt
			fullReceipt, err := termTree.GenerateVerificationReceipt(studentID, nil)
			if err != nil {
				return fmt.Errorf("failed to generate full receipt for %s/%s: %w", studentID, termID, err)
			}
			
			fmt.Printf("    ✓ Generated full receipt for %s (%d courses)\n", 
				studentID, len(fullReceipt.RevealedCourses))
			
			// Verify receipt off-chain
			result, err := verkle.VerifyReceiptOffChain(fullReceipt, termTree.VerkleRoot)
			if err != nil {
				return fmt.Errorf("failed to verify receipt for %s/%s: %w", studentID, termID, err)
			}
			
			if !result.Valid {
				return fmt.Errorf("receipt verification failed for %s/%s: %v", studentID, termID, result.Errors)
			}
			
			fmt.Printf("    ✅ Receipt verified successfully for %s\n", studentID)
			
			// Test selective disclosure (first 2 courses only)
			studentCourses := termData[termID][studentID]
			if len(studentCourses) >= 2 {
				courseIDs := []string{
					studentCourses[0].CourseID,
					studentCourses[1].CourseID,
				}
				
				selectiveReceipt, err := termTree.GenerateVerificationReceipt(studentID, courseIDs)
				if err != nil {
					return fmt.Errorf("failed to generate selective receipt for %s/%s: %w", studentID, termID, err)
				}
				
				fmt.Printf("    ✓ Generated selective receipt for %s (%d/%d courses)\n", 
					studentID, len(selectiveReceipt.RevealedCourses), len(studentCourses))
				
				// Verify selective receipt
				result, err := verkle.VerifyReceiptOffChain(selectiveReceipt, termTree.VerkleRoot)
				if err != nil {
					return fmt.Errorf("failed to verify selective receipt for %s/%s: %w", studentID, termID, err)
				}
				
				if !result.Valid {
					return fmt.Errorf("selective receipt verification failed for %s/%s: %v", studentID, termID, result.Errors)
				}
				
				fmt.Printf("    ✅ Selective receipt verified successfully for %s\n", studentID)
			}
		}
	}
	
	// Step 6: Save integration results
	fmt.Println("\n💾 Step 6: Saving integration results...")
	
	outputDir := filepath.Join(dataDir, "integration_results")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Save term Verkle tree summaries
	for termID, termTree := range termVerkleTrees {
		summary := map[string]interface{}{
			"term_id":      termID,
			"verkle_root":  fmt.Sprintf("%x", termTree.VerkleRoot),
			"total_students": len(termData[termID]),
			"students":     getStudentIDsFromTerm(termData[termID]),
		}
		
		filename := fmt.Sprintf("verkle_summary_%s.json", termID)
		filepath := filepath.Join(outputDir, filename)
		
		if err := saveIntegrationJSONFile(filepath, summary); err != nil {
			return fmt.Errorf("failed to save verkle summary for %s: %w", termID, err)
		}
		
		fmt.Printf("  ✅ Saved %s\n", filename)
	}
	
	// Step 7: Display integration statistics
	fmt.Println("\n📊 Step 7: Integration Statistics")
	fmt.Printf("  📚 Total Terms Processed: %d\n", len(termVerkleTrees))
	fmt.Printf("  👥 Total Students: %d\n", len(studentData))
	fmt.Printf("  🌲 Merkle Trees Built: %d\n", countMerkleTrees(studentTrees))
	fmt.Printf("  🌳 Verkle Trees Built: %d\n", len(termVerkleTrees))
	fmt.Printf("  📁 Results Directory: %s\n", outputDir)
	
	return nil
}

func loadStudentDataFromFiles(dataDir string) (map[string]map[string]interface{}, error) {
	studentsDir := filepath.Join(dataDir, "students")
	
	// Check if directory exists
	if _, err := os.Stat(studentsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("students directory not found: %s", studentsDir)
	}
	
	// Read all student files
	files, err := os.ReadDir(studentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read students directory: %w", err)
	}
	
	studentData := make(map[string]map[string]interface{})
	
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		filepath := filepath.Join(studentsDir, file.Name())
		data, err := os.ReadFile(filepath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", filepath, err)
		}
		
		var student map[string]interface{}
		if err := json.Unmarshal(data, &student); err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s: %w", filepath, err)
		}
		
		studentID := student["student_id"].(string)
		studentData[studentID] = student
	}
	
	return studentData, nil
}

func convertToCourseCompletion(course map[string]interface{}) (merkle.CourseCompletion, error) {
	// Note: This is a simplified conversion that creates valid CourseCompletion structs
	// The Merkle/Verkle integration doesn't need the exact original data,
	// just a valid structure that matches the schema
	
	return merkle.CourseCompletion{
		IssuerID:    getStringValue(course, "issuer_id", "IU-CS"),
		StudentID:   getStringValue(course, "student_id", "STU001"),
		TermID:      getStringValue(course, "term_id", "Fall_2024"),
		CourseID:    getStringValue(course, "course_id", "CS101"),
		CourseName:  getStringValue(course, "course_name", "Test Course"),
		AttemptNo:   getUint8Value(course, "attempt_no", 1),
		StartedAt:   getTimeValue(course, "started_at"),
		CompletedAt: getTimeValue(course, "completed_at"),
		AssessedAt:  getTimeValue(course, "assessed_at"),
		IssuedAt:    getTimeValue(course, "issued_at"),
		Grade:       getStringValue(course, "grade", "B"),
		Credits:     getUint8Value(course, "credits", 3),
		Instructor:  getStringValue(course, "instructor", "Prof. Test"),
	}, nil
}

func getStringValue(data map[string]interface{}, key, defaultValue string) string {
	if val, exists := data[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getUint8Value(data map[string]interface{}, key string, defaultValue uint8) uint8 {
	if val, exists := data[key]; exists {
		if num, ok := val.(float64); ok {
			return uint8(num)
		}
	}
	return defaultValue
}

func getTimeValue(data map[string]interface{}, key string) time.Time {
	if val, exists := data[key]; exists {
		if timeStr, ok := val.(string); ok {
			if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
				return t
			}
		}
	}
	// Return a default time if parsing fails
	return time.Now().AddDate(0, 0, -30) // 30 days ago as default
}

func getStudentIDsFromTerm(termData map[string][]merkle.CourseCompletion) []string {
	students := make([]string, 0, len(termData))
	for studentID := range termData {
		students = append(students, studentID)
	}
	return students
}

func countMerkleTrees(studentTrees map[string]map[string]*merkle.StudentTermMerkle) int {
	count := 0
	for _, termTrees := range studentTrees {
		count += len(termTrees)
	}
	return count
}

func saveIntegrationJSONFile(filepath string, data interface{}) error {
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
	cryptoIntegrationCmd.Flags().String("data-dir", "data/generated_student_data", "Directory containing generated student data")
	rootCmd.AddCommand(cryptoIntegrationCmd)
}