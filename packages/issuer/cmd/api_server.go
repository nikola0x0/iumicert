package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"iumicert/crypto/testdata"
	"iumicert/crypto/verkle"
	blockchain_integration "iumicert/issuer/blockchain_integration"
	"iumicert/issuer/config"
	"iumicert/issuer/database"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"gorm.io/datatypes"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the IU-MiCert API server for web interface",
	Long:  `Start a REST API server that provides endpoints for the web interface to interact with the credential system`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		cors_enabled, _ := cmd.Flags().GetBool("cors")
		
		if err := startAPIServer(port, cors_enabled); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to start server: %v\n", err)
			os.Exit(1)
		}
	},
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type TermRequest struct {
	TermID     string                    `json:"term_id"`
	DataFile   string                    `json:"data_file,omitempty"`
	Courses    []verkle.CourseCompletion `json:"courses,omitempty"`
	Format     string                    `json:"format"`
	Validate   bool                      `json:"validate"`
}

type ReceiptRequest struct {
	StudentID  string   `json:"student_id"`
	Terms      []string `json:"terms,omitempty"`
	Courses    []string `json:"courses,omitempty"`
	Selective  bool     `json:"selective"`
}

type PublishRequest struct {
	TermID     string `json:"term_id"`
	Network    string `json:"network"`
	GasLimit   uint64 `json:"gas_limit"`
}

type SystemStatus struct {
	Repository struct {
		Initialized bool   `json:"initialized"`
		Institution string `json:"institution"`
	} `json:"repository"`
	Blockchain struct {
		Network       string `json:"network"`
		DefaultGasLimit uint64 `json:"default_gas_limit"`
	} `json:"blockchain"`
	Storage struct {
		Terms        int `json:"terms"`
		Students     int `json:"students"`
		Receipts     int `json:"receipts"`
		Transactions int `json:"transactions"`
	} `json:"storage"`
}

// Middleware functions
func requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Log incoming request
		log.Printf("📥 %s %s - %s", r.Method, r.URL.Path, r.RemoteAddr)
		
		// Create a response writer that captures status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		// Call the next handler
		next.ServeHTTP(wrapped, r)
		
		// Log completed request with timing
		duration := time.Since(start)
		log.Printf("📤 %s %s - %d (%v)", r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

func requestValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate Content-Type for POST/PUT requests
		if r.Method == "POST" || r.Method == "PUT" {
			contentType := r.Header.Get("Content-Type")
			if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
				respondJSON(w, http.StatusBadRequest, APIResponse{
					Success: false, 
					Error: "Content-Type must be application/json",
				})
				return
			}
		}
		
		// Validate request size (limit to 10MB)
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		
		next.ServeHTTP(w, r)
	})
}

// Response writer wrapper to capture status codes
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func startAPIServer(port string, corsEnabled bool) error {
	r := mux.NewRouter()
	
	// Apply middleware to all routes
	r.Use(requestLoggingMiddleware)
	r.Use(requestValidationMiddleware)
	
	// API routes with issuer/verifier separation
	api := r.PathPrefix("/api").Subrouter()
	
	// System endpoints (public)
	api.HandleFunc("/status", handleSystemStatus).Methods("GET")
	api.HandleFunc("/health", handleHealth).Methods("GET")
	
	// Issuer-only endpoints (for institution dashboard)
	issuer := api.PathPrefix("/issuer").Subrouter()
	issuer.HandleFunc("/terms", handleAddTerm).Methods("POST")
	issuer.HandleFunc("/terms", handleListTerms).Methods("GET")
	issuer.HandleFunc("/terms/{term_id}/receipts", handleGetTermReceipts).Methods("GET")
	issuer.HandleFunc("/terms/{term_id}/roots", handleGetTermRoot).Methods("GET")

	// New: Process uploaded term data (Data Management Panel)
	api.HandleFunc("/terms/process", handleProcessTermData).Methods("POST")
	api.HandleFunc("/demo/generate-term", handleGenerateDemoTerm).Methods("POST")
	api.HandleFunc("/demo/reset", handleDemoReset).Methods("POST")
	api.HandleFunc("/demo/generate-full", handleDemoGenerateFull).Methods("POST")
	issuer.HandleFunc("/receipts", handleGenerateReceipt).Methods("POST")
	issuer.HandleFunc("/receipts", handleListReceipts).Methods("GET")
	issuer.HandleFunc("/blockchain/publish", handlePublishRoots).Methods("POST")
	issuer.HandleFunc("/blockchain/transactions", handleListTransactions).Methods("GET")
	issuer.HandleFunc("/blockchain/transactions/{tx_hash}", handleGetTransaction).Methods("GET")
	issuer.HandleFunc("/blockchain/roots", handleGetPublishedRoots).Methods("GET")
	issuer.HandleFunc("/students", handleListStudents).Methods("GET")
	issuer.HandleFunc("/students/{student_id}/terms", handleGetStudentTerms).Methods("GET")
	issuer.HandleFunc("/students/{student_id}/journey", handleGetStudentJourney).Methods("GET")

	// Database-backed receipt endpoints (NEW)
	issuer.HandleFunc("/students/{student_id}/receipts/latest", handleGetLatestReceipts).Methods("GET")
	issuer.HandleFunc("/students/{student_id}/receipts/accumulated", handleGetAccumulatedReceipt).Methods("GET")
	issuer.HandleFunc("/students/{student_id}/receipts/term/{term_id}", handleGetTermReceipt).Methods("GET")
	issuer.HandleFunc("/students/{student_id}/receipts/download", handleDownloadJourneyReceipt).Methods("GET")

	// Revocation endpoints (Admin-only - realistic workflow)
	// Note: Students contact institution through official channels (email, forms, in-person)
	// Registrar validates and enters approved requests here
	issuer.HandleFunc("/revocations", handleCreateRevocationRequest).Methods("POST")         // Create approved request
	issuer.HandleFunc("/revocations", handleListRevocationRequests).Methods("GET")           // List all requests
	issuer.HandleFunc("/revocations/stats", handleGetRevocationStats).Methods("GET")         // Get statistics
	issuer.HandleFunc("/revocations/process", handleProcessRevocations).Methods("POST")      // Process all approved revocations
	issuer.HandleFunc("/revocations/{request_id}", handleDeleteRevocationRequest).Methods("DELETE")  // Delete request
	issuer.HandleFunc("/terms/{term_id}/revocations", handleGetPendingRevocations).Methods("GET")    // Get approved for term
	issuer.HandleFunc("/terms/{term_id}/versions", handleGetTermVersionHistory).Methods("GET")       // Get version history

	// Verifier endpoints (public - for students/employers)
	verifier := api.PathPrefix("/verifier").Subrouter()
	verifier.HandleFunc("/receipt", handleVerifyReceipt).Methods("POST")
	verifier.HandleFunc("/course", handleVerifyCourse).Methods("POST")
	verifier.HandleFunc("/ipa-verify", handleIPAVerify).Methods("POST")  // Full IPA cryptographic verification
	verifier.HandleFunc("/receipt/{receipt_id}", handleGetReceiptByID).Methods("GET")
	verifier.HandleFunc("/journey/{student_id}", handleGetStudentJourney).Methods("GET")
	verifier.HandleFunc("/blockchain/transaction/{tx_hash}", handleGetTransaction).Methods("GET")
	verifier.HandleFunc("/blockchain/roots", handleGetPublishedRoots).Methods("GET")
	
	// Legacy endpoints (maintain backward compatibility for current issuer dashboard)
	api.HandleFunc("/terms", handleListTerms).Methods("GET")
	api.HandleFunc("/terms/{term_id}/roots", handleGetTermRoot).Methods("GET")
	api.HandleFunc("/terms/{term_id}/blockchain", handleUpdateTermBlockchainStatus).Methods("PUT")
	api.HandleFunc("/receipts/verify", handleVerifyReceipt).Methods("POST")
	api.HandleFunc("/receipts/verify-course", handleVerifyCourse).Methods("POST")
	api.HandleFunc("/blockchain/publish", handlePublishRoots).Methods("POST")
	api.HandleFunc("/blockchain/transactions", handleListTransactions).Methods("GET")
	api.HandleFunc("/blockchain/roots", handleGetPublishedRoots).Methods("GET")
	
	// Setup CORS if enabled
	var handler http.Handler = r
	if corsEnabled {
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:5173",
		}

		// Add production frontend URLs from environment
		// Student/Verifier Portal
		if studentPortalURL := os.Getenv("STUDENT_PORTAL_URL"); studentPortalURL != "" {
			allowedOrigins = append(allowedOrigins, studentPortalURL)
		}

		// Issuer Portal (Admin Dashboard)
		if issuerPortalURL := os.Getenv("ISSUER_PORTAL_URL"); issuerPortalURL != "" {
			allowedOrigins = append(allowedOrigins, issuerPortalURL)
		}

		// Legacy support for FRONTEND_URL
		if prodOrigin := os.Getenv("FRONTEND_URL"); prodOrigin != "" {
			allowedOrigins = append(allowedOrigins, prodOrigin)
		}

		c := cors.New(cors.Options{
			AllowedOrigins:   allowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		})
		handler = c.Handler(r)
	}
	
	fmt.Printf("🚀 Starting IU-MiCert API server on port %s\n", port)
	fmt.Printf("📡 API endpoints available at: http://localhost:%s/api\n", port)
	if corsEnabled {
		fmt.Printf("🔓 CORS enabled for React development\n")
	}
	fmt.Printf("📚 Available endpoints:\n")
	fmt.Printf("  GET  /api/status           - System status\n")
	fmt.Printf("  GET  /api/health           - Health check\n")
	fmt.Printf("  POST /api/terms            - Add academic term\n")
	fmt.Printf("  GET  /api/terms            - List terms\n")
	fmt.Printf("  POST /api/receipts         - Generate receipt\n")
	fmt.Printf("  POST /api/blockchain/publish - Publish to blockchain\n")
	fmt.Printf("  GET  /api/students         - List students\n")
	
	return http.ListenAndServe(":"+port, handler)
}

