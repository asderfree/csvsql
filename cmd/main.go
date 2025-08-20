package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"csvsql/config"
	"csvsql/internal/database"
	"csvsql/internal/importer"
	"csvsql/internal/mapping"
	"csvsql/internal/repl"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Expect file paths as command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: csvsql <file1.csv> [file2.xlsx] ...")
		os.Exit(1)
	}

	// Use configured database (default: in-memory SQLite)
	db, err := sql.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Initialize components with dependency injection
	mapper := mapping.NewMapper()
	dbManager := database.NewManager(db, mapper)
	processor := importer.NewProcessor(dbManager)
	commands := repl.NewCommands(dbManager)
	formatter := repl.NewFormatter()
	session := repl.NewSession(commands, formatter)

	// Load all files provided as arguments
	for _, filePath := range os.Args[1:] {
		if err := processor.LoadFile(filePath); err != nil {
			log.Printf("Warning: Failed to load file %s: %v", filePath, err)
		}
	}

	// Start the interactive Read-Eval-Print Loop (REPL)
	session.Run()
}
