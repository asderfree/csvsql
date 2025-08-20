package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Session manages the REPL session
type Session struct {
	commands    *Commands
	formatter   *Formatter
	lastResults [][]string
}

// NewSession creates a new REPL session
func NewSession(commands *Commands, formatter *Formatter) *Session {
	return &Session{
		commands:  commands,
		formatter: formatter,
	}
}

// Run starts the interactive REPL session
func (s *Session) Run() {
	fmt.Println("\nEnter SQL commands or type .help for help.")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("sql> ")

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		result, err := s.commands.ProcessCommand(input)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Print("sql> ")
			continue
		}

		// Check if we should exit
		if result.Type == ExitCommand {
			return
		}

		s.handleCommandResult(result)
		fmt.Print("sql> ")
	}
}

// handleCommandResult processes the result of a command
func (s *Session) handleCommandResult(result CommandResult) {
	switch result.Type {
	case ExitCommand:
		return
	case HelpCommand:
		if helpText, ok := result.Data.(string); ok {
			fmt.Println(helpText)
		}
	case TablesCommand, SchemaCommand, SQLQueryCommand:
		if data, ok := result.Data.([][]string); ok {
			s.lastResults = data
			s.formatter.PrintResults(data)
		}
	case MappingsCommand:
		if mappings, ok := result.Data.(map[string]map[string]string); ok {
			s.formatter.PrintMappings(mappings)
		}
	case ExportCommand:
		if filename, ok := result.Data.(string); ok {
			if s.lastResults != nil {
				if err := s.formatter.ExportToCSV(filename, s.lastResults); err != nil {
					fmt.Println("Error exporting to CSV:", err)
				}
			} else {
				fmt.Println("No results to export. Run a SELECT query first.")
			}
		}
	}
}
