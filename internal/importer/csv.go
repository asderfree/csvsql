package importer

import (
	"encoding/csv"
	"os"
)

// ReadCSV reads all records from a CSV file
func ReadCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}
