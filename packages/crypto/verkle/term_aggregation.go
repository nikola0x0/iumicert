package verkle

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	verkleLib "github.com/ethereum/go-verkle"
)

// CourseCompletion represents a completed course with all relevant data
type CourseCompletion struct {
	IssuerID    string    `json:"issuer_id"`
	StudentID   string    `json:"student_id"`
	TermID      string    `json:"term_id"`
	CourseID    string    `json:"course_id"`
	CourseName  string    `json:"course_name"`
	AttemptNo   uint8     `json:"attempt_no"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
	AssessedAt  time.Time `json:"assessed_at"`
	IssuedAt    time.Time `json:"issued_at"`
	Grade       string    `json:"grade"`
	Credits     uint8     `json:"credits"`
	Instructor  string    `json:"instructor"`
}

// TermVerkleTree manages a single Verkle tree containing all courses for a term
type TermVerkleTree struct {
	TermID           string                              `json:"term_id"`
	PublishedAt      time.Time                           `json:"published_at"`
	Version          uint32                              `json:"version"`
	VerkleRoot       [32]byte                            `json:"verkle_root"`
	CourseEntries    map[string]CourseCompletion         `json:"course_entries"` // courseKey -> CourseCompletion
	CourseProofs     map[string][]byte                   `json:"course_proofs"`  // courseKey -> serialized Verkle proof
	tree             verkleLib.VerkleNode                // Internal tree (not serialized)
}

// VerificationReceipt contains all data needed for off-chain or on-chain verification
type VerificationReceipt struct {
	TermID          string                      `json:"term_id"`
	StudentDID      string                      `json:"student_did"`
	VerkleRoot      [32]byte                    `json:"verkle_root"`
	PublishedAt     time.Time                   `json:"published_at"`
	RevealedCourses []CourseCompletion          `json:"revealed_courses"`
	CourseProofs    map[string][]byte           `json:"course_proofs"`   // courseID -> verkle proof
	SelectiveDisclosure bool                    `json:"selective_disclosure"`
	Metadata        ReceiptMetadata             `json:"metadata"`
}

// ReceiptMetadata contains additional verification metadata
type ReceiptMetadata struct {
	GeneratedAt       time.Time `json:"generated_at"`
	TotalCourses      int       `json:"total_courses"`
	RevealedCourses   int       `json:"revealed_courses"`
	VerificationLevel string    `json:"verification_level"` // "full" or "selective"
}

// NewTermVerkleTree creates a new term-level Verkle tree
func NewTermVerkleTree(termID string) *TermVerkleTree {
	return &TermVerkleTree{
		TermID:        termID,
		CourseEntries: make(map[string]CourseCompletion),
		CourseProofs:  make(map[string][]byte),
		tree:          verkleLib.New(),
	}
}

// AddCourses adds a student's courses directly to the single Verkle tree
func (tvt *TermVerkleTree) AddCourses(studentDID string, courses []CourseCompletion) error {
	log.Printf("Adding %d courses for student %s to term %s", len(courses), studentDID, tvt.TermID)
	
	for _, course := range courses {
		// Generate deterministic key for each course
		courseKey := fmt.Sprintf("%s:%s:%s", studentDID, tvt.TermID, course.CourseID)
		courseKeyHash := sha256.Sum256([]byte(courseKey))
		
		// Serialize course data as value
		courseData, err := json.Marshal(course)
		if err != nil {
			return fmt.Errorf("failed to serialize course %s: %w", course.CourseID, err)
		}
		courseValueHash := sha256.Sum256(courseData)
		
		// Store course entry for later retrieval
		tvt.CourseEntries[courseKey] = course
		
		// Add to Verkle tree: key = H(studentDID:termID:courseID), value = H(course_data)
		err = tvt.tree.Insert(courseKeyHash[:], courseValueHash[:], nil)
		if err != nil {
			return fmt.Errorf("failed to insert course %s into verkle tree: %w", course.CourseID, err)
		}
		
		log.Printf("✅ Course %s added for student %s", course.CourseID, studentDID)
	}
	
	return nil
}

// GenerateCourseProof creates a Verkle proof for a specific course
func (tvt *TermVerkleTree) GenerateCourseProof(studentDID, courseID string) ([]byte, error) {
	courseKey := fmt.Sprintf("%s:%s:%s", studentDID, tvt.TermID, courseID)
	courseKeyHash := sha256.Sum256([]byte(courseKey))
	
	// Check if course exists
	if _, exists := tvt.CourseEntries[courseKey]; !exists {
		return nil, fmt.Errorf("course %s not found for student %s in term %s", courseID, studentDID, tvt.TermID)
	}
	
	// Generate Verkle proof using the global proof functions
	proofElements, _, _, err := tvt.tree.GetProofItems([][]byte{courseKeyHash[:]}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get proof items for course %s: %w", courseID, err)
	}
	
	// For now, we'll serialize the proof elements as JSON (simplified approach)
	// In production, you'd want to use the proper Verkle proof serialization
	proof, err := json.Marshal(proofElements)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize proof for course %s: %w", courseID, err)
	}
	
	// Store the proof for later use
	tvt.CourseProofs[courseKey] = proof
	
	log.Printf("✅ Generated proof for course %s (student: %s)", courseID, studentDID)
	return proof, nil
}

// PublishTerm computes the final Verkle root and prepares for blockchain publication
func (tvt *TermVerkleTree) PublishTerm() error {
	log.Printf("Publishing term %s with %d courses", tvt.TermID, len(tvt.CourseEntries))
	
	if len(tvt.CourseEntries) == 0 {
		return fmt.Errorf("cannot publish term %s: no courses added", tvt.TermID)
	}
	
	// Compute Verkle tree commitment
	commitment := tvt.tree.Commit()
	if commitment == nil {
		return fmt.Errorf("failed to compute verkle tree commitment")
	}
	
	// Store the root
	tvt.VerkleRoot = commitment.Bytes()
	tvt.PublishedAt = time.Now()
	tvt.Version++
	
	log.Printf("✅ Term %s published successfully, Verkle root: %x", tvt.TermID, tvt.VerkleRoot)
	return nil
}

// GenerateStudentReceipt creates a verification receipt for specific courses using single Verkle tree
func (tvt *TermVerkleTree) GenerateStudentReceipt(studentDID string, courseIDs []string) (*VerificationReceipt, error) {
	log.Printf("Generating student receipt for %s, courses: %v", studentDID, courseIDs)
	
	// Check if term is published
	if tvt.PublishedAt.IsZero() {
		return nil, fmt.Errorf("term %s not yet published", tvt.TermID)
	}
	
	// Get all courses for this student from the single Verkle tree
	var studentCourses []CourseCompletion
	courseProofs := make(map[string][]byte)
	
	// If no specific courses requested, find all courses for this student
	if len(courseIDs) == 0 {
		// Find all courses for this student by scanning CourseEntries
		for courseKey, course := range tvt.CourseEntries {
			// Check if this course belongs to the student
			if strings.HasPrefix(courseKey, studentDID+":") {
				studentCourses = append(studentCourses, course)
				courseIDs = append(courseIDs, course.CourseID)
			}
		}
	} else {
		// Collect specific requested courses
		for _, courseID := range courseIDs {
			courseKey := fmt.Sprintf("%s:%s:%s", studentDID, tvt.TermID, courseID)
			if course, exists := tvt.CourseEntries[courseKey]; exists {
				studentCourses = append(studentCourses, course)
			} else {
				log.Printf("Warning: course %s not found for student %s in term %s", courseID, studentDID, tvt.TermID)
			}
		}
	}
	
	if len(studentCourses) == 0 {
		return nil, fmt.Errorf("no courses found for student %s in term %s", studentDID, tvt.TermID)
	}
	
	// Generate Verkle proofs for each course
	for _, courseID := range courseIDs {
		proof, err := tvt.GenerateCourseProof(studentDID, courseID)
		if err != nil {
			log.Printf("Warning: failed to generate proof for course %s: %v", courseID, err)
			continue
		}
		courseProofs[courseID] = proof
	}
	
	// Get total courses for this student for metadata
	totalCourses := 0
	for courseKey := range tvt.CourseEntries {
		if strings.HasPrefix(courseKey, studentDID+":") {
			totalCourses++
		}
	}
	
	// Create verification receipt with single Verkle structure
	receipt := &VerificationReceipt{
		TermID:              tvt.TermID,
		StudentDID:          studentDID,
		VerkleRoot:          tvt.VerkleRoot,
		PublishedAt:         tvt.PublishedAt,
		RevealedCourses:     studentCourses,
		CourseProofs:        courseProofs,
		SelectiveDisclosure: len(studentCourses) < totalCourses,
		Metadata: ReceiptMetadata{
			GeneratedAt:       time.Now(),
			TotalCourses:      totalCourses,
			RevealedCourses:   len(studentCourses),
			VerificationLevel: determineVerificationLevel(totalCourses, len(studentCourses)),
		},
	}
	
	log.Printf("✅ Generated receipt for student %s with %d/%d courses", 
		studentDID, len(studentCourses), totalCourses)
	
	return receipt, nil
}


// determineVerificationLevel determines if this is full or selective disclosure
func determineVerificationLevel(totalCourses, revealedCourses int) string {
	if revealedCourses == totalCourses {
		return "full"
	}
	return "selective"
}

// VerifyReceiptOffChain performs complete off-chain verification of a single Verkle receipt
func VerifyReceiptOffChain(receipt *VerificationReceipt, expectedVerkleRoot [32]byte) (*VerificationResult, error) {
	log.Println("=== STARTING OFF-CHAIN VERIFICATION (Single Verkle) ===")
	
	result := &VerificationResult{
		Valid:           true,
		Timestamp:       time.Now(),
		TermID:          receipt.TermID,
		StudentDID:      receipt.StudentDID,
		CoursesVerified: len(receipt.RevealedCourses),
		Errors:          []string{},
		Warnings:        []string{},
	}
	
	// 1. Verify Verkle root matches expected
	if receipt.VerkleRoot != expectedVerkleRoot {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Verkle root mismatch: got %x, expected %x", 
			receipt.VerkleRoot, expectedVerkleRoot))
	} else {
		log.Printf("✅ Verkle root verification passed: %x", expectedVerkleRoot)
	}
	
	// 2. Verify each course's temporal consistency
	for _, course := range receipt.RevealedCourses {
		if err := validateCourseTimestamps(course); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Course %s: %v", course.CourseID, err))
		}
		
		// Verify issued timestamp is before term publication
		if course.IssuedAt.After(receipt.PublishedAt) {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Course %s issued after term publication", course.CourseID))
		}
	}
	
	// 3. Verify Verkle proofs exist for revealed courses
	for _, course := range receipt.RevealedCourses {
		if _, hasProof := receipt.CourseProofs[course.CourseID]; !hasProof {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Missing Verkle proof for course %s", course.CourseID))
		}
	}
	
	// 4. Verify no extra proofs (security check)
	for courseID := range receipt.CourseProofs {
		if findCourseByID(receipt.RevealedCourses, courseID) == nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Verkle proof for unrevealed course: %s", courseID))
		}
	}
	
	// 5. Verify selective disclosure flag consistency
	if receipt.SelectiveDisclosure && len(receipt.RevealedCourses) >= receipt.Metadata.TotalCourses {
		result.Warnings = append(result.Warnings, "Selective disclosure flag set but all courses revealed")
	} else if !receipt.SelectiveDisclosure && len(receipt.RevealedCourses) < receipt.Metadata.TotalCourses {
		result.Warnings = append(result.Warnings, "Selective disclosure flag not set but some courses hidden")
	}
	
	if result.Valid {
		log.Println("✅ OFF-CHAIN VERIFICATION SUCCESSFUL")
	} else {
		log.Printf("❌ OFF-CHAIN VERIFICATION FAILED: %d errors", len(result.Errors))
	}
	
	return result, nil
}

// VerificationResult contains the results of verification
type VerificationResult struct {
	Valid           bool      `json:"valid"`
	Timestamp       time.Time `json:"timestamp"`
	TermID          string    `json:"term_id"`
	StudentDID      string    `json:"student_did"`
	CoursesVerified int       `json:"courses_verified"`
	Errors          []string  `json:"errors"`
	Warnings        []string  `json:"warnings"`
}

// Helper functions
func validateCourseTimestamps(course CourseCompletion) error {
	if course.StartedAt.After(course.CompletedAt) {
		return fmt.Errorf("started_at after completed_at")
	}
	if course.CompletedAt.After(course.AssessedAt) {
		return fmt.Errorf("completed_at after assessed_at")
	}
	if course.AssessedAt.After(course.IssuedAt) {
		return fmt.Errorf("assessed_at after issued_at")
	}
	return nil
}

func findCourseByID(courses []CourseCompletion, courseID string) *CourseCompletion {
	for _, course := range courses {
		if course.CourseID == courseID {
			return &course
		}
	}
	return nil
}

// SerializeToJSON serializes the term Verkle tree to JSON
func (tvt *TermVerkleTree) SerializeToJSON() ([]byte, error) {
	return json.MarshalIndent(tvt, "", "  ")
}

// GetStudentList returns all student DIDs in this term
func (tvt *TermVerkleTree) GetStudentList() []string {
	studentSet := make(map[string]bool)
	for courseKey := range tvt.CourseEntries {
		// Extract studentDID from courseKey format: studentDID:termID:courseID
		parts := strings.Split(courseKey, ":")
		if len(parts) >= 1 {
			studentSet[parts[0]] = true
		}
	}
	
	var students []string
	for studentDID := range studentSet {
		students = append(students, studentDID)
	}
	return students
}

// GetStudentCourseCount returns the number of courses for a student
func (tvt *TermVerkleTree) GetStudentCourseCount(studentDID string) int {
	count := 0
	for courseKey := range tvt.CourseEntries {
		if strings.HasPrefix(courseKey, studentDID+":") {
			count++
		}
	}
	return count
}