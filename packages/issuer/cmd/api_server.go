package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"iumicert/crypto/merkle"
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
	Courses    []merkle.CourseCompletion `json:"courses,omitempty"`
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

func startAPIServer(port string, corsEnabled bool) error {
	r := mux.NewRouter()
	
	// API routes
	api := r.PathPrefix("/api").Subrouter()
	
	// System endpoints
	api.HandleFunc("/status", handleSystemStatus).Methods("GET")
	api.HandleFunc("/health", handleHealth).Methods("GET")
	
	// Term management
	api.HandleFunc("/terms", handleAddTerm).Methods("POST")
	api.HandleFunc("/terms", handleListTerms).Methods("GET")
	api.HandleFunc("/terms/{term_id}/receipts", handleGetTermReceipts).Methods("GET")
	api.HandleFunc("/terms/{term_id}/roots", handleGetTermRoot).Methods("GET")
	
	// Receipt generation
	api.HandleFunc("/receipts", handleGenerateReceipt).Methods("POST")
	api.HandleFunc("/receipts/verify", handleVerifyReceipt).Methods("POST")
	api.HandleFunc("/receipts", handleListReceipts).Methods("GET")
	
	// Blockchain operations
	api.HandleFunc("/blockchain/publish", handlePublishRoots).Methods("POST")
	api.HandleFunc("/blockchain/transactions", handleListTransactions).Methods("GET")
	api.HandleFunc("/blockchain/transactions/{tx_hash}", handleGetTransaction).Methods("GET")
	
	// Student data operations
	api.HandleFunc("/students", handleListStudents).Methods("GET")
	api.HandleFunc("/students/{student_id}/terms", handleGetStudentTerms).Methods("GET")
	api.HandleFunc("/students/{student_id}/journey", handleGetStudentJourney).Methods("GET")
	
	// Setup CORS if enabled
	var handler http.Handler = r
	if corsEnabled {
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"}, // React dev servers
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"*"},
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
		{"blockchain_ready/receipts", &status.Storage.Receipts},
		{"blockchain_ready/transactions", &status.Storage.Transactions},
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
	var completions []merkle.CourseCompletion
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

func handleListTerms(w http.ResponseWriter, r *http.Request) {
	// Load terms from generated data
	terms := []map[string]interface{}{}
	
	// Check generated terms data (now with IU Vietnam realistic data)
	termFiles, err := filepath.Glob("data/generated_student_data/terms/summary_*.json")
	if err == nil && len(termFiles) > 0 {
		for _, termFile := range termFiles {
			if termData, err := os.ReadFile(termFile); err == nil {
				var termSummary map[string]interface{}
				if err := json.Unmarshal(termData, &termSummary); err == nil {
					termID, ok := termSummary["term_id"].(string)
					if !ok {
						continue
					}
					
					// Parse term info from ID (e.g., "Fall_2024" -> "Fall 2024")
					nameParts := strings.Split(termID, "_")
					termName := strings.Join(nameParts, " ")
					
					// Determine dates based on term
					var startDate, endDate string
					if strings.Contains(termID, "Fall") {
						year := "2024"
						if len(nameParts) > 1 {
							year = nameParts[1]
						}
						startDate = fmt.Sprintf("%s-08-15", year)
						endDate = fmt.Sprintf("%s-12-15", year)
					} else if strings.Contains(termID, "Spring") {
						year := "2025"
						if len(nameParts) > 1 {
							year = nameParts[1]
						}
						startDate = fmt.Sprintf("%s-01-15", year)
						endDate = fmt.Sprintf("%s-05-15", year)
					}
					
					// Check if term has merkle trees
					treeFiles, _ := filepath.Glob(fmt.Sprintf("data/merkle_trees/%s/merkle_*.json", termID))
					status := "completed"
					if len(treeFiles) == 0 {
						status = "pending"
					}
					
					terms = append(terms, map[string]interface{}{
						"id": termID,
						"name": termName,
						"start_date": startDate,
						"end_date": endDate,
						"status": status,
						"student_count": termSummary["total_students"],
						"total_courses": termSummary["total_courses"],
					})
				}
			}
		}
	}
	
	// If no generated terms, fall back to merkle trees directory
	if len(terms) == 0 {
		if files, err := filepath.Glob("data/merkle_trees/*/"); err == nil {
			for _, dir := range files {
				termID := filepath.Base(dir)
				if termID == "" {
					continue
				}
				
				// Count merkle trees in this term
				treeFiles, _ := filepath.Glob(filepath.Join(dir, "merkle_*.json"))
				
				terms = append(terms, map[string]interface{}{
					"id": termID,
					"name": strings.ReplaceAll(termID, "_", " "),
					"start_date": "2024-01-01",
					"end_date": "2024-05-31",
					"status": "active",
					"student_count": len(treeFiles),
				})
			}
		}
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: terms})
}

