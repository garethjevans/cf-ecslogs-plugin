package ecslogs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/plugin"
)

// ANSI color codes
const (
	colorRed    = "\033[31m"   // ERROR
	colorOrange = "\033[33m"   // WARN (yellow/orange in terminals)
	colorGreen  = "\033[32m"   // INFO
	colorBlue   = "\033[34m"   // DEBUG
	colorGrey   = "\033[2;37m" // TRACE (dim grey)
	colorReset  = "\033[0m"    // Reset to default
)

// ECSLogsPlugin is the struct that implements the plugin interface
type ECSLogsPlugin struct{}

// GetMetadata returns the plugin metadata
func (p *ECSLogsPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "ECSLogsPlugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "ecslogs",
				HelpText: "Display ECS formatted logs in clear text",
				UsageDetails: plugin.Usage{
					Usage: "cf ecslogs APP_NAME [--recent]",
					Options: map[string]string{
						"recent": "Dump recent logs instead of tailing",
					},
				},
			},
		},
	}
}

// Run is the entry point when the plugin is invoked
func (p *ECSLogsPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] != "ecslogs" {
		return
	}

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: APP_NAME is required")
		fmt.Fprintln(os.Stderr, "Usage: cf ecslogs APP_NAME [--recent]")
		os.Exit(1)
	}

	appName := args[1]
	recent := false

	// Parse flags
	for i := 2; i < len(args); i++ {
		if args[i] == "--recent" {
			recent = true
		}
	}

	// Build the cf logs command
	var output []string
	var err error

	if recent {
		output, err = cliConnection.CliCommandWithoutTerminalOutput("logs", appName, "--recent")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving logs: %v\n", err)
			os.Exit(1)
		}
		// Process and display logs
		for _, line := range output {
			p.processLogLine(line)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Retrieving logs for app %s...\n\n", appName)
		// For streaming logs, we need to spawn cf logs command and process output line by line
		// Note: This will stream until interrupted (Ctrl+C)
		p.streamLogs(appName)
	}
}

