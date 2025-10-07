package database

import (
	"time"

	"gorm.io/datatypes"
)

// Student represents a student in the system
type Student struct {
	ID                 uint      `gorm:"primaryKey"`
	StudentID          string    `gorm:"uniqueIndex;not null;size:50"` // ITITIU00001
	Name               string    `gorm:"size:255"`
	Email              string    `gorm:"size:255"`
	DID                string    `gorm:"index;size:255"` // Decentralized identifier
	EnrollmentDate     time.Time
	ExpectedGraduation time.Time
	Status             string    `gorm:"index;size:50;default:'active'"` // "active", "graduated", "withdrawn"
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// Term represents an academic term
type Term struct {
	ID               uint      `gorm:"primaryKey"`
	TermID           string    `gorm:"uniqueIndex;not null;size:50"` // Semester_1_2023
	StartDate        time.Time
	EndDate          time.Time
	VerkleRootHex    string    `gorm:"index;size:64"` // Hex string for indexing
	VerkleRootBytes  []byte    `gorm:"type:bytea"`    // Binary for verification
	BlockchainTxHash string    `gorm:"index;size:66"` // Ethereum tx hash
	BlockNumber      uint64
	PublishedAt      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// TermReceipt represents a single term's receipt for a student
type TermReceipt struct {
	ID          uint   `gorm:"primaryKey"`
	ReceiptID   string `gorm:"uniqueIndex;not null;size:255"` // receipt_ITITIU00001_Semester_1_2023_20251006
	StudentID   string `gorm:"uniqueIndex:idx_student_term_unique;not null;size:50"`
	TermID      string `gorm:"uniqueIndex:idx_student_term_unique;not null;size:50"`

	// Proof Data (JSONB for efficient storage and querying)
	VerkleProof     datatypes.JSON `gorm:"type:jsonb;not null"` // Full VerkleProof structure
	StateDiff       datatypes.JSON `gorm:"type:jsonb;not null"` // StateDiff array
	RevealedCourses datatypes.JSON `gorm:"type:jsonb;not null"` // Array of course completions

	// Metadata
	CourseCount   int
	VerkleRootHex string    `gorm:"index;size:64"` // For quick blockchain verification
	GeneratedAt   time.Time `gorm:"index"`
	IsSelective   bool      `gorm:"default:false"` // true if not all courses revealed

	CreatedAt time.Time
	UpdatedAt time.Time
}

// AccumulatedReceipt represents a diploma or progress receipt (multiple terms)
type AccumulatedReceipt struct {
	ID                   uint   `gorm:"primaryKey"`
	AccumulatedReceiptID string `gorm:"uniqueIndex;not null;size:255"` // diploma_ITITIU00001_20251206
	StudentID            string `gorm:"index:idx_student_type;not null;size:50"`

	// Receipt Type
	Type string `gorm:"index:idx_student_type;size:50"` // "progress", "diploma", "custom"

	// Accumulated Data
	TermReceiptIDs datatypes.JSON `gorm:"type:jsonb"` // Array of term receipt IDs included
	TermsIncluded  datatypes.JSON `gorm:"type:jsonb"` // Array of term IDs (e.g., ["Semester_1_2023", ...])
	AllCourses     datatypes.JSON `gorm:"type:jsonb"` // All courses from all terms

	// Aggregated Proofs (optional - for batch verification)
	AggregatedProofData datatypes.JSON `gorm:"type:jsonb"` // Optimized combined proof structure

	// Summary Statistics
	TotalCourses   int
	TotalCredits   int
	GPA            float64
	CompletedTerms int

	// Metadata
	GeneratedAt time.Time `gorm:"index"`
	ValidFrom   time.Time // Start of first term
	ValidUntil  *time.Time // End of last term (null for progress receipts)

	// Blockchain Anchoring
	BlockchainVerified bool   `gorm:"default:false"`
	BlockchainTxHash   string `gorm:"index;size:66"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// VerificationLog represents a verification attempt
type VerificationLog struct {
	ID               uint   `gorm:"primaryKey"`
	ReceiptID        string `gorm:"index;not null;size:255"` // Can be TermReceipt or AccumulatedReceipt
	ReceiptType      string `gorm:"size:50"`                 // "term", "accumulated"
	VerifierID       string `gorm:"size:255"`                // Who verified (employer, university, etc.)
	VerificationMode string `gorm:"size:50"`                 // "local", "blockchain", "full_ipa"
	Success          bool
	ErrorMessage     string `gorm:"type:text"`
	VerifiedAt       time.Time `gorm:"index"`
	IPAddress        string    `gorm:"size:45"` // IPv6 compatible
	UserAgent        string    `gorm:"type:text"`

	CreatedAt time.Time
}

// BlockchainTransaction represents a blockchain transaction for publishing roots
type BlockchainTransaction struct {
	ID          uint   `gorm:"primaryKey"`
	TxHash      string `gorm:"uniqueIndex;not null;size:66"`
	TermID      string `gorm:"index;size:50"`
	VerkleRoot  []byte `gorm:"type:bytea"`
	BlockNumber uint64 `gorm:"index"`
	GasUsed     uint64
	Status      string    `gorm:"size:50"` // "pending", "confirmed", "failed"
	SubmittedAt time.Time
	ConfirmedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