func handleGetTermReceipts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	termID := vars["term_id"]
	
	receipts := []map[string]interface{}{}
	
	// Look for all receipt files and filter by term
	if files, err := filepath.Glob("blockchain_ready/receipts/receipt_*.json"); err == nil {
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
	
	rootFile := fmt.Sprintf("blockchain_ready/roots/root_%s.json", termID)
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

func handleListReceipts(w http.ResponseWriter, r *http.Request) {
	receipts := []map[string]interface{}{}
	
	if files, err := filepath.Glob("blockchain_ready/receipts/receipt_*.json"); err == nil {
		for _, file := range files {
			if receiptData, err := os.ReadFile(file); err == nil {
				var receipt map[string]interface{}
				if err := json.Unmarshal(receiptData, &receipt); err == nil {
					receipts = append(receipts, map[string]interface{}{
						"filename": filepath.Base(file),
						"student_id": receipt["student_id"],
						"timestamp": receipt["generation_timestamp"],
						"selective": receipt["receipt_type"],
					})
				}
			}
		}
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: receipts})
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
	
	// Call existing publishTermRoots function
	if err := publishTermRoots(req.TermID, network, "", gasLimit); err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	
	// Find the latest transaction file for this term
	pattern := fmt.Sprintf("blockchain_ready/transactions/tx_%s_*.json", req.TermID)
	if files, err := filepath.Glob(pattern); err == nil && len(files) > 0 {
		// Get the most recent file
		latestFile := files[len(files)-1]
		if txData, err := os.ReadFile(latestFile); err == nil {
			var tx map[string]interface{}
			if err := json.Unmarshal(txData, &tx); err == nil {
				respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: tx})
				return
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
	
	if files, err := filepath.Glob("blockchain_ready/transactions/tx_*.json"); err == nil {
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
	if files, err := filepath.Glob("blockchain_ready/transactions/tx_*.json"); err == nil {
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

func handleListStudents(w http.ResponseWriter, r *http.Request) {
	students := []map[string]interface{}{}
	
	// Get students from generated data if available
	if files, err := filepath.Glob("data/generated_student_data/students/journey_*.json"); err == nil {
		for _, file := range files {
			if studentData, err := os.ReadFile(file); err == nil {
				var student map[string]interface{}
				if err := json.Unmarshal(studentData, &student); err == nil {
					// Extract student ID from filename
					filename := filepath.Base(file)
					studentID := strings.TrimSuffix(strings.TrimPrefix(filename, "journey_"), ".json")
					
					students = append(students, map[string]interface{}{
						"student_id": studentID,
						"student_did": student["student_id"],
						"terms": getTermsFromStudent(student),
					})
				}
			}
		}
	}
	
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Data: students})
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

func init() {
	serveCmd.Flags().String("port", "8080", "Port to serve the API on")
	serveCmd.Flags().Bool("cors", true, "Enable CORS for React development")
	rootCmd.AddCommand(serveCmd)
}