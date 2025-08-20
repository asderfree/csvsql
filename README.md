# CSVSQL - use csv,xlsx as the sql in your terminal

A Go-based tool for loading CSV and Excel files into an in-memory SQLite database with support for Chinese headers and an interactive SQL REPL.

## Features

- **Multi-format Support**: Load CSV and Excel (.xlsx) files
- **Chinese Header Support**: Automatically handles Chinese column headers by mapping them to sanitized column names
- **Interactive SQL REPL**: Query your data with SQL commands
- **Export Results**: Export query results to CSV files
- **Configuration**: Environment-based configuration

## Project Structure

```
dana/
├── cmd/                # Application entry point
├── config/             # Configuration management
├── internal/           # Internal application code
│   ├── database/       # Database operations
│   ├── importer/       # File import and processing
│   ├── mapping/        # Chinese header mapping
│   └── repl/           # REPL interface
├── pkg/                # Public utilities
│   └── utils/          # String sanitization utilities
├── go.mod              # Go module dependencies
└── README.md           # This file
```

## Installation

```bash
go mod tidy
go build -o csvsql ./cmd
```

## Usage

### Basic Usage

```bash
# Load CSV files
./csvsql file1.csv file2.csv

# Load Excel files
./csvsql data.xlsx

# Load mixed file types
./csvsql users.csv resources.xlsx
```

### REPL Commands

Once the files are loaded, you'll enter an interactive SQL REPL:

- `.help` - Show available commands
- `.tables` - List available tables
- `.schema <table>` - Show table schema
- `.mappings` - Show Chinese header mappings
- `.exit` or `.quit` - Exit the application
- `EXPORT <filename.csv>` - Export last query results
- Any other input is treated as an SQL query

### Chinese Header Support

The tool automatically detects Chinese characters in column headers and:

1. Maps them to sanitized column names (e.g., `资源ID` → `_1`)
2. Maintains the mapping for query translation
3. Shows original Chinese headers in query results
4. Allows you to use Chinese field names in SQL queries

Example:
```sql
-- This query will automatically translate to use the correct column names
SELECT 资源ID, 访问地址 FROM resources WHERE 资源状态 = 'active';
```

## Configuration

Set environment variables to customize behavior:

- `DANA_DB_PATH` - Database path (default: `:memory:`)
- `DANA_MAX_FILE_SIZE` - Maximum file size in bytes (default: 100MB) NOT IMPLEMENT YET
- `DANA_VERBOSE` - Enable verbose logging (default: false) NOT IMPLEMENT YET

## Development

### Architecture

The application follows clean architecture principles:

- **Dependency Injection**: Components are injected rather than using globals
- **Interface-based Design**: Clear interfaces between components
- **Separation of Concerns**: Each module has a single responsibility
- **Testability**: Easy to test individual components

### Key Components

- **Mapper**: Handles Chinese header to column name mappings
- **Database Manager**: Manages database operations and query execution
- **File Processor**: Handles file loading and processing
- **REPL Session**: Manages the interactive session and command processing

### Testing

```bash
go test ./...
```

## Dependencies

- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/xuri/excelize/v2` - Excel file processing

## License

This project is open source and available under the MIT License.
