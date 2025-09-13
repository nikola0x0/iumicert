package verkle

import (
	"bytes"
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

// VerkleProofBundle holds all data needed for cryptographic verification
type VerkleProofBundle struct {
	VerkleProof *verkleLib.VerkleProof `json:"verkle_proof"`
	StateDiff   verkleLib.StateDiff    `json:"state_diff"`
	CourseKey   string                 `json:"course_key"`
	CourseID    string                 `json:"course_id"`
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

// GenerateCourseProof creates a proper cryptographic Verkle proof for a specific course
func (tvt *TermVerkleTree) GenerateCourseProof(studentDID, courseID string) ([]byte, error) {
	courseKey := fmt.Sprintf("%s:%s:%s", studentDID, tvt.TermID, courseID)
	courseKeyHash := sha256.Sum256([]byte(courseKey))
	
	// Check if course exists
	if _, exists := tvt.CourseEntries[courseKey]; !exists {
		return nil, fmt.Errorf("course %s not found for student %s in term %s", courseID, studentDID, tvt.TermID)
	}
	
	// Generate proper Verkle proof using MakeVerkleMultiProof (following Duc's approach)
	proof, _, _, _, err := verkleLib.MakeVerkleMultiProof(tvt.tree, nil, [][]byte{courseKeyHash[:]}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verkle proof for course %s: %w", courseID, err)
	}
	
	// Serialize the proof using proper Verkle serialization
	verkleProof, stateDiff, err := verkleLib.SerializeProof(proof)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize verkle proof for course %s: %w", courseID, err)
	}
	
	// Create proof bundle with all necessary verification data
	proofBundle := VerkleProofBundle{
		VerkleProof: verkleProof,
		StateDiff:   stateDiff,
		CourseKey:   courseKey,
		CourseID:    courseID,
	}
	
	// Serialize the bundle for storage
	proofData, err := json.Marshal(proofBundle)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize proof bundle for course %s: %w", courseID, err)
	}
	
	// Store the proof for later use
	tvt.CourseProofs[courseKey] = proofData
	
	log.Printf("✅ Generated cryptographic Verkle proof for course %s (student: %s)", courseID, studentDID)
	return proofData, nil
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

// VerifyCourseProof performs full cryptographic verification of a course proof against the Verkle root
func VerifyCourseProof(courseKey string, course CourseCompletion, proofData []byte, verkleRoot [32]byte) error {
	// Deserialize the proof bundle
	var proofBundle VerkleProofBundle
	if err := json.Unmarshal(proofData, &proofBundle); err != nil {
		return fmt.Errorf("failed to deserialize proof bundle: %w", err)
	}
	
	// Verify the course key matches
	if proofBundle.CourseKey != courseKey {
		return fmt.Errorf("proof bundle course key mismatch: expected %s, got %s", courseKey, proofBundle.CourseKey)
	}
	
	// Recreate the key hash exactly as it was during insertion
	courseKeyHash := sha256.Sum256([]byte(courseKey))
	var keyHash32 [32]byte
	copy(keyHash32[:], courseKeyHash[:])
	
	// Recreate the value hash from the course data
	courseData, err := json.Marshal(course)
	if err != nil {
		return fmt.Errorf("failed to serialize course data: %w", err)
	}
	courseValueHash := sha256.Sum256(courseData)
	
	// Perform cryptographic verification following Duc's approach
	// Check the StateDiff contains the expected key-value pair
	foundInDiff := false
	keyStem := keyHash32[:verkleLib.StemSize]
	keySuffix := keyHash32[verkleLib.StemSize]
	
	for _, stemDiff := range proofBundle.StateDiff {
		if bytes.Equal(keyStem, stemDiff.Stem[:]) {
			// Found the correct stem
			for _, suffixDiff := range stemDiff.SuffixDiffs {
				if keySuffix == suffixDiff.Suffix {
					// Found the exact key
					foundInDiff = true
					
					// Verify the value matches
					if suffixDiff.CurrentValue == nil {
						return fmt.Errorf("proof shows nil value for course %s, but course exists", course.CourseID)
					}
					
					// Compare the stored value with our computed value hash
					if !bytes.Equal((*suffixDiff.CurrentValue)[:], courseValueHash[:]) {
						return fmt.Errorf("value mismatch in proof for course %s", course.CourseID)
					}
					
					log.Printf("✅ Cryptographic verification successful for course %s", course.CourseID)
					break
				}
			}
			if foundInDiff {
				break
			}
		}
	}
	
	if !foundInDiff {
		return fmt.Errorf("course %s not found in proof's state diff", course.CourseID)
	}
	
	// Perform full cryptographic verification using go-verkle's Verify function
	// This verifies the IPA proof mathematically without needing the full tree!
	err = verkleLib.Verify(
		proofBundle.VerkleProof,
		verkleRoot[:],      // preStateRoot
		verkleRoot[:],      // postStateRoot (same as pre for proof of existence)
		proofBundle.StateDiff,
	)
	
	if err != nil {
		return fmt.Errorf("cryptographic verification failed for course %s: %w", course.CourseID, err)
	}
	
	log.Printf("✅ Full cryptographic IPA verification successful for course %s!", course.CourseID)
	return nil
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
	
	// 2. Verify each course proof against the Verkle root
	for _, course := range receipt.RevealedCourses {
		courseKey := fmt.Sprintf("%s:%s:%s", receipt.StudentDID, receipt.TermID, course.CourseID)
		proof, exists := receipt.CourseProofs[course.CourseID]
		if !exists {
			result.Warnings = append(result.Warnings, fmt.Sprintf("No proof found for course %s", course.CourseID))
			continue
		}
		
		// Verify the proof cryptographically
		if err := VerifyCourseProof(courseKey, course, proof, receipt.VerkleRoot); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Course %s proof verification failed: %v", course.CourseID, err))
		} else {
			log.Printf("✅ Course %s proof verified successfully", course.CourseID)
		}
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

// RebuildVerkleTree reconstructs the internal Verkle tree from saved course entries
// This is needed after deserializing from JSON since the tree field is not serialized
func (tvt *TermVerkleTree) RebuildVerkleTree() error {
	log.Printf("Rebuilding Verkle tree for term %s with %d course entries", tvt.TermID, len(tvt.CourseEntries))
	
	// Create new Verkle tree
	tvt.tree = verkleLib.New()
	
	// Re-insert all course entries
	for courseKey, course := range tvt.CourseEntries {
		courseKeyHash := sha256.Sum256([]byte(courseKey))
		
		// Serialize course data as value (same as original insertion)
		courseData, err := json.Marshal(course)
		if err != nil {
			return fmt.Errorf("failed to serialize course %s: %w", course.CourseID, err)
		}
		courseValueHash := sha256.Sum256(courseData)
		
		// Re-insert into Verkle tree
		err = tvt.tree.Insert(courseKeyHash[:], courseValueHash[:], nil)
		if err != nil {
			return fmt.Errorf("failed to re-insert course %s into verkle tree: %w", course.CourseID, err)
		}
	}
	
	log.Printf("✅ Verkle tree rebuilt successfully for term %s", tvt.TermID)
	return nil
}