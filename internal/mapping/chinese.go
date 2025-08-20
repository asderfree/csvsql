package mapping

import (
	"fmt"
	"maps"
	"regexp"
)

// Mapper handles Chinese header to column name mappings
type Mapper struct {
	chineseToColumn map[string]map[string]string // table -> chineseHeader -> columnName
}

// NewMapper creates a new Chinese header mapper
func NewMapper() *Mapper {
	return &Mapper{
		chineseToColumn: make(map[string]map[string]string),
	}
}

// AddMapping adds a mapping for a table
func (m *Mapper) AddMapping(tableName, chineseHeader, columnName string) {
	if m.chineseToColumn[tableName] == nil {
		m.chineseToColumn[tableName] = make(map[string]string)
	}
	m.chineseToColumn[tableName][chineseHeader] = columnName
}

// GetColumnName gets the column name for a Chinese header in a table
func (m *Mapper) GetColumnName(tableName, chineseHeader string) (string, bool) {
	if tableMappings, exists := m.chineseToColumn[tableName]; exists {
		if columnName, found := tableMappings[chineseHeader]; found {
			return columnName, true
		}
	}
	return "", false
}

// GetChineseHeader gets the Chinese header for a column name in a table
func (m *Mapper) GetChineseHeader(tableName, columnName string) (string, bool) {
	if tableMappings, exists := m.chineseToColumn[tableName]; exists {
		for chineseHeader, colName := range tableMappings {
			if colName == columnName {
				return chineseHeader, true
			}
		}
	}
	return "", false
}

// TranslateQuery replaces Chinese field names in SQL queries with their corresponding column names
func (m *Mapper) TranslateQuery(query string) string {
	// For each table, check if any Chinese headers are used in the query
	for _, tableMappings := range m.chineseToColumn {
		for chineseHeader, columnName := range tableMappings {
			// Replace Chinese field names in the query
			// Handle different SQL contexts: SELECT, WHERE, ORDER BY, etc.
			patterns := []string{
				fmt.Sprintf("`%s`", regexp.QuoteMeta(chineseHeader)),   // Backticks
				fmt.Sprintf("\"%s\"", regexp.QuoteMeta(chineseHeader)), // Double quotes
				fmt.Sprintf("'%s'", regexp.QuoteMeta(chineseHeader)),   // Single quotes
				chineseHeader,
			}

			// Handle quoted patterns first
			for _, pattern := range patterns {
				re := regexp.MustCompile(pattern)
				query = re.ReplaceAllString(query, columnName)
			}
		}
	}
	return query
}

// RestoreHeaders replaces sanitized column names with original Chinese headers
func (m *Mapper) RestoreHeaders(columns []string) []string {
	restoredColumns := make([]string, len(columns))
	copy(restoredColumns, columns)

	// Check all tables for mappings
	for _, tableMappings := range m.chineseToColumn {
		for i, col := range restoredColumns {
			// Find the Chinese header that maps to this column name
			for chineseHeader, columnName := range tableMappings {
				if columnName == col {
					restoredColumns[i] = chineseHeader
					break
				}
			}
		}
	}

	return restoredColumns
}

// GetMappings returns all mappings for debugging
func (m *Mapper) GetMappings() map[string]map[string]string {
	result := make(map[string]map[string]string)
	for tableName, tableMappings := range m.chineseToColumn {
		result[tableName] = make(map[string]string)
		maps.Copy(result[tableName], tableMappings)
	}
	return result
}

// GetTableMappings returns mappings for a specific table
func (m *Mapper) GetTableMappings(tableName string) map[string]string {
	if tableMappings, exists := m.chineseToColumn[tableName]; exists {
		result := make(map[string]string)
		maps.Copy(result, tableMappings)
		return result
	}
	return nil
}
