package main

import (
	"encoding/json"
	"fmt"
	"iumicert/issuer/database"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"gorm.io/datatypes"
)

var dbImportCmd = &cobra.Command{
	Use:   "db-import",
	Short: "Import generated data into database",
	Long:  "Import students, terms, and receipts from JSON files into PostgreSQL database",
	Run: func(cmd *cobra.Command, args []string) {
		runDBImport()
	},
}

func init() {
	rootCmd.AddCommand(dbImportCmd)
}

func runDBImport() {
	fmt.Println("📥 Importing data into database...")

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	repo := database.NewReceiptRepository(db)

	// Step 1: Import Students
	fmt.Println("\n👥 Step 1: Importing students...")
	if err := importStudents(repo); err != nil {
		log.Printf("⚠️  Warning: %v", err)
	}

	// Step 2: Import Terms
	fmt.Println("\n📚 Step 2: Importing terms...")
	if err := importTerms(repo); err != nil {
		log.Printf("⚠️  Warning: %v", err)
	}

	// Step 3: Import Term Receipts
	fmt.Println("\n🎓 Step 3: Importing term receipts...")
	if err := importTermReceipts(repo); err != nil {
		log.Fatalf("❌ Failed to import receipts: %v", err)
	}

	// Step 4: Generate Accumulated Receipts
	fmt.Println("\n📚 Step 4: Generating accumulated receipts (full academic journey)...")
	if err := generateAccumulatedReceipts(repo); err != nil {
		log.Printf("⚠️  Warning: %v", err)
	}

	fmt.Println("\n✅ Database import completed successfully!")
	fmt.Println("")
	fmt.Println("📊 Summary:")

	// Count records
	var studentCount, termCount, receiptCount, accumulatedCount int64
	db.Model(&database.Student{}).Count(&studentCount)
	db.Model(&database.Term{}).Count(&termCount)
	db.Model(&database.TermReceipt{}).Count(&receiptCount)
	db.Model(&database.AccumulatedReceipt{}).Count(&accumulatedCount)

	fmt.Printf("  • Students: %d\n", studentCount)
	fmt.Printf("  • Terms: %d\n", termCount)
	fmt.Printf("  • Term Receipts: %d\n", receiptCount)
	fmt.Printf("  • Accumulated Receipts: %d\n", accumulatedCount)
}

func importStudents(repo *database.ReceiptRepository) error {
	studentFiles, err := filepath.Glob("../data/student_journeys/students/journey_*.json")
	if err != nil {
		return fmt.Errorf("failed to find student files: %w", err)
	}

	if len(studentFiles) == 0 {
		return fmt.Errorf("no student files found in ../data/student_journeys/students/")
	}

	for _, file := range studentFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("  ⚠️  Failed to read %s: %v", file, err)
			continue
		}

		var journey map[string]interface{}
		if err := json.Unmarshal(data, &journey); err != nil {
			log.Printf("  ⚠️  Failed to parse %s: %v", file, err)
			continue
		}

		studentID := journey["student_id"].(string)

		// Extract just the student ID (remove any existing DID prefix)
		cleanStudentID := studentID
		if len(studentID) > 12 && studentID[:12] == "did:example:" {
			cleanStudentID = studentID[12:]
		}

		// Create student record with did:key format (cryptographic key-based DID)
		student := &database.Student{
			StudentID:          cleanStudentID,
			Name:               fmt.Sprintf("Student %s", cleanStudentID),
			DID:                fmt.Sprintf("did:key:%s", cleanStudentID),
			EnrollmentDate:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			ExpectedGraduation: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			Status:             "active",
		}

		if err := repo.CreateStudent(student); err != nil {
			// Ignore duplicate errors (student already exists)
			if err.Error() != "ERROR: duplicate key value violates unique constraint \"idx_students_student_id\" (SQLSTATE 23505)" {
				log.Printf("  ⚠️  Failed to create student %s: %v", studentID, err)
			}
			continue
		}

		fmt.Printf("  ✓ Imported student: %s\n", studentID)
	}

	return nil
}

func importTerms(repo *database.ReceiptRepository) error {
	rootFiles, err := filepath.Glob("../publish_ready/roots/root_*.json")
	if err != nil {
		return fmt.Errorf("failed to find root files: %w", err)
	}

	if len(rootFiles) == 0 {
		return fmt.Errorf("no root files found in ../publish_ready/roots/")
	}

	for _, file := range rootFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("  ⚠️  Failed to read %s: %v", file, err)
			continue
		}

		var rootData map[string]interface{}
		if err := json.Unmarshal(data, &rootData); err != nil {
			log.Printf("  ⚠️  Failed to parse %s: %v", file, err)
			continue
		}

		termID := rootData["term_id"].(string)
		verkleRootHex := rootData["verkle_root"].(string)

		// Convert hex to bytes
		verkleRootBytes, err := hexToBytes(verkleRootHex)
		if err != nil {
			log.Printf("  ⚠️  Invalid verkle root for %s: %v", termID, err)
			continue
		}

		// Create term record
		term := &database.Term{
			TermID:          termID,
			StartDate:       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), // Placeholder
			EndDate:         time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC), // Placeholder
			VerkleRootHex:   verkleRootHex,
			VerkleRootBytes: verkleRootBytes,
		}

		if err := repo.CreateTerm(term); err != nil {
			// Ignore duplicate errors
			if err.Error() != "ERROR: duplicate key value violates unique constraint \"idx_terms_term_id\" (SQLSTATE 23505)" {
				log.Printf("  ⚠️  Failed to create term %s: %v", termID, err)
			}
			continue
		}

		fmt.Printf("  ✓ Imported term: %s\n", termID)
	}

	return nil
}