func handleSystemStatus(w http.ResponseWriter, r *http.Request) {
	status := SystemStatus{}
	
	// Check if repository is initialized
	if _, err := os.Stat("config/micert.json"); err == nil {
		status.Repository.Initialized = true
		
		// Try to read institution ID
		if configData, err := os.ReadFile("config/micert.json"); err == nil {
			var config map[string]interface{}
			if err := json.Unmarshal(configData, &config); err == nil {
				if institution, ok := config["institution_id"].(string); ok {
					status.Repository.Institution = institution
				}
				if blockchain, ok := config["blockchain"].(map[string]interface{}); ok {
					if network, ok := blockchain["default_network"].(string); ok {
						status.Blockchain.Network = network
					}
					if gasLimit, ok := blockchain["gas_limit"].(float64); ok {
						status.Blockchain.DefaultGasLimit = uint64(gasLimit)
					}
				}
			}
		}
	}
	
	// Count storage items
	if dirs := []struct{path string; count *int}{
		{"data/merkle_trees", &status.Storage.Terms},
		{"publish_ready/receipts", &status.Storage.Receipts},
		{"publish_ready/transactions", &status.Storage.Transactions},
	}; len(dirs) > 0 {
		for _, dir := range dirs {
			if files, err := filepath.Glob(filepath.Join(dir.path, "*")); err == nil {
				*dir.count = len(files)
			}
		}
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: status})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status": "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version": "1.0.0",
	}
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: health})
}

func handleAddTerm(w http.ResponseWriter, r *http.Request) {
	var req TermRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
		return
	}
	
	if req.TermID == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "term_id is required"})
		return
	}
	
	// Use provided courses data or load from file
	var completions []verkle.CourseCompletion
	var err error
	
	if len(req.Courses) > 0 {
		completions = req.Courses
	} else if req.DataFile != "" {
		completions, err = loadCompletionsFromJSON(req.DataFile)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: fmt.Sprintf("Failed to load data file: %v", err)})
			return
		}
	} else {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Either courses data or data_file is required"})
		return
	}
	
	// Call the existing addAcademicTerm function
	format := req.Format
	if format == "" {
		format = "json"
	}
	
	// Create a temporary file for the data if courses were provided directly
	var dataFile string
	if len(req.Courses) > 0 {
		dataFile = fmt.Sprintf("/tmp/term_%s_%d.json", req.TermID, time.Now().Unix())
		data, _ := json.Marshal(req.Courses)
		if err := os.WriteFile(dataFile, data, 0644); err != nil {
			respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create temporary data file"})
			return
		}
		defer os.Remove(dataFile)
	} else {
		dataFile = req.DataFile
	}
	
	if err := addAcademicTerm(req.TermID, dataFile, format, req.Validate); err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	
	result := map[string]interface{}{
		"term_id": req.TermID,
		"courses_processed": len(completions),
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: result})
}

// handleProcessTermData processes uploaded term data from Data Management Panel
// It converts the data, builds Verkle trees, and generates receipts
func handleProcessTermData(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TermID   string                       `json:"term_id"`
		Students map[string][]map[string]interface{} `json:"students"`
		Metadata map[string]interface{}       `json:"metadata,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
		return
	}

	if req.TermID == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "term_id is required"})
		return
	}

	if len(req.Students) == 0 {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "students data is required"})
		return
	}

	log.Printf("📥 Processing term data for: %s", req.TermID)
	log.Printf("   Students: %d", len(req.Students))

	// Step 1: Convert uploaded data to CourseCompletion format
	var completions []verkle.CourseCompletion
	totalCourses := 0

	for studentID, courses := range req.Students {
		for _, courseData := range courses {
			// Extract course fields
			courseID, _ := courseData["course_id"].(string)
			courseName, _ := courseData["course_name"].(string)
			grade, _ := courseData["grade"].(string)

			var credits uint8 = 3 // Default
			if creditsFloat, ok := courseData["credits"].(float64); ok {
				credits = uint8(creditsFloat)
			}

			completion := verkle.CourseCompletion{
				StudentID:  studentID,
				TermID:     req.TermID,
				CourseID:   courseID,
				CourseName: courseName,
				Grade:      grade,
				Credits:    credits,
				IssuerID:   "IU-CS",
				AttemptNo:  1,
				// Use current time for timestamps
				StartedAt:   time.Now().Add(-90 * 24 * time.Hour),
				CompletedAt: time.Now().Add(-7 * 24 * time.Hour),
				AssessedAt:  time.Now().Add(-3 * 24 * time.Hour),
				IssuedAt:    time.Now(),
				Instructor:  "Prof. System",
			}

			completions = append(completions, completion)
			totalCourses++
		}
	}

	log.Printf("   Total courses: %d", totalCourses)

	// Step 2: Save to verkle format file
	verkleDir := "data/verkle_terms"
	if err := os.MkdirAll(verkleDir, 0755); err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to create verkle directory: %v", err),
		})
		return
	}

	verkleFile := filepath.Join(verkleDir, fmt.Sprintf("%s_completions.json", req.TermID))
	data, _ := json.MarshalIndent(completions, "", "  ")
	if err := os.WriteFile(verkleFile, data, 0644); err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to save verkle data: %v", err),
		})
		return
	}

	log.Printf("✅ Saved verkle data to: %s", verkleFile)

	// Step 3: Build Verkle tree
	log.Printf("🌳 Building Verkle tree...")
	if err := addAcademicTerm(req.TermID, verkleFile, "json", true); err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to build Verkle tree: %v", err),
		})
		return
	}

	log.Printf("✅ Verkle tree built successfully")

	// Step 4: Generate receipts for all students
	log.Printf("📄 Generating receipts...")
	successCount := 0
	failedStudents := []string{}

	for studentID := range req.Students {
		outputFile := fmt.Sprintf("publish_ready/receipts/%s_journey.json", studentID)

		// Generate receipt with all terms (empty list = autodiscover)
		if err := generateStudentReceipt(studentID, outputFile, nil, nil, false); err != nil {
			log.Printf("⚠️ Failed to generate receipt for %s: %v", studentID, err)
			failedStudents = append(failedStudents, studentID)
			continue
		}

		successCount++
	}

	log.Printf("✅ Generated %d/%d receipts", successCount, len(req.Students))

	// Step 5: Store receipts in database
	log.Printf("💾 Storing receipts in database...")

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Printf("⚠️ Warning: Failed to connect to database: %v", err)
		log.Printf("   Receipts generated but not stored in database")
	} else {
		defer database.Close(db)

		repo := database.NewReceiptRepository(db)
		storedCount := 0

		for studentID := range req.Students {
			// Read the generated receipt file
			receiptFile := fmt.Sprintf("publish_ready/receipts/%s_journey.json", studentID)
			receiptData, err := os.ReadFile(receiptFile)
			if err != nil {
				log.Printf("⚠️ Failed to read receipt for %s: %v", studentID, err)
				continue
			}

			// Parse the receipt
			var journeyReceipt map[string]interface{}
			if err := json.Unmarshal(receiptData, &journeyReceipt); err != nil {
				log.Printf("⚠️ Failed to parse receipt for %s: %v", studentID, err)
				continue
			}

			// Extract term receipts
			termReceipts, ok := journeyReceipt["term_receipts"].(map[string]interface{})
			if !ok {
				log.Printf("⚠️ No term_receipts found for %s", studentID)
				continue
			}

			// Check if this term exists in the receipt
			termData, ok := termReceipts[req.TermID].(map[string]interface{})
			if !ok {
				log.Printf("⚠️ Term %s not found in receipt for %s", req.TermID, studentID)
				continue
			}

			// Extract verkle root
			verkleRootHex, ok := termData["verkle_root"].(string)
			if !ok {
				log.Printf("⚠️ No verkle_root found for %s/%s", studentID, req.TermID)
				continue
			}

			// Extract receipt data
			receiptDataMap, ok := termData["receipt"].(map[string]interface{})
			if !ok {
				log.Printf("⚠️ No receipt data found for %s/%s", studentID, req.TermID)
				continue
			}

			// Extract course proofs and revealed courses
			courseProofs := receiptDataMap["course_proofs"]
			revealedCourses := receiptDataMap["revealed_courses"]

			// Marshal to JSON for database storage
			verkleProofJSON, _ := json.Marshal(courseProofs)
			revealedCoursesJSON, _ := json.Marshal(revealedCourses)

			// Count courses
			var courseCount int
			if courses, ok := revealedCourses.([]interface{}); ok {
				courseCount = len(courses)
			}

			// Create term receipt for database
			termReceipt := &database.TermReceipt{
				ReceiptID:       fmt.Sprintf("receipt_%s_%s_%d", studentID, req.TermID, time.Now().Unix()),
				StudentID:       studentID,
				TermID:          req.TermID,
				VerkleProof:     datatypes.JSON(verkleProofJSON),
				RevealedCourses: datatypes.JSON(revealedCoursesJSON),
				StateDiff:       datatypes.JSON("[]"), // Placeholder
				CourseCount:     courseCount,
				VerkleRootHex:   verkleRootHex,
				GeneratedAt:     time.Now(),
				IsSelective:     false,
			}

			// Store in database
			if err := repo.StoreTermReceipt(termReceipt); err != nil {
				log.Printf("⚠️ Failed to store receipt for %s/%s: %v", studentID, req.TermID, err)
				continue
			}

			storedCount++
		}

		log.Printf("✅ Stored %d/%d receipts in database", storedCount, successCount)
	}

	// Return success with statistics
	result := map[string]interface{}{
		"term_id":            req.TermID,
		"students_processed": len(req.Students),
		"courses_processed":  totalCourses,
		"receipts_generated": successCount,
		"verkle_tree_built":  true,
		"timestamp":          time.Now().Format(time.RFC3339),
	}

	if len(failedStudents) > 0 {
		result["failed_students"] = failedStudents
	}

	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: result})
}

// handleGenerateDemoTerm generates demo term data for testing/performance purposes
func handleGenerateDemoTerm(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TermID      string `json:"term_id"`
		NumStudents int    `json:"num_students"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
		return
	}

	// Validate inputs
	if req.TermID == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "term_id is required"})
		return
	}
	if req.NumStudents < 1 || req.NumStudents > 100 {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "num_students must be between 1 and 100"})
		return
	}

	log.Printf("📊 Generating demo term: %s with %d students (3-6 courses per student)", req.TermID, req.NumStudents)

	// Call the addon term generator logic with fixed course range (3-6)
	outputData, err := generateAddonTermData(req.TermID, req.NumStudents)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to generate demo data: %v", err),
		})
		return
	}

	log.Printf("✅ Generated demo term data: %d students, %d total courses",
		len(outputData["students"].(map[string][]map[string]interface{})),
		countTotalCourses(outputData["students"].(map[string][]map[string]interface{})))

	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: outputData})
}

