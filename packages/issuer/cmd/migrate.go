package main

import (
	"fmt"
	"iumicert/issuer/database"
	"log"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  "Create or update the database schema based on the defined models",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func runMigrations() {
	fmt.Println("🔄 Starting database migration...")

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	// Create additional indexes
	if err := database.CreateIndexes(db); err != nil {
		log.Printf("⚠️  Warning: Some indexes may not have been created: %v", err)
	}

	fmt.Println("✅ Database migration completed successfully!")
	fmt.Println("")
	fmt.Println("📋 Tables created:")
	fmt.Println("  - students")
	fmt.Println("  - terms")
	fmt.Println("  - term_receipts")
	fmt.Println("  - accumulated_receipts")
	fmt.Println("  - verification_logs")
	fmt.Println("  - blockchain_transactions")
}
