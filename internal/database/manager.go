package database

import (
	"database/sql"
	"fmt"
	"strings"

	"csvsql/internal/mapping"
	"csvsql/pkg/utils"
)

// Manager handles database operations
type Manager struct {
	db     *sql.DB
	mapper *mapping.Mapper
}

// NewManager creates a new database manager
func NewManager(db *sql.DB, mapper *mapping.Mapper) *Manager {
	return &Manager{
		db:     db,
		mapper: mapper,
	}
}

// CreateAndInsert creates a table and inserts data using a transaction
func (m *Manager) CreateAndInsert(tableName string, data [][]string) error {
	if len(data) == 0 {
		return fmt.Errorf("cannot create table from empty data")
	}

	headers := data[0]

	// Initialize mapping for this table (no need to add empty mapping)

	// Sanitize headers for use as column names
	sanitizedHeaders := make(map[string]int)

	for i, h := range headers {
		originalHeader := h

		// Check if header contains Chinese characters
		if utils.ContainsChinese(h) {
			// Generate a new column name with index
			columnIndex := i + 1
			newColumnName := fmt.Sprintf("_%d", columnIndex)

			// Store the mapping from Chinese header to column name
			m.mapper.AddMapping(tableName, originalHeader, newColumnName)

			// Update the header to use the new column name
			headers[i] = newColumnName
		} else {
			// Handle non-Chinese headers as before
			sanitizedHeader := utils.SanitizeColumnName(h)
			if count, ok := sanitizedHeaders[sanitizedHeader]; ok {
				sanitizedHeaders[sanitizedHeader] = count + 1
				headers[i] = fmt.Sprintf("%s_%d", sanitizedHeader, count+1)
			} else {
				sanitizedHeaders[sanitizedHeader] = 1
				headers[i] = sanitizedHeader
			}
		}
	}

	// Create the table
	columnDefs := make([]string, len(headers))
	for i, h := range headers {
		columnDefs[i] = h + " TEXT"
	}
	query := fmt.Sprintf("CREATE TABLE %s (%s);", tableName, strings.Join(columnDefs, ", "))
	if _, err := m.db.Exec(query); err != nil {
		return fmt.Errorf("create table failed: %w", err)
	}

	// Insert data in a transaction for efficiency
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	placeholders := strings.Repeat("?,", len(headers))
	placeholders = placeholders[:len(placeholders)-1] // remove trailing comma
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s VALUES (%s)", tableName, placeholders))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, row := range data[1:] { // Skip header row
		rowInterface := make([]interface{}, len(row))
		for i, v := range row {
			rowInterface[i] = v
		}
		_, err := stmt.Exec(rowInterface...)
		if err != nil {
			tx.Rollback() // Rollback on any error
			return err
		}
	}

	return tx.Commit()
}

// ExecuteQuery runs the user's SQL query and returns the result
func (m *Manager) ExecuteQuery(query string) ([][]string, error) {
	// Translate Chinese field names before executing the query
	translatedQuery := m.mapper.TranslateQuery(query)

	trimmedQuery := strings.ToUpper(strings.TrimSpace(translatedQuery))
	// For non-SELECT queries (INSERT, UPDATE, DELETE)
	if !strings.HasPrefix(trimmedQuery, "SELECT") && !strings.HasPrefix(trimmedQuery, "PRAGMA") {
		res, err := m.db.Exec(translatedQuery)
		if err != nil {
			return nil, err
		}
		affected, _ := res.RowsAffected()
		return nil, fmt.Errorf("query OK, %d rows affected", affected)
	}

	// For SELECT queries
	rows, err := m.db.Query(translatedQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Restore Chinese headers by reversing the mapping
	restoredColumns := m.mapper.RestoreHeaders(columns)

	var resultsData [][]string
	resultsData = append(resultsData, restoredColumns) // Add restored headers as the first row

	for rows.Next() {
		rowValues := make([]interface{}, len(columns))
		rowScanners := make([]interface{}, len(columns))
		for i := range rowValues {
			rowScanners[i] = &rowValues[i]
		}

		if err := rows.Scan(rowScanners...); err != nil {
			return nil, err
		}

		rowStr := make([]string, len(columns))
		for i, val := range rowValues {
			if val == nil {
				rowStr[i] = "NULL"
			} else {
				rowStr[i] = fmt.Sprintf("%s", val)
			}
		}
		resultsData = append(resultsData, rowStr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resultsData, nil
}

// GetMapper returns the Chinese header mapper
func (m *Manager) GetMapper() *mapping.Mapper {
	return m.mapper
}