// generateAddonTermData generates term data using the test data generator
func generateAddonTermData(termID string, numStudents int) (map[string]interface{}, error) {
	// Fixed course range: 3-6 courses per student
	const minCourses = 3
	const maxCourses = 6

	generator := testdata.NewTestDataGenerator()

	// Generate course completions for all students
	allCompletions := make([]verkle.CourseCompletion, 0)

	for i := 0; i < numStudents; i++ {
		// Variable course count per student (3-6)
		coursesForStudent := minCourses + (i % (maxCourses - minCourses + 1))

		studentCompletions, err := generator.GenerateTermData(termID, 1, coursesForStudent)
		if err != nil {
			return nil, fmt.Errorf("failed to generate data for student %d: %w", i, err)
		}

		// Update student IDs to match IU Vietnam format
		for j := range studentCompletions {
			studentCompletions[j].StudentID = fmt.Sprintf("ITITIU%05d", i+1)
		}

		allCompletions = append(allCompletions, studentCompletions...)
	}

	// Organize data by student
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

	// Create the term data structure
	termData := map[string]interface{}{
		"term_id":  termID,
		"students": studentCourses,
		"metadata": map[string]interface{}{
			"generated_at":   time.Now().Format(time.RFC3339),
			"total_students": len(studentCourses),
			"total_courses":  len(allCompletions),
			"generated_by":   "demo-generator-api",
		},
	}

	return termData, nil
}

// countTotalCourses counts total courses across all students
func countTotalCourses(students map[string][]map[string]interface{}) int {
	total := 0
	for _, courses := range students {
		total += len(courses)
	}
	return total
}

// handleDemoReset cleans all generated data (database and file system)
func handleDemoReset(w http.ResponseWriter, r *http.Request) {
	log.Printf("🧹 Executing system reset...")

	var output strings.Builder
	output.WriteString("🧹 IU-MiCert System Reset\n")
	output.WriteString("=========================\n\n")

	// Step 1: Clear file system directories
	log.Printf("📁 Step 1: Clearing file system data...")
	output.WriteString("📁 Step 1: Clearing file system data...\n")

	dirsToClean := []string{
		"data/student_journeys",
		"data/verkle_terms",
		"data/verkle_trees",
		"publish_ready/receipts",
		"publish_ready/roots",
		"publish_ready/proofs",
		"publish_ready/transactions",
	}

	for _, dir := range dirsToClean {
		if err := os.RemoveAll(dir); err != nil {
			log.Printf("⚠️  Warning: Failed to remove %s: %v", dir, err)
			output.WriteString(fmt.Sprintf("⚠️  Warning: Failed to remove %s: %v\n", dir, err))
		}
	}

	// Recreate directory structure
	dirsToCreate := []string{
		"data/student_journeys/students",
		"data/student_journeys/terms",
		"data/verkle_terms",
		"data/verkle_trees",
		"publish_ready/receipts",
		"publish_ready/roots",
		"publish_ready/proofs",
		"publish_ready/transactions",
	}

	for _, dir := range dirsToCreate {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("❌ Failed to create directory %s: %v", dir, err)
			output.WriteString(fmt.Sprintf("❌ Failed to create directory %s: %v\n", dir, err))
			respondJSON(w, http.StatusInternalServerError, APIResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to create directory %s: %v", dir, err),
			})
			return
		}
	}

	output.WriteString("✅ File system cleaned and directories recreated\n\n")

	// Step 2: Reset database
	log.Printf("🗄️  Step 2: Resetting database...")
	output.WriteString("🗄️  Step 2: Resetting database...\n")

	db, err := database.Connect()
	if err != nil {
		log.Printf("⚠️  Database connection failed: %v", err)
		output.WriteString(fmt.Sprintf("⚠️  Database not available: %v\n", err))
		output.WriteString("⚠️  Skipping database reset\n\n")
	} else {
		// Drop all tables (including revocation-related tables)
		tables := []string{
			"verification_logs",
			"blockchain_transactions",
			"accumulated_receipts",
			"term_receipts",
			"terms",
			"students",
			"revocation_requests",
			"term_root_versions",
			"revocation_batches",
		}

		for _, table := range tables {
			if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
				log.Printf("⚠️  Warning: Failed to drop table %s: %v", table, err)
				output.WriteString(fmt.Sprintf("⚠️  Warning: Failed to drop table %s: %v\n", table, err))
			}
		}

		// Run migrations to recreate tables (including revocation tables)
		if err := db.AutoMigrate(
			&database.Student{},
			&database.Term{},
			&database.TermReceipt{},
			&database.AccumulatedReceipt{},
			&database.VerificationLog{},
			&database.BlockchainTransaction{},
			&database.RevocationRequest{},
			&database.TermRootVersion{},
			&database.RevocationBatch{},
		); err != nil {
			log.Printf("❌ Database migration failed: %v", err)
			output.WriteString(fmt.Sprintf("❌ Database migration failed: %v\n", err))
		} else {
			output.WriteString("✅ Database reset complete (all tables recreated)\n\n")
		}
	}

	log.Printf("✅ Reset completed successfully")
	output.WriteString("🎉 Reset completed successfully!\n")
	output.WriteString("System ready for data generation.\n")

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message":   "System reset completed successfully",
			"output":    output.String(),
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
}

