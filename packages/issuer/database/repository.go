package database

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReceiptRepository struct {
	db *gorm.DB
}

func NewReceiptRepository(db *gorm.DB) *ReceiptRepository {
	return &ReceiptRepository{db: db}
}

// ========== TERM RECEIPTS ==========

// StoreTermReceipt stores a single term receipt
func (r *ReceiptRepository) StoreTermReceipt(receipt *TermReceipt) error {
	return r.db.Create(receipt).Error
}

// BulkStoreTermReceipts efficiently stores multiple term receipts (for demo data)
func (r *ReceiptRepository) BulkStoreTermReceipts(receipts []*TermReceipt) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Use batch insert for better performance
		if err := tx.CreateInBatches(receipts, 100).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetTermReceipt retrieves a specific term receipt
func (r *ReceiptRepository) GetTermReceipt(receiptID string) (*TermReceipt, error) {
	var receipt TermReceipt
	err := r.db.Where("receipt_id = ?", receiptID).First(&receipt).Error
	return &receipt, err
}

// GetTermReceiptsForStudent gets all term receipts for a student
func (r *ReceiptRepository) GetTermReceiptsForStudent(studentID string) ([]*TermReceipt, error) {
	var receipts []*TermReceipt
	err := r.db.
		Where("student_id = ?", studentID).
		Order("generated_at ASC").
		Find(&receipts).Error
	return receipts, err
}

// GetTermReceiptsForStudentByTerms gets specific term receipts
func (r *ReceiptRepository) GetTermReceiptsForStudentByTerms(
	studentID string,
	termIDs []string,
) ([]*TermReceipt, error) {
	var receipts []*TermReceipt
	err := r.db.
		Where("student_id = ? AND term_id IN ?", studentID, termIDs).
		Order("generated_at ASC").
		Find(&receipts).Error
	return receipts, err
}

// ========== ACCUMULATED RECEIPTS ==========

// GenerateAccumulatedReceipt creates a new accumulated receipt from term receipts
func (r *ReceiptRepository) GenerateAccumulatedReceipt(
	studentID string,
	termIDs []string, // nil = all terms
	receiptType string, // "progress", "diploma", "custom"
) (*AccumulatedReceipt, error) {
	var accumulated *AccumulatedReceipt

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Step 1: Get term receipts
		var termReceipts []*TermReceipt
		query := tx.Where("student_id = ?", studentID)
		if termIDs != nil && len(termIDs) > 0 {
			query = query.Where("term_id IN ?", termIDs)
		}
		if err := query.Order("generated_at ASC").Find(&termReceipts).Error; err != nil {
			return err
		}

		if len(termReceipts) == 0 {
			return fmt.Errorf("no term receipts found for student %s", studentID)
		}

		// Step 2: Accumulate data
		accumulated = &AccumulatedReceipt{
			AccumulatedReceiptID: fmt.Sprintf("%s_%s_%d",
				receiptType, studentID, time.Now().Unix()),
			StudentID:   studentID,
			Type:        receiptType,
			GeneratedAt: time.Now(),
		}

		// Collect term IDs and receipt IDs
		var termIDList []string
		var receiptIDList []string
		var allCourses []interface{}
		totalCourses := 0
		totalCredits := 0

		for _, tr := range termReceipts {
			termIDList = append(termIDList, tr.TermID)
			receiptIDList = append(receiptIDList, tr.ReceiptID)
			totalCourses += tr.CourseCount

			// Extract courses from JSON
			var courses []map[string]interface{}
			if err := json.Unmarshal(tr.RevealedCourses, &courses); err == nil {
				for _, course := range courses {
					allCourses = append(allCourses, course)
					if credits, ok := course["credits"].(float64); ok {
						totalCredits += int(credits)
					}
				}
			}
		}

		// Marshal accumulated data to JSON
		termIDsJSON, _ := json.Marshal(termIDList)
		receiptIDsJSON, _ := json.Marshal(receiptIDList)
		allCoursesJSON, _ := json.Marshal(allCourses)

		accumulated.TermsIncluded = datatypes.JSON(termIDsJSON)
		accumulated.TermReceiptIDs = datatypes.JSON(receiptIDsJSON)
		accumulated.AllCourses = datatypes.JSON(allCoursesJSON)
		accumulated.TotalCourses = totalCourses
		accumulated.TotalCredits = totalCredits
		accumulated.CompletedTerms = len(termReceipts)

		// Set validity period
		accumulated.ValidFrom = termReceipts[0].GeneratedAt
		if receiptType == "diploma" {
			lastReceipt := termReceipts[len(termReceipts)-1]
			accumulated.ValidUntil = &lastReceipt.GeneratedAt
		}

		// Calculate GPA (if grade data available)
		accumulated.GPA = r.calculateGPA(allCourses)

		// Step 3: Store accumulated receipt
		if err := tx.Create(accumulated).Error; err != nil {
			return err
		}

		return nil
	})

	return accumulated, err
}

