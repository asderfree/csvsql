package importer

import (
	"fmt"
	"path/filepath"
	"strings"

	"csvsql/internal/database"
	"csvsql/pkg/utils"
)

// Processor handles file loading and processing
type Processor struct {
	dbManager *database.Manager
}

// NewProcessor creates a new file processor
func NewProcessor(dbManager *database.Manager) *Processor {
	return &Processor{
		dbManager: dbManager,
	}
}

// LoadFile dispatches to the correct parser based on file extension
func (p *Processor) LoadFile(filePath string) error {
	// Get the original file path for reading, only sanitize the table name
	tableName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	// Sanitize table name to be valid SQL
	tableName = utils.SanitizeTableName(tableName)

	var data [][]string
	var err error

	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".csv":
		data, err = p.readCSV(filePath)
	case ".xlsx":
		data, err = p.readXLSX(filePath)
	default:
		return fmt.Errorf("unsupported file type: %s", filePath)
	}

	if err != nil {
		return err
	}

	if len(data) < 1 {
		return fmt.Errorf("no data found in file: %s", filePath)
	}

	if err := p.dbManager.CreateAndInsert(tableName, data); err != nil {
		return fmt.Errorf("failed to load data into table %s: %v", tableName, err)
	}

	fmt.Printf("Successfully loaded table '%s' from %s.\n", tableName, filePath)
	return nil
}

// readCSV reads all records from a CSV file
func (p *Processor) readCSV(filePath string) ([][]string, error) {
	return ReadCSV(filePath)
}

// readXLSX reads all records from the first sheet of an Excel file
func (p *Processor) readXLSX(filePath string) ([][]string, error) {
	return ReadXLSX(filePath)
}
