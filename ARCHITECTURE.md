# CF ECS Log Plugin - Project Structure

This document describes the organization of the CF ECS Log Plugin codebase.

## Directory Structure

```
cf-ecslogs-plugin/
├── bin/                    # Compiled binaries (created by build)
│   └── cf-ecslogs-plugin   # The plugin binary
├── examples/              # Example ECS log formats and usage
│   └── README.md         # Detailed examples of ECS logs
├── pkg/                   # Go packages
│   └── ecslogs/          # Main plugin package
│       ├── plugin.go     # Plugin implementation
│       ├── plugin_test.go # Unit tests
│       └── types.go      # ECS log data structures
├── .gitignore            # Git ignore rules
├── .tool-versions        # asdf tool versions
├── CONTRIBUTING.md       # Contribution guidelines
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── LICENSE               # Apache 2.0 license
├── main.go               # Plugin entry point
├── Makefile              # Build automation
├── README.md             # Main documentation
└── test-manual.sh        # Manual testing helper script
```

## Key Files

### main.go
The entry point for the plugin. Starts the CF CLI plugin interface.

### pkg/ecslogs/plugin.go
Core plugin implementation containing:
- `ECSLogsPlugin` struct implementing the CF CLI plugin interface
- `GetMetadata()` - Returns plugin metadata and command definitions
- `Run()` - Main command execution logic
- `processLogLine()` - Parses and formats individual log lines
- `parseTraditionalLogFormat()` - Handles traditional CF log format with embedded ECS JSON
- `displayECSLog()` - Formats and displays ECS logs in clear text

### pkg/ecslogs/types.go
ECS log data structures matching the Elastic Common Schema format:
- `ECSLog` - Main log entry structure
- `Labels`, `Process`, `Error`, `Service`, `Cloud`, `Host`, `ECS` - Supporting structures

### pkg/ecslogs/plugin_test.go
Unit tests for the plugin functionality:
- ECS JSON parsing tests
- Error handling tests
- Plugin metadata tests

## Build System

### Makefile Targets

- `make build` - Build the plugin binary
- `make build-all` - Build for multiple platforms
- `make test` - Run unit tests
- `make test-coverage` - Run tests with coverage report
- `make fmt` - Format code
- `make vet` - Run go vet
- `make deps` - Install/update dependencies
- `make clean` - Remove build artifacts
- `make install` - Build and install plugin to CF CLI
- `make uninstall` - Remove plugin from CF CLI
- `make reinstall` - Reinstall the plugin
- `make lint` - Run golangci-lint (if installed)
- `make help` - Show available targets

## Plugin Architecture

### Command Flow

1. User runs `cf ecslogs APP_NAME [--recent]`
2. CF CLI invokes the plugin's `Run()` method
3. Plugin parses arguments and flags
4. Plugin calls CF CLI's internal `logs` command via `CliConnection`
5. Plugin receives log output line by line
6. Each line is parsed and formatted:
   - Try to parse as pure ECS JSON
   - If that fails, check for traditional CF format with embedded JSON
   - If still not JSON, pass through as-is
7. Formatted output is printed to stdout

### Log Processing Pipeline

```
Raw Log Line
    ↓
processLogLine()
    ↓
Is it ECS JSON? ─── Yes ──→ displayECSLog()
    ↓ No                          ↓
    ↓                        Format & Print
    ↓
parseTraditionalLogFormat()
    ↓
Extract JSON from message
    ↓
displayECSLog() or Print as-is
```

## ECS Format Support

The plugin supports the Elastic Common Schema (ECS) format with the following fields:

- `@timestamp` - Log timestamp (ISO 8601)
- `message` - Log message text
- `log.level` - Log level (info, error, warn, etc.)
- `process.type` - Process type (web, worker, etc.)
- `process.index` - Process instance index
- `process.pid` - Process ID
- `error.message` - Error message
- `error.type` - Error type
- `error.stack_trace` - Stack trace
- `labels.*` - Various labels (app_name, org_name, space_name, etc.)
- `service.*` - Service information
- `cloud.*` - Cloud provider information
- `host.*` - Host information

## Output Format

The plugin formats logs as:

```
TIMESTAMP [SOURCE] LEVEL MESSAGE
```

Where:
- `TIMESTAMP` - ISO 8601 formatted timestamp
- `SOURCE` - Derived from process information (e.g., APP/PROC/WEB/0)
- `LEVEL` - Log level in uppercase
- `MESSAGE` - The log message

For errors with stack traces, additional lines are printed with indentation.

## Testing

### Unit Tests
Located in `pkg/ecslogs/plugin_test.go`, tests cover:
- ECS JSON parsing
- Error handling
- Plugin metadata

### Manual Testing
Use `test-manual.sh` to verify plugin installation and get testing instructions.

## Dependencies

- `code.cloudfoundry.org/cli` - CF CLI plugin framework
- `github.com/onsi/ginkgo` - BDD testing framework
- `github.com/onsi/gomega` - Matcher library for tests

## Development Workflow

1. Make changes to code
2. Run `make fmt` to format
3. Run `make vet` to check for issues
4. Run `make build` to compile
5. Run `make install` to install locally
6. Test with a CF environment
7. Commit changes
8. Submit pull request

## License

Apache License 2.0 - See LICENSE file for details.