// GetAccumulatedReceipt retrieves an accumulated receipt
func (r *ReceiptRepository) GetAccumulatedReceipt(receiptID string) (*AccumulatedReceipt, error) {
	var receipt AccumulatedReceipt
	err := r.db.Where("accumulated_receipt_id = ?", receiptID).First(&receipt).Error
	return &receipt, err
}

// GetLatestDiplomaReceipt gets the most recent diploma receipt for a student
func (r *ReceiptRepository) GetLatestDiplomaReceipt(studentID string) (*AccumulatedReceipt, error) {
	var receipt AccumulatedReceipt
	err := r.db.
		Where("student_id = ? AND type = ?", studentID, "diploma").
		Order("generated_at DESC").
		First(&receipt).Error
	return &receipt, err
}

// GetCurrentProgressReceipt generates the latest progress receipt on-the-fly
func (r *ReceiptRepository) GetCurrentProgressReceipt(studentID string) (*AccumulatedReceipt, error) {
	// Always generate fresh progress receipt (ensures it's current)
	return r.GenerateAccumulatedReceipt(studentID, nil, "progress")
}

// StoreAccumulatedReceipt stores an accumulated receipt in the database
func (r *ReceiptRepository) StoreAccumulatedReceipt(receipt *AccumulatedReceipt) error {
	return r.db.Create(receipt).Error
}

// ========== VERIFICATION LOGS ==========

// LogVerification records a verification attempt
func (r *ReceiptRepository) LogVerification(log *VerificationLog) error {
	return r.db.Create(log).Error
}

// GetVerificationHistory gets verification history for a receipt
func (r *ReceiptRepository) GetVerificationHistory(receiptID string) ([]*VerificationLog, error) {
	var logs []*VerificationLog
	err := r.db.
		Where("receipt_id = ?", receiptID).
		Order("verified_at DESC").
		Limit(100). // Limit for performance
		Find(&logs).Error
	return logs, err
}

// ========== STUDENTS ==========

// CreateStudent creates a new student
func (r *ReceiptRepository) CreateStudent(student *Student) error {
	return r.db.Create(student).Error
}

// GetStudent retrieves a student by ID
func (r *ReceiptRepository) GetStudent(studentID string) (*Student, error) {
	var student Student
	err := r.db.Where("student_id = ?", studentID).First(&student).Error
	return &student, err
}

// GetAllStudents retrieves all students
func (r *ReceiptRepository) GetAllStudents() ([]*Student, error) {
	var students []*Student
	err := r.db.Find(&students).Error
	return students, err
}

// ========== TERMS ==========

// CreateTerm creates a new term
func (r *ReceiptRepository) CreateTerm(term *Term) error {
	return r.db.Create(term).Error
}

// GetTerm retrieves a term by ID
func (r *ReceiptRepository) GetTerm(termID string) (*Term, error) {
	var term Term
	err := r.db.Where("term_id = ?", termID).First(&term).Error
	return &term, err
}

// GetAllTerms retrieves all terms
func (r *ReceiptRepository) GetAllTerms() ([]*Term, error) {
	var terms []*Term
	err := r.db.Order("start_date ASC").Find(&terms).Error
	return terms, err
}

// ========== HELPER FUNCTIONS ==========

// calculateGPA calculates GPA from courses
func (r *ReceiptRepository) calculateGPA(courses []interface{}) float64 {
	totalPoints := 0.0
	totalCredits := 0.0

	for _, courseInterface := range courses {
		course, ok := courseInterface.(map[string]interface{})
		if !ok {
			continue
		}

		grade, hasGrade := course["grade"].(string)
		credits, hasCredits := course["credits"].(float64)

		if hasGrade && hasCredits {
			gradePoint := r.gradeToPoint(grade)
			totalPoints += gradePoint * credits
			totalCredits += credits
		}
	}

	if totalCredits == 0 {
		return 0
	}
	return totalPoints / totalCredits
}

// gradeToPoint converts letter grade to GPA point
func (r *ReceiptRepository) gradeToPoint(grade string) float64 {
	gradeMap := map[string]float64{
		"A":  4.0,
		"A-": 3.7,
		"B+": 3.3,
		"B":  3.0,
		"B-": 2.7,
		"C+": 2.3,
		"C":  2.0,
		"C-": 1.7,
		"D":  1.0,
		"F":  0.0,
	}

	if point, ok := gradeMap[grade]; ok {
		return point
	}
	return 0.0
}
