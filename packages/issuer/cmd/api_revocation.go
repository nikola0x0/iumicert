package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"iumicert/issuer/database"
	blockchain_integration "iumicert/issuer/blockchain_integration"
	"iumicert/issuer/config"
)

// ===== REVOCATION API HANDLERS (Issuer Dashboard Only) =====

// validateCredentialExists checks if a student has a specific course in a term
func validateCredentialExists(studentID, termID, courseID string) error {
	// Load term tree (built Verkle tree with course entries)
	termFile := filepath.Join("data/verkle_trees", fmt.Sprintf("%s_verkle_tree.json", termID))

	// Check if term file exists
	if _, err := os.Stat(termFile); os.IsNotExist(err) {
		return fmt.Errorf("term %s not found", termID)
	}

	// Read term tree data
	data, err := os.ReadFile(termFile)
	if err != nil {
		return fmt.Errorf("failed to read term tree: %w", err)
	}

	// Parse term tree - CourseEntries has keys like "did:example:ITITIU00001:Semester_1_2023:IT013IU"
	var termTree struct {
		TermID        string                 `json:"term_id"`
		CourseEntries map[string]interface{} `json:"course_entries"`
	}

	if err := json.Unmarshal(data, &termTree); err != nil {
		return fmt.Errorf("failed to parse term tree: %w", err)
	}

	// Build expected course key
	studentDID := fmt.Sprintf("did:example:%s", studentID)
	courseKey := fmt.Sprintf("%s:%s:%s", studentDID, termID, courseID)

	// Check if credential exists in tree
	if _, exists := termTree.CourseEntries[courseKey]; !exists {
		return fmt.Errorf("credential not found: student %s does not have course %s in term %s",
			studentID, courseID, termID)
	}

	return nil
}

// handleCreateRevocationRequest creates a new revocation request (ADMIN ONLY)
// This is called after registrar validates student complaint through official channels
func handleCreateRevocationRequest(w http.ResponseWriter, r *http.Request) {
	var request struct {
		StudentID   string `json:"student_id"`
		TermID      string `json:"term_id"`
		CourseID    string `json:"course_id"`
		Reason      string `json:"reason"`
		RequestedBy string `json:"requested_by"` // Admin username
		Notes       string `json:"notes"`        // Additional context
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	// Validate required fields
	if request.StudentID == "" || request.TermID == "" || request.CourseID == "" || request.Reason == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Missing required fields: student_id, term_id, course_id, reason",
		})
		return
	}

	// VALIDATION 1: Check if credential actually exists in the term
	if err := validateCredentialExists(request.StudentID, request.TermID, request.CourseID); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Validation failed: %v", err),
		})
		return
	}

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database connection failed",
		})
		return
	}

	// VALIDATION 2: Check if this credential was already revoked
	var existingRevocation database.RevocationRequest
	err = db.Where("student_id = ? AND term_id = ? AND course_id = ? AND status IN (?)",
		request.StudentID, request.TermID, request.CourseID,
		[]string{"approved", "processed"}).First(&existingRevocation).Error

	if err == nil {
		// Found existing revocation
		respondJSON(w, http.StatusConflict, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("This credential already has a %s revocation request (ID: %s)",
				existingRevocation.Status, existingRevocation.RequestID),
		})
		return
	}

	// Create revocation request (already validated by registrar)
	revocationReq := &database.RevocationRequest{
		RequestID:   fmt.Sprintf("revoke_req_%s", uuid.New().String()),
		StudentID:   request.StudentID,
		TermID:      request.TermID,
		CourseID:    request.CourseID,
		Reason:      request.Reason,
		RequestedBy: request.RequestedBy,
		Status:      "approved", // Directly approved since registrar already validated
		ApprovedBy:  request.RequestedBy,
		ApprovedAt:  timePtr(time.Now()),
		Notes:       request.Notes,
	}

	if err := database.CreateRevocationRequest(db, revocationReq); err != nil {
		log.Printf("❌ Failed to create revocation request: %v", err)
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to create revocation request",
		})
		return
	}

	log.Printf("✅ Revocation request created (approved): %s for %s/%s/%s",
		revocationReq.RequestID, request.StudentID, request.TermID, request.CourseID)

	respondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message":    "Revocation request created and approved. Will be processed during next term publication.",
			"request_id": revocationReq.RequestID,
			"status":     "approved",
			"note":       "This will be automatically processed when the next term is published.",
		},
	})
}

// handleListRevocationRequests lists all revocation requests with optional filters
func handleListRevocationRequests(w http.ResponseWriter, r *http.Request) {
	termID := r.URL.Query().Get("term_id")
	status := r.URL.Query().Get("status")

	db, err := database.Connect()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database connection failed",
		})
		return
	}

	requests, err := database.GetAllRevocationRequests(db, termID, status)
	if err != nil {
		log.Printf("❌ Failed to get revocation requests: %v", err)
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to retrieve revocation requests",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"requests": requests,
			"count":    len(requests),
		},
	})
}