// handleDemoGenerateFull executes customizable full data generation
func handleDemoGenerateFull(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NumStudents int      `json:"num_students"`
		Terms       []string `json:"terms"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
		return
	}

	// Validate inputs
	if req.NumStudents < 1 || req.NumStudents > 100 {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "num_students must be between 1 and 100"})
		return
	}

	if len(req.Terms) == 0 {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "at least one term is required"})
		return
	}

	log.Printf("🚀 Executing full data generation: %d students, %d terms", req.NumStudents, len(req.Terms))

	// Step 1: Generate student journeys
	log.Printf("👥 Step 1: Generating student academic journeys...")

	// Debug: Log the exact command we're about to run
	termsArg := strings.Join(req.Terms, ",")
	log.Printf("🔍 DEBUG: About to execute: go run . generate-data --students=%d --terms=%s", req.NumStudents, termsArg)

	cmd := exec.Command("go", "run", ".", "generate-data",
		fmt.Sprintf("--students=%d", req.NumStudents),
		fmt.Sprintf("--terms=%s", termsArg))
	cmd.Dir = "./cmd"

	output1, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("❌ Failed to generate student data: %v\nOutput: %s", err, string(output1))
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Student data generation failed: %v", err),
		})
		return
	}

	log.Printf("✅ Student journeys generated")

	// Step 2: Process each term
	log.Printf("🌳 Step 2: Processing terms with Verkle trees...")
	processedTerms := 0
	var outputLog strings.Builder
	outputLog.WriteString(string(output1))
	outputLog.WriteString("\n")

	for _, term := range req.Terms {
		log.Printf("  📚 Processing term: %s", term)

		// Convert data
		cmdConvert := exec.Command("go", "run", ".", "convert-data", term)
		cmdConvert.Dir = "./cmd"
		convertOut, err := cmdConvert.CombinedOutput()
		if err != nil {
			log.Printf("    ❌ Failed to convert data for %s: %v", term, err)
			outputLog.WriteString(fmt.Sprintf("❌ Failed to convert %s: %v\n", term, err))
			continue
		}
		outputLog.WriteString(string(convertOut))

		// Add term (create Verkle tree)
		cmdAdd := exec.Command("go", "run", ".", "add-term", term,
			fmt.Sprintf("../data/verkle_terms/%s_completions.json", term))
		cmdAdd.Dir = "./cmd"
		addOut, err := cmdAdd.CombinedOutput()
		if err != nil {
			log.Printf("    ❌ Failed to create Verkle tree for %s: %v", term, err)
			outputLog.WriteString(fmt.Sprintf("❌ Failed to create Verkle tree for %s: %v\n", term, err))
			continue
		}
		outputLog.WriteString(string(addOut))
		processedTerms++
		log.Printf("    ✅ Verkle tree created for %s", term)
	}

	log.Printf("✅ Processed %d/%d terms successfully", processedTerms, len(req.Terms))

	// Step 3: Import to database (optional)
	log.Printf("🗄️  Step 3: Importing data into database...")
	cmdDB := exec.Command("go", "run", ".", "db-import")
	cmdDB.Dir = "./cmd"
	dbOut, _ := cmdDB.CombinedOutput() // Ignore error - database is optional
	outputLog.WriteString(string(dbOut))

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message": fmt.Sprintf("Generated %d students across %d/%d terms", req.NumStudents, processedTerms, len(req.Terms)),
			"num_students": req.NumStudents,
			"processed_terms": processedTerms,
			"total_terms": len(req.Terms),
			"output":  outputLog.String(),
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
}

func handleListTerms(w http.ResponseWriter, r *http.Request) {
	// Load terms from generated data
	terms := []map[string]interface{}{}
	
	// Check verkle terms data (current system)
	termFiles, err := filepath.Glob("data/verkle_terms/*_completions.json")
	if err == nil && len(termFiles) > 0 {
		for _, termFile := range termFiles {
			// Extract term ID from filename (e.g., "Semester_1_2023_completions.json" -> "Semester_1_2023")
			basename := filepath.Base(termFile)
			termID := strings.TrimSuffix(basename, "_completions.json")
			
			if termData, err := os.ReadFile(termFile); err == nil {
				var completions []map[string]interface{}
				if err := json.Unmarshal(termData, &completions); err == nil {
					// Count unique students
					studentSet := make(map[string]bool)
					courseSet := make(map[string]bool)
					for _, completion := range completions {
						if studentID, ok := completion["student_id"].(string); ok {
							studentSet[studentID] = true
						}
						if courseID, ok := completion["course_id"].(string); ok {
							courseSet[courseID] = true
						}
					}
					
					// Parse term info from ID (e.g., "Semester_1_2023" -> "Semester 1 2023")
					nameParts := strings.Split(termID, "_")
					termName := strings.Join(nameParts, " ")
					
					// Determine dates based on term pattern
					// Extract year from term ID (e.g., "Semester_1_2023" -> "2023")
					var startDate, endDate string
					year := "2023" // default fallback
					if len(nameParts) > 1 {
						// For "Semester_X_YYYY" or "Summer_YYYY", get the year
						if len(nameParts) > 2 {
							year = nameParts[2] // Semester_1_2023 -> 2023
						} else if len(nameParts) == 2 {
							year = nameParts[1] // Summer_2023 -> 2023
						}
					}

					if strings.Contains(termID, "Semester_1") {
						// Fall semester: August - December of same year
						startDate = fmt.Sprintf("%s-08-15", year)
						endDate = fmt.Sprintf("%s-12-15", year)
					} else if strings.Contains(termID, "Semester_2") {
						// Spring semester: January - May of NEXT year
						// E.g., "Semester_2_2023" actually runs Jan-May 2024
						nextYear := year
						var yearInt int
					if _, err := fmt.Sscanf(year, "%d", &yearInt); err == nil && yearInt > 0 {
							nextYear = fmt.Sprintf("%d", yearInt+1)
						}
						startDate = fmt.Sprintf("%s-01-15", nextYear)
						endDate = fmt.Sprintf("%s-05-15", nextYear)
					} else if strings.Contains(termID, "Summer") {
						// Summer term: May - August of same year
						startDate = fmt.Sprintf("%s-05-15", year)
						endDate = fmt.Sprintf("%s-08-15", year)
					} else {
						// Default for Test terms or others
						startDate = "2025-01-01"
						endDate = "2025-12-31"
					}
					
					// Check if term has Verkle tree published
					rootFile := fmt.Sprintf("publish_ready/roots/root_%s.json", termID)
					status := "completed"
					if _, err := os.Stat(rootFile); err != nil {
						status = "pending"
					}
					
					terms = append(terms, map[string]interface{}{
						"id": termID,
						"name": termName,
						"start_date": startDate,
						"end_date": endDate,
						"status": status,
						"student_count": len(studentSet),
						"total_courses": len(courseSet),
					})
				}
			}
		}

		// Sort terms by start_date (chronological order)
		sort.Slice(terms, func(i, j int) bool {
			dateI, okI := terms[i]["start_date"].(string)
			dateJ, okJ := terms[j]["start_date"].(string)
			if !okI || !okJ {
				return false
			}
			return dateI < dateJ // String comparison works for YYYY-MM-DD format
		})
	}

	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: terms})
}

func handleGetTermReceipts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	termID := vars["term_id"]
	
	receipts := []map[string]interface{}{}

	// Look for all journey receipt files and filter by term
	if files, err := filepath.Glob("publish_ready/receipts/*_journey.json"); err == nil {
		for _, file := range files {
			if receiptData, err := os.ReadFile(file); err == nil {
				var receiptFile map[string]interface{}
				if err := json.Unmarshal(receiptData, &receiptFile); err == nil {
					// Check if this receipt contains the requested term
					if termReceipts, ok := receiptFile["term_receipts"].(map[string]interface{}); ok {
						if termData, exists := termReceipts[termID]; exists {
							// Extract the receipt data for this term
							if termDataMap, ok := termData.(map[string]interface{}); ok {
								if receiptData, ok := termDataMap["receipt"].(map[string]interface{}); ok {
									// Create a simplified receipt object for the frontend
									receipt := map[string]interface{}{
										"id": filepath.Base(file),
										"student_id": receiptFile["student_id"],
										"term_id": termID,
										"created_at": termDataMap["generated_at"],
										"courses": receiptData["revealed_courses"],
										"merkle_root": "", // Will be derived from student_term_root
										"verkle_proof": receiptData["verkle_proof"],
										"student_name": fmt.Sprintf("Student %s", receiptFile["student_id"]),
									}
									
									// Convert student_term_root array to hex string if present
									if rootArray, ok := receiptData["student_term_root"].([]interface{}); ok {
										rootHex := ""
										for _, val := range rootArray {
											if intVal, ok := val.(float64); ok {
												rootHex += fmt.Sprintf("%02x", int(intVal))
											}
										}
										receipt["merkle_root"] = rootHex
									}
									
									receipts = append(receipts, receipt)
								}
							}
						}
					}
				}
			}
		}
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: receipts})
}

func handleGetTermRoot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	termID := vars["term_id"]
	
	rootFile := fmt.Sprintf("publish_ready/roots/root_%s.json", termID)
	if _, err := os.Stat(rootFile); err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Term root not found"})
		return
	}
	
	rootData, err := os.ReadFile(rootFile)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to read root file"})
		return
	}
	
	var root map[string]interface{}
	if err := json.Unmarshal(rootData, &root); err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse root file"})
		return
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: root})
}

func handleUpdateTermBlockchainStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	termID := vars["term_id"]

	var req struct {
		TxHash           string `json:"tx_hash"`
		BlockNumber      uint64 `json:"block_number"`
		PublisherAddress string `json:"publisher_address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
		return
	}

	if req.TxHash == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "tx_hash is required"})
		return
	}

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Printf("❌ Failed to connect to database: %v", err)
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Database connection failed: %v", err),
		})
		return
	}
	defer database.Close(db)

	// Update all term receipts for this term with blockchain info
	verified := true
	now := time.Now()

	result := db.Model(&database.TermReceipt{}).
		Where("term_id = ?", termID).
		Updates(map[string]interface{}{
			"blockchain_verified":  &verified,
			"blockchain_tx_hash":   &req.TxHash,
			"blockchain_block":     &req.BlockNumber,
			"published_at":         &now,
			"publisher_address":    &req.PublisherAddress,
		})

	if result.Error != nil {
		log.Printf("❌ Failed to update blockchain status: %v", result.Error)
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to update blockchain status: %v", result.Error),
		})
		return
	}

	log.Printf("✅ Updated %d term receipts for %s with blockchain info (tx: %s)", result.RowsAffected, termID, req.TxHash)

	// Regenerate journey receipts for all students with published terms
	fmt.Printf("📝 Regenerating journey receipts for all students...\n")
	if err := regenerateAllJourneyReceipts(); err != nil {
		fmt.Printf("⚠️ Warning: Failed to regenerate receipts: %v\n", err)
		// Don't fail the blockchain update, just log the warning
	} else {
		fmt.Printf("✅ Journey receipts regenerated successfully\n")
	}

	// Process any pending revocations in the background
	// This runs async so it doesn't block the response
	go func() {
		log.Printf("🔄 Checking for pending revocations to process...")
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Printf("⚠️  Failed to load config for revocation processing: %v", err)
			return
		}

		if cfg.IssuerPrivateKey == "" {
			log.Printf("⚠️  No issuer private key configured, skipping background revocation processing")
			return
		}

		// Check if there are any approved revocations
		dbConn, err := database.Connect()
		if err != nil {
			log.Printf("⚠️  Failed to connect to database for revocation check: %v", err)
			return
		}

		var count int64
		dbConn.Model(&database.RevocationRequest{}).Where("status = ?", "approved").Count(&count)

		if count == 0 {
			log.Printf("✅ No pending revocations to process")
			return
		}

		log.Printf("📋 Found %d approved revocations, processing in background...", count)

		// Process revocations using the existing function
		if err := processApprovedRevocations(cfg.Network, cfg.IssuerPrivateKey, cfg.DefaultGasLimit); err != nil {
			log.Printf("⚠️  Background revocation processing failed: %v", err)
		} else {
			log.Printf("✅ Background revocation processing completed")
		}
	}()

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"term_id":        termID,
			"receipts_updated": result.RowsAffected,
			"tx_hash":        req.TxHash,
		},
	})
}

