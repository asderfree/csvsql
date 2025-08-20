package importer

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// ReadXLSX reads all records from the first sheet of an Excel file
func ReadXLSX(filePath string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in excel file")
	}
	// Read from the first sheet
	return f.GetRows(sheets[0])
}