// handleGetPendingRevocations gets approved (not yet processed) revocations for a term
func handleGetPendingRevocations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	termID := vars["term_id"]

	db, err := database.Connect()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database connection failed",
		})
		return
	}

	// Get approved but not processed revocations
	var requests []database.RevocationRequest
	err = db.Where("term_id = ? AND status = ?", termID, "approved").Find(&requests).Error
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to retrieve pending revocations",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"term_id":  termID,
			"requests": requests,
			"count":    len(requests),
			"note":     "These will be automatically processed during next term publication",
		},
	})
}

// handleGetTermVersionHistory gets all versions for a term
func handleGetTermVersionHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	termID := vars["term_id"]

	// Get from database
	db, err := database.Connect()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database connection failed",
		})
		return
	}

	versions, err := database.GetTermVersionHistory(db, termID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to retrieve version history",
		})
		return
	}

	// Also get from blockchain for comparison
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("⚠️  Failed to load config for blockchain check: %v", err)
	}

	var blockchainData interface{}
	if err == nil {
		blockchainIntegration, err := blockchain_integration.NewBlockchainIntegration(
			cfg.Network,
			cfg.IssuerPrivateKey,
			cfg.ContractAddress,
		)
		if err == nil {
			defer blockchainIntegration.Close()

			ctx := r.Context()
			blockchainVersions, blockchainRoots, err := blockchainIntegration.GetTermHistory(ctx, termID)

			if err == nil {
				blockchainData = map[string]interface{}{
					"versions": blockchainVersions,
					"roots":    blockchainRoots,
				}
			}
		}
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"term_id":    termID,
			"versions":   versions,
			"count":      len(versions),
			"blockchain": blockchainData,
		},
	})
}

// handleGetRevocationStats gets revocation statistics
func handleGetRevocationStats(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database connection failed",
		})
		return
	}

	stats, err := database.GetRevocationStats(db)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to retrieve revocation statistics",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    stats,
	})
}

// handleDeleteRevocationRequest allows admin to delete a revocation request
func handleDeleteRevocationRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["request_id"]

	db, err := database.Connect()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database connection failed",
		})
		return
	}

	err = db.Where("request_id = ?", requestID).Delete(&database.RevocationRequest{}).Error
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to delete revocation request",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message": "Revocation request deleted",
		},
	})
}

// handleProcessRevocations processes all approved revocations
// This is called automatically after any term is published via the dashboard
func handleProcessRevocations(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔄 API: Processing approved revocations...")

	// Load configuration to get network settings and private key
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("❌ Failed to load config: %v", err)
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to load configuration",
		})
		return
	}

	// Check if private key is configured
	if cfg.IssuerPrivateKey == "" {
		log.Printf("⚠️  No issuer private key configured, skipping revocation processing")
		respondJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Data: map[string]interface{}{
				"message":   "Skipped - no private key configured",
				"processed": 0,
			},
		})
		return
	}

	// Get all approved revocations
	db, err := database.Connect()
	if err != nil {
		log.Printf("❌ Database connection failed: %v", err)
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database connection failed",
		})
		return
	}

	var approvedRevocations []database.RevocationRequest
	err = db.Where("status = ?", "approved").Find(&approvedRevocations).Error
	if err != nil {
		log.Printf("❌ Failed to get approved revocations: %v", err)
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to get approved revocations",
		})
		return
	}

	if len(approvedRevocations) == 0 {
		log.Printf("✅ No pending revocations to process")
		respondJSON(w, http.StatusOK, APIResponse{
			Success: true,
			Data: map[string]interface{}{
				"message":   "No pending revocations to process",
				"processed": 0,
			},
		})
		return
	}

	// Group revocations by term
	revocationsByTerm := make(map[string][]database.RevocationRequest)
	for _, rev := range approvedRevocations {
		revocationsByTerm[rev.TermID] = append(revocationsByTerm[rev.TermID], rev)
	}

	log.Printf("📋 Found %d approved revocations across %d terms",
		len(approvedRevocations), len(revocationsByTerm))

	// Process each term with revocations
	processedCount := 0
	var errors []string

	for termID, revocations := range revocationsByTerm {
		log.Printf("🔄 Processing %d revocations for term: %s", len(revocations), termID)

		// Execute revocation by rebuilding tree and publishing new version
		err := supersedeTermWithRevocations(termID, revocations, cfg, db)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to process revocations for term %s: %v", termID, err)
			log.Printf("❌ %s", errMsg)
			errors = append(errors, errMsg)
			continue
		}

		processedCount += len(revocations)
		log.Printf("✅ Successfully processed revocations for term %s", termID)
	}

	// Return result
	result := map[string]interface{}{
		"message":      fmt.Sprintf("Processed %d revocations across %d terms", processedCount, len(revocationsByTerm)),
		"processed":    processedCount,
		"terms_affected": len(revocationsByTerm),
	}

	if len(errors) > 0 {
		result["errors"] = errors
		result["partial_success"] = true
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    result,
	})
}

// Helper function
func timePtr(t time.Time) *time.Time {
	return &t
}