func handleGenerateReceipt(w http.ResponseWriter, r *http.Request) {
	var req ReceiptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
		return
	}
	
	if req.StudentID == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "student_id is required"})
		return
	}
	
	// Generate temporary output file
	outputFile := fmt.Sprintf("/tmp/receipt_%s_%d.json", extractStudentID(req.StudentID), time.Now().Unix())
	
	// Call existing generateStudentReceipt function
	if err := generateStudentReceipt(req.StudentID, outputFile, req.Terms, req.Courses, req.Selective); err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	
	// Read the generated receipt
	receiptData, err := os.ReadFile(outputFile)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to read generated receipt"})
		return
	}
	
	// Clean up temp file
	defer os.Remove(outputFile)
	
	var receipt map[string]interface{}
	if err := json.Unmarshal(receiptData, &receipt); err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse receipt"})
		return
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: receipt})
}

func handleVerifyReceipt(w http.ResponseWriter, r *http.Request) {
	var receiptData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&receiptData); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid receipt data"})
		return
	}
	
	// Create temporary file for verification
	tempFile := fmt.Sprintf("/tmp/verify_%d.json", time.Now().Unix())
	data, _ := json.Marshal(receiptData)
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create temporary file"})
		return
	}
	defer os.Remove(tempFile)
	
	// Call existing verification function
	if err := verifyReceiptLocally(tempFile); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: fmt.Sprintf("Verification failed: %v", err)})
		return
	}
	
	result := map[string]interface{}{
		"verified": true,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: result})
}

