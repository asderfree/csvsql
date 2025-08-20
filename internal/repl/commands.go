package repl

import (
	"fmt"
	"strings"

	"csvsql/internal/database"
	"csvsql/internal/mapping"
)

// Commands handles REPL command processing
type Commands struct {
	dbManager *database.Manager
	mapper    *mapping.Mapper
}

// NewCommands creates a new commands handler
func NewCommands(dbManager *database.Manager) *Commands {
	return &Commands{
		dbManager: dbManager,
		mapper:    dbManager.GetMapper(),
	}
}

// ProcessCommand processes a REPL command and returns results
func (c *Commands) ProcessCommand(input string) (CommandResult, error) {
	input = strings.TrimSpace(input)

	switch strings.ToLower(input) {
	case "exit", "quit", ".exit", ".quit":
		return CommandResult{Type: ExitCommand}, nil
	case ".help":
		return CommandResult{Type: HelpCommand, Data: getHelpText()}, nil
	case ".tables":
		return c.handleTablesCommand()
	case ".mappings":
		return c.handleMappingsCommand()
	}

	if strings.HasPrefix(strings.ToLower(input), ".schema ") {
		return c.handleSchemaCommand(input)
	}

	if strings.HasPrefix(strings.ToUpper(input), "EXPORT ") {
		return c.handleExportCommand(input)
	}

	// Default: treat as SQL query
	return c.handleSQLQuery(input)
}

// CommandType represents the type of command
type CommandType int

const (
	SQLQueryCommand CommandType = iota
	HelpCommand
	TablesCommand
	SchemaCommand
	MappingsCommand
	ExportCommand
	ExitCommand
)

// CommandResult represents the result of processing a command
type CommandResult struct {
	Type  CommandType
	Data  interface{}
	Error error
}

func (c *Commands) handleTablesCommand() (CommandResult, error) {
	results, err := c.dbManager.ExecuteQuery("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		return CommandResult{}, err
	}
	return CommandResult{Type: TablesCommand, Data: results}, nil
}

func (c *Commands) handleSchemaCommand(input string) (CommandResult, error) {
	tableName := strings.TrimSpace(strings.TrimPrefix(input, ".schema "))
	query := fmt.Sprintf("PRAGMA table_info(%s);", tableName)
	results, err := c.dbManager.ExecuteQuery(query)
	if err != nil {
		return CommandResult{}, err
	}
	return CommandResult{Type: SchemaCommand, Data: results}, nil
}

func (c *Commands) handleMappingsCommand() (CommandResult, error) {
	mappings := c.mapper.GetMappings()
	return CommandResult{Type: MappingsCommand, Data: mappings}, nil
}

func (c *Commands) handleExportCommand(input string) (CommandResult, error) {
	parts := strings.Fields(input)
	if len(parts) != 2 {
		return CommandResult{}, fmt.Errorf("invalid EXPORT command. Usage: EXPORT <filename.csv>")
	}
	return CommandResult{Type: ExportCommand, Data: parts[1]}, nil
}

func (c *Commands) handleSQLQuery(query string) (CommandResult, error) {
	results, err := c.dbManager.ExecuteQuery(query)
	if err != nil {
		return CommandResult{}, err
	}
	return CommandResult{Type: SQLQueryCommand, Data: results}, nil
}

func getHelpText() string {
	return `Commands:
  .help              Show this help message.
  .tables            List available tables.
  .schema <table>    Show the schema for a table.
  .mappings          Show Chinese header to column name mappings.
  .exit, .quit       Exit the application.
  EXPORT <file.csv>  Export the last SELECT query results to a CSV file.
  Any other text is treated as an SQL query.`
}
