package repl

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// Formatter handles output formatting for the REPL
type Formatter struct{}

// NewFormatter creates a new formatter
func NewFormatter() *Formatter {
	return &Formatter{}
}

// PrintResults formats and prints data to the console
func (f *Formatter) PrintResults(data [][]string) {
	if len(data) <= 1 {
		fmt.Println("Query OK, 0 rows returned.")
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	// Write headers
	fmt.Fprintln(writer, strings.Join(data[0], "\t"))
	// Write separator
	headerSeparators := make([]string, len(data[0]))
	for i, h := range data[0] {
		headerSeparators[i] = strings.Repeat("-", len(h))
	}
	fmt.Fprintln(writer, strings.Join(headerSeparators, "\t"))

	// Write rows
	for _, row := range data[1:] {
		fmt.Fprintln(writer, strings.Join(row, "\t"))
	}
	writer.Flush()
	fmt.Printf("\n(%d rows)\n", len(data)-1)
}

// PrintMappings displays the current Chinese header mappings for all tables
func (f *Formatter) PrintMappings(mappings map[string]map[string]string) {
	if len(mappings) == 0 {
		fmt.Println("No Chinese header mappings found.")
		return
	}

	fmt.Print("Chinese Header Mappings:\n")
	fmt.Print("========================\n")

	for tableName, tableMappings := range mappings {
		if len(tableMappings) > 0 {
			fmt.Printf("Table: %s\n", tableName)
			fmt.Print("-------------------\n")
			for chineseHeader, columnName := range tableMappings {
				fmt.Printf("  %s -> %s\n", chineseHeader, columnName)
			}
		}
	}
}

// ExportToCSV saves query results to a CSV file
func (f *Formatter) ExportToCSV(filename string, data [][]string) error {
	if len(data) <= 1 {
		return fmt.Errorf("no results to export. run a SELECT query first")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	fmt.Printf("Exporting %d rows to %s...\n", len(data)-1, filename)
	return writer.WriteAll(data)
}