// streamLogs spawns the cf logs command and processes output line by line in real-time
func (p *ECSLogsPlugin) streamLogs(appName string) {
	// Spawn cf logs command
	cmd := exec.Command("cf", "logs", appName)

	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stdout pipe: %v\n", err)
		os.Exit(1)
	}

	// Get stderr pipe
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stderr pipe: %v\n", err)
		os.Exit(1)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting cf logs command: %v\n", err)
		os.Exit(1)
	}

	// Process stdout in real-time
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			p.processLogLine(line)
		}
	}()

	// Process stderr in real-time (typically status messages)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			// Print stderr messages as-is (they're usually CF status messages)
			fmt.Fprintln(os.Stderr, line)
		}
	}()

	// Wait for command to finish (will run until Ctrl+C)
	if err := cmd.Wait(); err != nil {
		// Don't treat Ctrl+C as an error
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 0 && exitErr.ExitCode() != 130 { // 130 is Ctrl+C
				fmt.Fprintf(os.Stderr, "Error running cf logs: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

// ProcessLine is a public wrapper for testing that processes a single log line
func (p *ECSLogsPlugin) ProcessLine(line string) {
	p.processLogLine(line)
}

// processLogLine parses an ECS JSON log line and displays it in clear text
func (p *ECSLogsPlugin) processLogLine(line string) {
	// Trim whitespace
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// Try to parse as ECS JSON
	var ecsLog ECSLog
	if err := json.Unmarshal([]byte(line), &ecsLog); err != nil {
		// If it's not JSON, it might be a regular log line or metadata
		// Check if it looks like a CF log line format
		if strings.Contains(line, "[") && (strings.Contains(line, "APP") || strings.Contains(line, "RTR") || strings.Contains(line, "API")) {
			// Try to parse as traditional CF log format and extract JSON from message
			p.parseTraditionalLogFormat(line)
		} else {
			// Print as-is if it's not JSON (could be status messages from CF)
			fmt.Println(line)
		}
		return
	}

	// Format and display the ECS log
	p.displayECSLog(&ecsLog)
}

// parseTraditionalLogFormat handles CF logs in traditional format that may contain ECS JSON in the message
// Format: TIMESTAMP [SOURCE] OUT/ERR JSON
// Example: 2026-01-22T17:21:36.27+0000 [APP/REV/47/PROC/WEB/0] OUT {"@timestamp":"...","message":"..."}
func (p *ECSLogsPlugin) parseTraditionalLogFormat(line string) {
	// Regular expression to parse CF log format
	// Captures: timestamp, source, output type (OUT/ERR), and the rest
	// Note: There may or may not be a space before OUT/ERR
	cfLogRegex := regexp.MustCompile(`^\s*(\S+)\s+\[([^\]]+)\]\s*(OUT|ERR)\s+(.*)$`)

	matches := cfLogRegex.FindStringSubmatch(line)
	if matches == nil || len(matches) < 5 {
		// Doesn't match expected format, print as-is
		fmt.Println(line)
		return
	}

	cfTimestamp := matches[1]
	cfSource := matches[2]
	outputType := matches[3]
	content := matches[4]

	// Try to find JSON in the content
	jsonStart := strings.Index(content, "{")
	if jsonStart == -1 {
		// No JSON found, format as plain text log
		fmt.Printf("%s [%s] %s %s\n", cfTimestamp, cfSource, outputType, content)
		return
	}

	// Extract JSON part
	jsonPart := content[jsonStart:]

	var ecsLog ECSLog
	if err := json.Unmarshal([]byte(jsonPart), &ecsLog); err != nil {
		// Not valid ECS JSON, print the whole line formatted
		fmt.Printf("%s [%s] %s %s\n", cfTimestamp, cfSource, outputType, content)
		return
	}

	// Display the parsed ECS log with CF source info
	p.displayECSLogWithCFInfo(&ecsLog, cfSource, cfTimestamp)
}

// displayECSLogWithCFInfo formats and displays an ECS log with CF source information
func (p *ECSLogsPlugin) displayECSLogWithCFInfo(log *ECSLog, cfSource, cfTimestamp string) {
	timestamp := log.Timestamp
	if timestamp == "" {
		timestamp = time.Now().Format(time.RFC3339)
	}

	// Parse timestamp and format it nicely
	t, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		t, _ = time.Parse(time.RFC3339, timestamp)
	}
	formattedTime := t.Format("2006-01-02T15:04:05.000Z07:00")

	// Use the CF source from the log line (already formatted like APP/REV/47/PROC/WEB/0)
	source := cfSource

	// Get log level from the nested log object
	level := ""
	if log.Log != nil && log.Log.Level != "" {
		level = strings.ToUpper(log.Log.Level)
	} else {
		level = "INFO"
	}

	// Get the message
	message := log.Message
	if message == "" && log.Error != nil {
		message = log.Error.Message
	}

	// Add logger information if available
	loggerInfo := ""
	if log.Log != nil && log.Log.Logger != "" {
		// Extract just the class name from the full package path
		loggerParts := strings.Split(log.Log.Logger, ".")
		loggerName := loggerParts[len(loggerParts)-1]
		loggerInfo = fmt.Sprintf(" [%s]", loggerName)
	}

	// Add thread information if available
	threadInfo := ""
	if log.Process != nil && log.Process.Thread != nil && log.Process.Thread.Name != "" {
		threadInfo = fmt.Sprintf(" {%s}", log.Process.Thread.Name)
	}

	// Determine color based on log level
	var color string
	switch level {
	case "ERROR":
		color = colorRed
	case "WARN":
		color = colorOrange
	case "INFO":
		color = colorGreen
	case "DEBUG":
		color = colorBlue
	case "TRACE":
		color = colorGrey
	default:
		color = "" // No color for unknown levels
	}

	// Print with or without color
	if color != "" {
		// Apply color to the entire line
		fmt.Printf("%s%s [%s] %s%s%s %s%s\n",
			color, formattedTime, source, level, loggerInfo, threadInfo, message, colorReset)
	} else {
		// Normal output without color
		fmt.Printf("%s [%s] %s%s%s %s\n",
			formattedTime, source, level, loggerInfo, threadInfo, message)
	}

	// If there's additional context, display it (with color if applicable)
	if log.Error != nil && log.Error.StackTrace != "" {
		if color != "" {
			fmt.Printf("%s   Stack Trace:\n%s%s\n", color, log.Error.StackTrace, colorReset)
		} else {
			fmt.Printf("   Stack Trace:\n%s\n", log.Error.StackTrace)
		}
	}
}

// displayECSLog formats and displays an ECS log in clear text (without CF source info)
func (p *ECSLogsPlugin) displayECSLog(log *ECSLog) {
	p.displayECSLogWithCFInfo(log, "APP", "")
}