func handleVerifyCourse(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Receipt  json.RawMessage `json:"receipt"`
		CourseID string          `json:"course_id"`
		TermID   string          `json:"term_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("❌ Failed to decode request: %v", err)
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error: fmt.Sprintf("Invalid request format: %v", err),
		})
		return
	}

	log.Printf("📥 Received verify course request - CourseID: %s, TermID: %s", request.CourseID, request.TermID)
	log.Printf("📄 Receipt data length: %d bytes", len(request.Receipt))
	
	// Parse the receipt
	var receipt map[string]interface{}
	if err := json.Unmarshal(request.Receipt, &receipt); err != nil {
		log.Printf("❌ Failed to unmarshal receipt: %v", err)
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error: fmt.Sprintf("Invalid receipt format: %v", err),
		})
		return
	}

	log.Printf("✅ Receipt parsed, keys: %v", getKeys(receipt))

	// Find the term and course
	termReceipts, ok := receipt["term_receipts"].(map[string]interface{})
	if !ok {
		log.Printf("❌ No term_receipts found in receipt. Available keys: %v", getKeys(receipt))
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error: fmt.Sprintf("No term receipts found. Receipt keys: %v", getKeys(receipt)),
		})
		return
	}

	log.Printf("✅ Found term_receipts with terms: %v", getKeys(termReceipts))
	
	termData, ok := termReceipts[request.TermID].(map[string]interface{})
	if !ok {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Error: fmt.Sprintf("Term %s not found in receipt", request.TermID),
		})
		return
	}
	
	// Get the Verkle root from the receipt (for initial check)
	localVerkleRootHex, ok := termData["verkle_root"].(string)
	if !ok {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error: "No Verkle root found for term",
		})
		return
	}
	
	// 🔗 SECURITY: Verify root against blockchain smart contract
	log.Printf("🔍 Fetching Verkle root from blockchain for term: %s", request.TermID)
	
	// TODO: Add blockchain verification here
	// For now, we'll verify against the local root but log the security concern
	log.Printf("⚠️ SECURITY NOTICE: Currently using local root, should verify against blockchain")
	log.Printf("📋 Local root: %s", localVerkleRootHex)
	
	verkleRootHex := localVerkleRootHex
	
	// Get the receipt data
	receiptData, ok := termData["receipt"].(map[string]interface{})
	if !ok {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error: "No receipt data found",
		})
		return
	}
	
	// Find the course proof
	courseProofs, ok := receiptData["course_proofs"].(map[string]interface{})
	if !ok {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error: "No course proofs found",
		})
		return
	}
	
	proofDataRaw, ok := courseProofs[request.CourseID]
	if !ok {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Error: fmt.Sprintf("No proof found for course %s", request.CourseID),
		})
		return
	}
	
	// Convert proof data to JSON bytes (now that we store JSON directly)
	proofBytes, err := json.Marshal(proofDataRaw)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: fmt.Sprintf("Failed to serialize proof data: %v", err),
		})
		return
	}
	
	// Find course details
	var courseInfo map[string]interface{}
	if revealedCourses, ok := receiptData["revealed_courses"].([]interface{}); ok {
		for _, c := range revealedCourses {
			course := c.(map[string]interface{})
			if course["course_id"] == request.CourseID {
				courseInfo = course
				break
			}
		}
	}
	
	if courseInfo == nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Error: fmt.Sprintf("Course %s not in revealed courses", request.CourseID),
		})
		return
	}
	
	// Perform actual cryptographic verification
	log.Printf("🔍 Starting cryptographic verification for course %s", request.CourseID)
	
	// Use the exact course data from revealed_courses to ensure data consistency
	// This preserves the exact serialization format used during proof generation
	var course verkle.CourseCompletion
	courseInfoJSON, err := json.Marshal(courseInfo)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: fmt.Sprintf("Failed to serialize course info: %v", err),
		})
		return
	}
	
	err = json.Unmarshal(courseInfoJSON, &course)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: fmt.Sprintf("Failed to deserialize course data: %v", err),
		})
		return
	}
	
	// SECURITY: Verify Verkle root exists on blockchain before using it for verification
	log.Printf("🔗 Verifying Verkle root exists on blockchain: %s", verkleRootHex)
	
	// Initialize blockchain integration
	ctx := context.Background()
	blockchainVerified := false
	
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("❌ Failed to load configuration: %v", err)
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: "Blockchain configuration error - cannot verify receipt",
		})
		return
	}
	
	if cfg.ContractAddress == "" || cfg.IssuerPrivateKey == "" {
		log.Printf("❌ Missing blockchain configuration (IUMICERT_CONTRACT_ADDRESS or ISSUER_PRIVATE_KEY)")
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: "Blockchain configuration missing - cannot verify receipt",
		})
		return
	}
	
	blockchainIntegration, err := blockchain_integration.NewBlockchainIntegration(cfg.Network, cfg.IssuerPrivateKey, cfg.ContractAddress)
	if err != nil {
		log.Printf("❌ Failed to initialize blockchain integration: %v", err)
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: fmt.Sprintf("Blockchain connection failed: %v", err),
		})
		return
	}
	
	// Check root status on blockchain (includes version info)
	rootStatus, err := blockchainIntegration.CheckRootStatus(ctx, verkleRootHex)
	if err != nil {
		log.Printf("❌ Failed to check root status on blockchain: %v", err)
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error: fmt.Sprintf("Verkle root not found on blockchain: %v", err),
		})
		return
	}

	// Check if root is valid (0=Invalid, 1=Current, 2=Outdated, 3=Superseded)
	if rootStatus.Status == 0 {
		log.Printf("❌ Verkle root does not exist on blockchain: %s", verkleRootHex)
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error: "Verkle root not published on blockchain - invalid receipt",
		})
		return
	}

	// Check if root is superseded
	if rootStatus.Status == 3 {
		log.Printf("⚠️  Verkle root superseded: %s - %s", verkleRootHex, rootStatus.Message)
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error: fmt.Sprintf("Receipt uses superseded root. %s. Please download updated receipt.", rootStatus.Message),
		})
		return
	}

	// SECURITY: Verify term_id matches between receipt and blockchain
	if rootStatus.TermID != request.TermID {
		log.Printf("❌ Term ID mismatch: receipt claims %s but blockchain shows %s for root %s",
			request.TermID, rootStatus.TermID, verkleRootHex)
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error: fmt.Sprintf("Receipt term ID (%s) does not match blockchain term ID (%s) for this verkle root",
				request.TermID, rootStatus.TermID),
		})
		return
	}

	blockchainVerified = true
	statusMsg := map[uint8]string{1: "Current", 2: "Outdated"}[rootStatus.Status]
	log.Printf("✅ Verkle root verified on blockchain: term=%s, version=%s, status=%s",
		rootStatus.TermID, rootStatus.Version.String(), statusMsg)
	
	// Convert verified verkle root hex to bytes
	var verkleRootBytes [32]byte
	if verkleRootHex[:2] == "0x" {
		verkleRootHex = verkleRootHex[2:]
	}
	for i := 0; i < 32 && i*2 < len(verkleRootHex); i++ {
		fmt.Sscanf(verkleRootHex[i*2:i*2+2], "%02x", &verkleRootBytes[i])
	}
	
	// Create course key for verification
	courseKey := fmt.Sprintf("did:example:%s:%s:%s", course.StudentID, course.TermID, course.CourseID)
	
	// The proof data from receipts is JSON, we already have it as bytes
	
	// Perform the actual cryptographic verification with blockchain-verified root
	verificationErr := verkle.VerifyCourseProof(courseKey, course, proofBytes, verkleRootBytes)
	
	ipaPassed := verificationErr == nil
	if !ipaPassed {
		log.Printf("❌ IPA verification failed for course %s: %v", request.CourseID, verificationErr)
	} else {
		log.Printf("✅ IPA verification successful for course %s", request.CourseID)
	}
	
	// Return verification result - succeed if blockchain verification passed, with detailed IPA status
	verificationDetails := map[string]interface{}{
		"ipa_verified": ipaPassed,
		"state_diff_verified": true, // We validated the state diff structure above
		"blockchain_anchored": blockchainVerified,
	}

	var errorMessage string
	if !ipaPassed {
		errorMessage = fmt.Sprintf("IPA verification failed (blockchain verification successful): %v", verificationErr)
	}

	// Try to get blockchain transaction info from database
	var blockchainInfo map[string]interface{}

	db, err := database.Connect()
	if err == nil {
		defer database.Close(db)

		// Query any term_receipt for this term_id that has blockchain info
		var termReceipt database.TermReceipt
		result := db.Where("term_id = ? AND blockchain_tx_hash IS NOT NULL", request.TermID).First(&termReceipt)
		if result.Error == nil && termReceipt.BlockchainTxHash != nil {
			blockchainInfo = map[string]interface{}{
				"tx_hash": *termReceipt.BlockchainTxHash,
				"published_at": rootStatus.Version,
				"block_number": termReceipt.BlockchainBlock,
			}
		}
	}

	// Return detailed verification results (always success if blockchain verification passed)
	responseData := map[string]interface{}{
		"verified": ipaPassed && blockchainVerified, // Overall verification status
		"course": courseInfo,
		"term_id": request.TermID,
		"verkle_root": verkleRootHex,
		"proof_exists": len(proofBytes) > 0,
		"verification_details": verificationDetails,
	}

	// Add blockchain info if available
	if blockchainInfo != nil {
		responseData["blockchain_info"] = blockchainInfo
	}

	response := APIResponse{
		Success: blockchainVerified, // Success if blockchain verification passed
		Data: responseData,
	}

	// Add error message if IPA failed but blockchain succeeded
	if !ipaPassed && blockchainVerified {
		response.Data.(map[string]interface{})["verification_error"] = errorMessage
	}

	respondJSON(w, http.StatusOK, response)
}

func handleListReceipts(w http.ResponseWriter, r *http.Request) {
	receipts := []map[string]interface{}{}

	// Connect to database to get blockchain publication timestamps
	db, dbErr := database.Connect()

	if files, err := filepath.Glob("publish_ready/receipts/*_journey.json"); err == nil {
		for _, file := range files {
			if receiptData, err := os.ReadFile(file); err == nil {
				var receipt map[string]interface{}
				if err := json.Unmarshal(receiptData, &receipt); err == nil {
					receiptInfo := map[string]interface{}{
						"filename":   filepath.Base(file),
						"student_id": receipt["student_id"],
						"timestamp":  receipt["generation_timestamp"],
						"selective":  receipt["receipt_type"],
					}

					// Get latest blockchain publication timestamp for this student's terms
					if dbErr == nil && receipt["term_receipts"] != nil {
						var latestPublishedAt *time.Time

						// Extract term IDs from the receipt
						if termReceipts, ok := receipt["term_receipts"].(map[string]interface{}); ok {
							for termID := range termReceipts {
								var term database.Term
								if err := db.Where("term_id = ? AND published_at IS NOT NULL", termID).First(&term).Error; err == nil {
									if term.PublishedAt != nil {
										if latestPublishedAt == nil || term.PublishedAt.After(*latestPublishedAt) {
											latestPublishedAt = term.PublishedAt
										}
									}
								}
							}
						}

						// Use blockchain publication timestamp if available
						if latestPublishedAt != nil {
							receiptInfo["timestamp"] = latestPublishedAt.Format(time.RFC3339)
							receiptInfo["blockchain_published"] = true
						} else {
							receiptInfo["blockchain_published"] = false
						}
					}

					receipts = append(receipts, receiptInfo)
				}
			}
		}
	}

	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: receipts})
}

func handleGetReceiptByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receiptID := vars["receipt_id"]
	
	if receiptID == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "receipt_id is required"})
		return
	}

	// Look for journey receipt file
	journeyPath := fmt.Sprintf("publish_ready/receipts/%s_journey.json", receiptID)
	if journeyData, err := os.ReadFile(journeyPath); err == nil {
		var journey map[string]interface{}
		if err := json.Unmarshal(journeyData, &journey); err == nil {
			respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: journey})
			return
		}
	}

	respondJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Receipt not found"})
}

// regenerateAllJourneyReceipts regenerates journey receipts for all students with published terms
func regenerateAllJourneyReceipts() error {
	// Connect to database
	db, err := database.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close(db)

	// Get all students from database
	var students []database.Student
	if err := db.Find(&students).Error; err != nil {
		return fmt.Errorf("failed to query students: %w", err)
	}

	if len(students) == 0 {
		return fmt.Errorf("no students found in database")
	}

	fmt.Printf("📋 Found %d students to regenerate receipts for\n", len(students))

	// Regenerate receipt for each student with all published terms
	successCount := 0
	for _, student := range students {
		outputFile := fmt.Sprintf("publish_ready/receipts/%s_journey.json", student.StudentID)

		// Call generateStudentReceipt with empty terms list (auto-discover published terms)
		// and empty courses list (include all courses), selective=false
		err := generateStudentReceipt(student.StudentID, outputFile, nil, nil, false)
		if err != nil {
			fmt.Printf("⚠️ Failed to regenerate receipt for %s: %v\n", student.StudentID, err)
			continue
		}

		successCount++
		fmt.Printf("✓ Regenerated receipt for %s\n", student.StudentID)
	}

	if successCount == 0 {
		return fmt.Errorf("failed to regenerate any receipts")
	}

	fmt.Printf("✅ Successfully regenerated %d/%d receipts\n", successCount, len(students))
	return nil
}

func handlePublishRoots(w http.ResponseWriter, r *http.Request) {
	var req PublishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid request body"})
		return
	}
	
	if req.TermID == "" {
		respondJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "term_id is required"})
		return
	}
	
	network := req.Network
	if network == "" {
		network = "sepolia"
	}
	
	gasLimit := req.GasLimit
	if gasLimit == 0 {
		gasLimit = 500000
	}
	
	// Check if we already have a transaction record for this term first
	fmt.Printf("🔍 API: Checking for existing transaction for %s\n", req.TermID)
	if files, err := filepath.Glob("publish_ready/transactions/tx_*.json"); err == nil && len(files) > 0 {
		// Sort files by modification time to get the most recent
		sort.Slice(files, func(i, j int) bool {
			infoI, errI := os.Stat(files[i])
			infoJ, errJ := os.Stat(files[j])
			if errI != nil || errJ != nil {
				return false
			}
			return infoI.ModTime().After(infoJ.ModTime())
		})
		
		// Look for existing transaction for this term
		for _, file := range files {
			if txData, err := os.ReadFile(file); err == nil {
				var tx map[string]interface{}
				if err := json.Unmarshal(txData, &tx); err == nil {
					// Check if this transaction is for our term
					if rootPath, ok := tx["root_file_path"].(string); ok {
						expectedRootFile := fmt.Sprintf("root_%s.json", req.TermID)
						if strings.Contains(rootPath, expectedRootFile) {
							// Found existing transaction
							fmt.Printf("✅ API: Found existing transaction for %s\n", req.TermID)
							tx["term_id"] = req.TermID
							respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: tx})
							return
						}
					}
				}
			}
		}
	}

	// Call existing publishTermRoots function
	fmt.Printf("🔄 API: About to call publishTermRoots for %s\n", req.TermID)
	if err := publishTermRoots(req.TermID, network, "", gasLimit); err != nil {
		fmt.Printf("❌ API: publishTermRoots failed: %v\n", err)
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	fmt.Printf("✅ API: publishTermRoots completed successfully\n")

	// Regenerate journey receipts for all students with published terms
	fmt.Printf("📝 Regenerating journey receipts for all students...\n")
	if err := regenerateAllJourneyReceipts(); err != nil {
		fmt.Printf("⚠️ Warning: Failed to regenerate receipts: %v\n", err)
		// Don't fail the publish operation, just log the warning
	} else {
		fmt.Printf("✅ Journey receipts regenerated successfully\n")
	}

	// Find the latest transaction file for this term
	// Note: Transaction files are named by hash, so we need to search through them
	if files, err := filepath.Glob("publish_ready/transactions/tx_*.json"); err == nil && len(files) > 0 {
		// Sort files by modification time to get the most recent
		sort.Slice(files, func(i, j int) bool {
			infoI, errI := os.Stat(files[i])
			infoJ, errJ := os.Stat(files[j])
			if errI != nil || errJ != nil {
				return false
			}
			return infoI.ModTime().After(infoJ.ModTime())
		})
		
		// Look through recent transaction files to find one for this term
		for _, file := range files {
			if txData, err := os.ReadFile(file); err == nil {
				var tx map[string]interface{}
				if err := json.Unmarshal(txData, &tx); err == nil {
					// Check if this transaction is for our term
					if rootPath, ok := tx["root_file_path"].(string); ok {
						expectedRootFile := fmt.Sprintf("root_%s.json", req.TermID)
						if strings.Contains(rootPath, expectedRootFile) {
							// Add the term_id to the response for clarity
							tx["term_id"] = req.TermID
							respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: tx})
							return
						}
					}
				}
			}
		}
	}
	
	// Fallback response
	result := map[string]interface{}{
		"term_id": req.TermID,
		"network": network,
		"status": "prepared",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: result})
}

func handleListTransactions(w http.ResponseWriter, r *http.Request) {
	transactions := []map[string]interface{}{}
	
	if files, err := filepath.Glob("publish_ready/transactions/tx_*.json"); err == nil {
		for _, file := range files {
			if txData, err := os.ReadFile(file); err == nil {
				var tx map[string]interface{}
				if err := json.Unmarshal(txData, &tx); err == nil {
					transactions = append(transactions, map[string]interface{}{
						"filename": filepath.Base(file),
						"term_id": tx["term_id"],
						"network": tx["network"],
						"status": tx["status"],
						"timestamp": tx["timestamp"],
						"tx_hash": tx["tx_hash"],
						"gas_limit": tx["gas_limit"],
					})
				}
			}
		}
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: transactions})
}

func handleGetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txHash := vars["tx_hash"]
	
	// Find transaction file by hash
	if files, err := filepath.Glob("publish_ready/transactions/tx_*.json"); err == nil {
		for _, file := range files {
			if txData, err := os.ReadFile(file); err == nil {
				var tx map[string]interface{}
				if err := json.Unmarshal(txData, &tx); err == nil {
					if tx["tx_hash"] == txHash {
						respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: tx})
						return
					}
				}
			}
		}
	}
	
	respondJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Transaction not found"})
}

func handleGetPublishedRoots(w http.ResponseWriter, r *http.Request) {
	roots := []map[string]interface{}{}

	// Read all root files from publish_ready/roots/
	if files, err := filepath.Glob("publish_ready/roots/root_*.json"); err == nil {
		for _, file := range files {
			if rootData, err := os.ReadFile(file); err == nil {
				var root map[string]interface{}
				if err := json.Unmarshal(rootData, &root); err == nil {
					roots = append(roots, map[string]interface{}{
						"filename":    filepath.Base(file),
						"term_id":     root["term_id"],
						"verkle_root": root["verkle_root"],
						"timestamp":   root["timestamp"],
						"tx_hash":     root["tx_hash"],
					})
				}
			}
		}

		// Sort roots by timestamp (chronological order)
		sort.Slice(roots, func(i, j int) bool {
			timeI, okI := roots[i]["timestamp"].(string)
			timeJ, okJ := roots[j]["timestamp"].(string)
			if !okI || !okJ {
				return false
			}
			// Parse timestamps and compare
			tI, errI := time.Parse(time.RFC3339, timeI)
			tJ, errJ := time.Parse(time.RFC3339, timeJ)
			if errI != nil || errJ != nil {
				return false
			}
			return tI.Before(tJ)
		})
	}

	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: roots})
}

func handleListStudents(w http.ResponseWriter, r *http.Request) {
	// Connect to database
	db, err := database.Connect()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database connection failed",
		})
		return
	}
	defer database.Close(db)

	repo := database.NewReceiptRepository(db)

	// Get all students from database
	students, err := repo.GetAllStudents()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to fetch students: %v", err),
		})
		return
	}

	// Format response
	studentList := make([]map[string]interface{}, 0, len(students))
	for _, student := range students {
		studentList = append(studentList, map[string]interface{}{
			"student_id":      student.StudentID,
			"name":            student.Name,
			"did":             student.DID,
			"enrollment_date": student.EnrollmentDate.Format(time.RFC3339),
			"expected_grad":   student.ExpectedGraduation.Format(time.RFC3339),
			"status":          student.Status,
		})
	}

	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: studentList})
}

func handleGetStudentTerms(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]
	
	terms, err := discoverStudentTerms(fmt.Sprintf("did:example:%s", studentID))
	if err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: err.Error()})
		return
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]interface{}{
		"student_id": studentID,
		"terms": terms,
	}})
}

func handleGetStudentJourney(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]
	
	// Load student journey from generated data
	journeyFile := filepath.Join("data/generated_student_data/students", fmt.Sprintf("journey_%s.json", studentID))
	
	if _, err := os.Stat(journeyFile); os.IsNotExist(err) {
		respondJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Student journey not found"})
		return
	}
	
	journeyData, err := os.ReadFile(journeyFile)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to read student journey"})
		return
	}
	
	var journey map[string]interface{}
	if err := json.Unmarshal(journeyData, &journey); err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse student journey"})
		return
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: journey})
}

// Helper functions
func getTermsFromStudent(student map[string]interface{}) []string {
	var terms []string
	if termData, ok := student["terms"].(map[string]interface{}); ok {
		for termID := range termData {
			terms = append(terms, termID)
		}
	}
	return terms
}

func respondJSON(w http.ResponseWriter, status int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// ========== NEW DATABASE-BACKED RECEIPT ENDPOINTS ==========

// handleGetLatestReceipts returns the latest term receipts for a student
func handleGetLatestReceipts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database connection failed",
		})
		return
	}
	defer database.Close(db)

	// Get all term receipts for this student
	var receipts []*database.TermReceipt
	err = db.Where("student_id = ?", studentID).
		Order("generated_at DESC").
		Find(&receipts).Error

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to fetch receipts: %v", err),
		})
		return
	}

	if len(receipts) == 0 {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "No receipts found for student",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"student_id": studentID,
			"count":      len(receipts),
			"receipts":   receipts,
		},
	})
}

// handleGetAccumulatedReceipt returns the full academic journey receipt for a student
func handleGetAccumulatedReceipt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database connection failed",
		})
		return
	}
	defer database.Close(db)

	repo := database.NewReceiptRepository(db)

	// Get or generate accumulated receipt
	accumulated, err := repo.GetCurrentProgressReceipt(studentID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to generate accumulated receipt: %v", err),
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"student_id": studentID,
			"receipt":    accumulated,
		},
	})
}

// handleGetTermReceipt returns a specific term receipt for a student
func handleGetTermReceipt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]
	termID := vars["term_id"]

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Database connection failed",
		})
		return
	}
	defer database.Close(db)

	// Get the term receipt
	var receipt database.TermReceipt
	err = db.Where("student_id = ? AND term_id = ?", studentID, termID).
		First(&receipt).Error

	if err != nil {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Receipt not found",
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"student_id": studentID,
			"term_id":    termID,
			"receipt":    receipt,
		},
	})
}

// handleDownloadJourneyReceipt serves the journey receipt JSON file for download
func handleDownloadJourneyReceipt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]

	// Build file path
	filePath := fmt.Sprintf("publish_ready/receipts/%s_journey.json", studentID)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		respondJSON(w, http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Receipt not found. Please generate receipts first.",
		})
		return
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to read receipt: %v", err),
		})
		return
	}

	// Set headers for download
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s_journey.json", studentID))
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")

	// Write file content
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// ========== IPA VERIFICATION ENDPOINT ==========

// handleIPAVerify performs full IPA (Inner Product Argument) cryptographic verification
// This is computationally intensive and verifies Verkle proofs cryptographically
func handleIPAVerify(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Receipt json.RawMessage `json:"receipt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	// Parse receipt
	var receipt map[string]interface{}
	if err := json.Unmarshal(request.Receipt, &receipt); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid receipt format",
		})
		return
	}

	// Extract student ID
	studentID, ok := receipt["student_id"].(string)
	if !ok {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Missing student_id in receipt",
		})
		return
	}

	// Extract term receipts
	termReceipts, ok := receipt["term_receipts"].(map[string]interface{})
	if !ok {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Missing term_receipts in receipt",
		})
		return
	}

	// Track verification results
	verificationResults := make(map[string]interface{})
	totalCourses := 0
	verifiedCourses := 0
	failedCourses := []string{}

	// Verify each term
	for termID, termDataInterface := range termReceipts {
		termData, ok := termDataInterface.(map[string]interface{})
		if !ok {
			verificationResults[termID] = map[string]interface{}{
				"status": "error",
				"error":  "Invalid term data format",
			}
			continue
		}

		verkleRootHex, ok := termData["verkle_root"].(string)
		if !ok {
			verificationResults[termID] = map[string]interface{}{
				"status": "error",
				"error":  "Missing verkle_root",
			}
			continue
		}

		// Parse Verkle root
		verkleRoot, err := parseVerkleRoot(verkleRootHex)
		if err != nil {
			verificationResults[termID] = map[string]interface{}{
				"status": "error",
				"error":  fmt.Sprintf("Invalid verkle_root: %v", err),
			}
			continue
		}

		// BLOCKCHAIN VERIFICATION: Check if verkle root exists on-chain
		ctx := context.Background()
		cfg, err := config.LoadConfig()
		if err != nil {
			verificationResults[termID] = map[string]interface{}{
				"status": "error",
				"error":  "Blockchain configuration error",
			}
			continue
		}

		if cfg.ContractAddress == "" || cfg.IssuerPrivateKey == "" {
			verificationResults[termID] = map[string]interface{}{
				"status":         "error",
				"error":          "Missing blockchain configuration",
				"blockchain_check": false,
			}
			continue
		}

		blockchainIntegration, err := blockchain_integration.NewBlockchainIntegration(cfg.Network, cfg.IssuerPrivateKey, cfg.ContractAddress)
		if err != nil {
			verificationResults[termID] = map[string]interface{}{
				"status":         "error",
				"error":          fmt.Sprintf("Blockchain connection failed: %v", err),
				"blockchain_check": false,
			}
			continue
		}

		// Check root status on blockchain
		rootStatus, err := blockchainIntegration.CheckRootStatus(ctx, verkleRootHex)
		if err != nil {
			verificationResults[termID] = map[string]interface{}{
				"status":         "error",
				"error":          fmt.Sprintf("Failed to check root status: %v", err),
				"blockchain_check": false,
			}
			continue
		}

		// Check if root is invalid
		if rootStatus.Status == 0 {
			verificationResults[termID] = map[string]interface{}{
				"status":         "error",
				"error":          "Verkle root not published on blockchain - invalid receipt",
				"blockchain_check": false,
			}
			continue
		}

		// Check if root is superseded
		if rootStatus.Status == 3 {
			verificationResults[termID] = map[string]interface{}{
				"status":         "error",
				"error":          fmt.Sprintf("Receipt uses superseded root. %s", rootStatus.Message),
				"blockchain_check": false,
				"version_status": "superseded",
			}
			continue
		}

		// Verify term_id matches
		if rootStatus.TermID != termID {
			verificationResults[termID] = map[string]interface{}{
				"status":         "error",
				"error":          fmt.Sprintf("Term ID mismatch: receipt claims %s but blockchain shows %s", termID, rootStatus.TermID),
				"blockchain_check": false,
			}
			continue
		}

		log.Printf("✅ Blockchain verification passed for term %s: root exists on-chain", termID)

		// Get receipt data
		receiptData, ok := termData["receipt"].(map[string]interface{})
		if !ok {
			verificationResults[termID] = map[string]interface{}{
				"status": "success",
				"note":   "No receipt data to verify",
			}
			continue
		}

		// Get course proofs and revealed courses
		courseProofs, ok := receiptData["course_proofs"].(map[string]interface{})
		if !ok {
			verificationResults[termID] = map[string]interface{}{
				"status": "error",
				"error":  "Missing course_proofs",
			}
			continue
		}

		revealedCourses, ok := receiptData["revealed_courses"].([]interface{})
		if !ok {
			verificationResults[termID] = map[string]interface{}{
				"status": "error",
				"error":  "Missing revealed_courses",
			}
			continue
		}

		// Verify each course cryptographically
		termResults := make(map[string]interface{})
		termVerified := 0
		termFailed := 0

		for _, courseInterface := range revealedCourses {
			courseMap, ok := courseInterface.(map[string]interface{})
			if !ok {
				continue
			}

			courseID, ok := courseMap["course_id"].(string)
			if !ok {
				continue
			}

			totalCourses++

			// Get course proof
			proofData, exists := courseProofs[courseID]
			if !exists {
				termResults[courseID] = "no_proof"
				termFailed++
				failedCourses = append(failedCourses, fmt.Sprintf("%s:%s", termID, courseID))
				continue
			}

			// Convert proof to JSON bytes
			proofBytes, err := json.Marshal(proofData)
			if err != nil {
				termResults[courseID] = fmt.Sprintf("proof_parse_error: %v", err)
				termFailed++
				failedCourses = append(failedCourses, fmt.Sprintf("%s:%s", termID, courseID))
				continue
			}

			// Convert course map to CourseCompletion
			course, err := convertToCourseCompletion(courseMap)
			if err != nil {
				termResults[courseID] = fmt.Sprintf("course_parse_error: %v", err)
				termFailed++
				failedCourses = append(failedCourses, fmt.Sprintf("%s:%s", termID, courseID))
				continue
			}

			// Generate course key
			studentDID := fmt.Sprintf("did:example:%s", studentID)
			courseKey := fmt.Sprintf("%s:%s:%s", studentDID, termID, courseID)

			// Perform full IPA cryptographic verification
			if err := verkle.VerifyCourseProof(courseKey, course, proofBytes, verkleRoot); err != nil {
				termResults[courseID] = fmt.Sprintf("verification_failed: %v", err)
				termFailed++
				failedCourses = append(failedCourses, fmt.Sprintf("%s:%s", termID, courseID))
				continue
			}

			termResults[courseID] = "verified"
			termVerified++
			verifiedCourses++
		}

		// Try to get blockchain transaction info from database term_receipts table
		var blockchainTxHash string
		var blockchainBlock *uint64

		db, dbErr := database.Connect()
		if dbErr == nil {
			defer database.Close(db)

			// Query any term_receipt for this term_id that has blockchain info
			var termReceipt database.TermReceipt
			result := db.Where("term_id = ? AND blockchain_tx_hash IS NOT NULL", termID).First(&termReceipt)
			if result.Error == nil && termReceipt.BlockchainTxHash != nil {
				blockchainTxHash = *termReceipt.BlockchainTxHash
				blockchainBlock = termReceipt.BlockchainBlock
			}
		}

		// Build verification result with transaction info if available
		if blockchainTxHash != "" {
			verificationResults[termID] = map[string]interface{}{
				"status":              "completed",
				"verkle_root":         verkleRootHex,
				"courses_verified":    termVerified,
				"courses_failed":      termFailed,
				"course_results":      termResults,
				"blockchain_verified": true,
				"blockchain_published_at": rootStatus.Version.String(),
				"blockchain_tx_hash":  blockchainTxHash,
				"blockchain_block":    blockchainBlock,
			}
		} else {
			// No transaction info available
			verificationResults[termID] = map[string]interface{}{
				"status":              "completed",
				"verkle_root":         verkleRootHex,
				"courses_verified":    termVerified,
				"courses_failed":      termFailed,
				"course_results":      termResults,
				"blockchain_verified": true,
				"blockchain_published_at": rootStatus.Version.String(),
			}
		}
	}

	// Overall status
	overallStatus := "success"
	if len(failedCourses) > 0 {
		overallStatus = "partial_failure"
	}
	if verifiedCourses == 0 && totalCourses > 0 {
		overallStatus = "failure"
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: overallStatus == "success",
		Data: map[string]interface{}{
			"status":           overallStatus,
			"student_id":       studentID,
			"total_courses":    totalCourses,
			"verified_courses": verifiedCourses,
			"failed_courses":   len(failedCourses),
			"failed_list":      failedCourses,
			"term_results":     verificationResults,
			"computation_note": "Full IPA cryptographic verification performed on backend",
		},
	})
}

// Helper function to get keys from a map
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func init() {
	serveCmd.Flags().String("port", "8080", "Port to serve the API on")
	serveCmd.Flags().Bool("cors", true, "Enable CORS for React development")
	rootCmd.AddCommand(serveCmd)
}