func importTermReceipts(repo *database.ReceiptRepository) error {
	receiptFiles, err := filepath.Glob("../publish_ready/receipts/*_journey.json")
	if err != nil {
		return fmt.Errorf("failed to find receipt files: %w", err)
	}

	if len(receiptFiles) == 0 {
		return fmt.Errorf("no receipt files found in ../publish_ready/receipts/")
	}

	totalReceipts := 0

	for _, file := range receiptFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("  ⚠️  Failed to read %s: %v", file, err)
			continue
		}

		var receipt map[string]interface{}
		if err := json.Unmarshal(data, &receipt); err != nil {
			log.Printf("  ⚠️  Failed to parse %s: %v", file, err)
			continue
		}

		studentID := receipt["student_id"].(string)
		termReceipts := receipt["term_receipts"].(map[string]interface{})

		// Import each term receipt
		for termID, termData := range termReceipts {
			termMap := termData.(map[string]interface{})
			receiptData := termMap["receipt"].(map[string]interface{})

			// Extract data
			courseProofs := receiptData["course_proofs"]
			revealedCourses := receiptData["revealed_courses"]
			verkleRootHex := receiptData["verkle_root"].(string)

			// Marshal to JSON
			verkleProofJSON, _ := json.Marshal(courseProofs)
			revealedCoursesJSON, _ := json.Marshal(revealedCourses)

			// Create term receipt
			termReceipt := &database.TermReceipt{
				ReceiptID:       fmt.Sprintf("receipt_%s_%s_%d", studentID, termID, time.Now().Unix()),
				StudentID:       studentID,
				TermID:          termID,
				VerkleProof:     datatypes.JSON(verkleProofJSON),
				RevealedCourses: datatypes.JSON(revealedCoursesJSON),
				StateDiff:       datatypes.JSON("[]"), // Placeholder
				CourseCount:     len(revealedCourses.([]interface{})),
				VerkleRootHex:   verkleRootHex,
				GeneratedAt:     time.Now(),
				IsSelective:     false,
			}

			if err := repo.StoreTermReceipt(termReceipt); err != nil {
				log.Printf("  ⚠️  Failed to store receipt for %s/%s: %v", studentID, termID, err)
				continue
			}

			totalReceipts++
		}

		fmt.Printf("  ✓ Imported receipts for student: %s\n", studentID)
	}

	fmt.Printf("\n  📊 Total term receipts imported: %d\n", totalReceipts)
	return nil
}

// Helper function to convert hex string to bytes
func hexToBytes(hexStr string) ([]byte, error) {
	// Remove 0x prefix if present
	if len(hexStr) >= 2 && hexStr[0:2] == "0x" {
		hexStr = hexStr[2:]
	}

	bytes := make([]byte, len(hexStr)/2)
	for i := 0; i < len(bytes); i++ {
		fmt.Sscanf(hexStr[i*2:i*2+2], "%02x", &bytes[i])
	}

	return bytes, nil
}

// generateAccumulatedReceipts creates full academic journey receipts for each student
func generateAccumulatedReceipts(repo *database.ReceiptRepository) error {
	// Get all students
	students, err := repo.GetAllStudents()
	if err != nil {
		return fmt.Errorf("failed to fetch students: %w", err)
	}

	generatedCount := 0
	for _, student := range students {
		// Strip "did:example:" prefix if present to match term_receipts format
		studentID := student.StudentID
		if len(studentID) > 12 && studentID[:12] == "did:example:" {
			studentID = studentID[12:]
		}

		// Generate accumulated receipt for full journey (all terms)
		accumulated, err := repo.GenerateAccumulatedReceipt(studentID, nil, "progress")
		if err != nil {
			log.Printf("  ⚠️  Failed to generate accumulated receipt for %s: %v", studentID, err)
			continue
		}

		// Store the accumulated receipt
		if err := repo.StoreAccumulatedReceipt(accumulated); err != nil {
			log.Printf("  ⚠️  Failed to store accumulated receipt for %s: %v", student.StudentID, err)
			continue
		}

		generatedCount++
		fmt.Printf("  ✓ Generated accumulated receipt for: %s (%d terms, %d courses)\n",
			student.StudentID, accumulated.CompletedTerms, accumulated.TotalCourses)
	}

	fmt.Printf("\n  📊 Total accumulated receipts generated: %d\n", generatedCount)
	return nil
}
