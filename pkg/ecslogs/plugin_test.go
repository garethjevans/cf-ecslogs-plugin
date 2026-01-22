package ecslogs

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

func TestECSLogParsing(t *testing.T) {
	// Sample ECS JSON log with nested log object
	jsonLog := `{
		"@timestamp": "2026-01-22T17:21:36.274136327Z",
		"log": {
			"level": "DEBUG",
			"logger": "org.springframework.security.web.FilterChainProxy"
		},
		"process": {
			"pid": 8,
			"thread": {
				"name": "http-nio-8080-exec-9"
			}
		},
		"service": {
			"name": "mcp-gateway",
			"version": "0.10.2"
		},
		"message": "Securing GET /actuator/health",
		"tags": ["COMMONS-LOGGING"],
		"ecs": {
			"version": "8.11"
		}
	}`

	var log ECSLog
	err := json.Unmarshal([]byte(jsonLog), &log)
	if err != nil {
		t.Fatalf("Failed to parse ECS log: %v", err)
	}

	if log.Timestamp != "2026-01-22T17:21:36.274136327Z" {
		t.Errorf("Expected timestamp '2026-01-22T17:21:36.274136327Z', got '%s'", log.Timestamp)
	}

	if log.Message != "Securing GET /actuator/health" {
		t.Errorf("Expected message 'Securing GET /actuator/health', got '%s'", log.Message)
	}

	if log.Log == nil {
		t.Fatal("Expected log to be non-nil")
	}

	if log.Log.Level != "DEBUG" {
		t.Errorf("Expected log level 'DEBUG', got '%s'", log.Log.Level)
	}

	if log.Log.Logger != "org.springframework.security.web.FilterChainProxy" {
		t.Errorf("Expected logger 'org.springframework.security.web.FilterChainProxy', got '%s'", log.Log.Logger)
	}

	if log.Process == nil {
		t.Fatal("Expected process to be non-nil")
	}

	if log.Process.PID != 8 {
		t.Errorf("Expected process PID 8, got %d", log.Process.PID)
	}

	if log.Process.Thread == nil {
		t.Fatal("Expected thread to be non-nil")
	}

	if log.Process.Thread.Name != "http-nio-8080-exec-9" {
		t.Errorf("Expected thread name 'http-nio-8080-exec-9', got '%s'", log.Process.Thread.Name)
	}

	if len(log.Tags) != 1 || log.Tags[0] != "COMMONS-LOGGING" {
		t.Errorf("Expected tags [COMMONS-LOGGING], got %v", log.Tags)
	}
}

func TestCFLogFormatParsing(t *testing.T) {
	// Sample CF log line with ECS JSON
	logLine := `   2026-01-22T17:21:36.27+0000 [APP/REV/47/PROC/WEB/0] OUT {"@timestamp":"2026-01-22T17:21:36.274136327Z","log":{"level":"DEBUG","logger":"org.springframework.security.web.FilterChainProxy"},"process":{"pid":8,"thread":{"name":"http-nio-8080-exec-9"}},"service":{"name":"mcp-gateway","version":"0.10.2","node":{}},"message":"Securing GET /actuator/health","tags":["COMMONS-LOGGING"],"ecs":{"version":"8.11"}}`

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	plugin := &ECSLogsPlugin{}
	plugin.ProcessLine(logLine)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that output contains key elements
	if !strings.Contains(output, "APP/REV/47/PROC/WEB/0") {
		t.Errorf("Output should contain source '[APP/REV/47/PROC/WEB/0]', got: %s", output)
	}

	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Output should contain 'DEBUG', got: %s", output)
	}

	if !strings.Contains(output, "Securing GET /actuator/health") {
		t.Errorf("Output should contain message, got: %s", output)
	}

	if !strings.Contains(output, "FilterChainProxy") {
		t.Errorf("Output should contain logger name, got: %s", output)
	}

	if !strings.Contains(output, "http-nio-8080-exec-9") {
		t.Errorf("Output should contain thread name, got: %s", output)
	}
}

func TestPluginMetadata(t *testing.T) {
	plugin := &ECSLogsPlugin{}
	metadata := plugin.GetMetadata()

	if metadata.Name != "ECSLogsPlugin" {
		t.Errorf("Expected plugin name 'ECSLogsPlugin', got '%s'", metadata.Name)
	}

	if len(metadata.Commands) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(metadata.Commands))
	}

	cmd := metadata.Commands[0]
	if cmd.Name != "ecslogs" {
		t.Errorf("Expected command name 'ecslogs', got '%s'", cmd.Name)
	}

	if _, ok := cmd.UsageDetails.Options["recent"]; !ok {
		t.Error("Expected 'recent' option to be present")
	}
}

func TestErrorColorization(t *testing.T) {
	// Sample CF log line with ERROR level
	logLine := `   2026-01-22T17:21:36.28+0000 [APP/REV/47/PROC/WEB/0] OUT {"@timestamp":"2026-01-22T17:21:36.280000000Z","log":{"level":"ERROR","logger":"com.example.service.PaymentService"},"process":{"pid":8,"thread":{"name":"http-nio-8080-exec-10"}},"service":{"name":"mcp-gateway","version":"0.10.2"},"message":"Payment processing failed","ecs":{"version":"8.11"}}`

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	plugin := &ECSLogsPlugin{}
	plugin.ProcessLine(logLine)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that output contains ANSI color codes for red
	if !strings.Contains(output, "\033[31m") {
		t.Errorf("ERROR log should contain red color code (\\033[31m), got: %s", output)
	}

	if !strings.Contains(output, "\033[0m") {
		t.Errorf("ERROR log should contain reset code (\\033[0m), got: %s", output)
	}

	if !strings.Contains(output, "ERROR") {
		t.Errorf("Output should contain 'ERROR', got: %s", output)
	}
}

func TestNonErrorNoColorization(t *testing.T) {
	// Sample CF log line with INFO level
	logLine := `   2026-01-22T17:21:36.29+0000 [APP/REV/47/PROC/WEB/0] OUT {"@timestamp":"2026-01-22T17:21:36.290000000Z","log":{"level":"INFO","logger":"org.springframework.web.servlet.DispatcherServlet"},"process":{"pid":8,"thread":{"name":"http-nio-8080-exec-11"}},"service":{"name":"mcp-gateway","version":"0.10.2"},"message":"Completed request processing","ecs":{"version":"8.11"}}`

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	plugin := &ECSLogsPlugin{}
	plugin.ProcessLine(logLine)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that output contains green color code for INFO
	if !strings.Contains(output, "\033[32m") {
		t.Errorf("INFO log should contain green color code (\\033[32m), got: %s", output)
	}

	if !strings.Contains(output, "INFO") {
		t.Errorf("Output should contain 'INFO', got: %s", output)
	}
}

func TestWarnColorization(t *testing.T) {
	// Sample CF log line with WARN level
	logLine := `   2026-01-22T17:21:36.31+0000 [APP/REV/47/PROC/WEB/0] OUT {"@timestamp":"2026-01-22T17:21:36.310000000Z","log":{"level":"WARN","logger":"org.springframework.security.web.FilterChainProxy"},"process":{"pid":8,"thread":{"name":"http-nio-8080-exec-12"}},"service":{"name":"mcp-gateway","version":"0.10.2"},"message":"Authentication attempt failed","ecs":{"version":"8.11"}}`

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	plugin := &ECSLogsPlugin{}
	plugin.ProcessLine(logLine)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that output contains orange/yellow color code for WARN
	if !strings.Contains(output, "\033[33m") {
		t.Errorf("WARN log should contain orange color code (\\033[33m), got: %s", output)
	}

	if !strings.Contains(output, "WARN") {
		t.Errorf("Output should contain 'WARN', got: %s", output)
	}
}

func TestDebugColorization(t *testing.T) {
	// Sample CF log line with DEBUG level
	logLine := `   2026-01-22T17:21:36.27+0000 [APP/REV/47/PROC/WEB/0] OUT {"@timestamp":"2026-01-22T17:21:36.274136327Z","log":{"level":"DEBUG","logger":"org.springframework.security.web.FilterChainProxy"},"process":{"pid":8,"thread":{"name":"http-nio-8080-exec-9"}},"service":{"name":"mcp-gateway","version":"0.10.2"},"message":"Securing GET /actuator/health","ecs":{"version":"8.11"}}`

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	plugin := &ECSLogsPlugin{}
	plugin.ProcessLine(logLine)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that output contains blue color code for DEBUG
	if !strings.Contains(output, "\033[34m") {
		t.Errorf("DEBUG log should contain blue color code (\\033[34m), got: %s", output)
	}

	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Output should contain 'DEBUG', got: %s", output)
	}
}

func TestTraceColorization(t *testing.T) {
	// Sample CF log line with TRACE level
	logLine := `   2026-01-22T17:21:36.26+0000 [APP/REV/47/PROC/WEB/0] OUT {"@timestamp":"2026-01-22T17:21:36.260000000Z","log":{"level":"TRACE","logger":"org.springframework.web.method.HandlerMethod"},"process":{"pid":8,"thread":{"name":"http-nio-8080-exec-1"}},"service":{"name":"mcp-gateway","version":"0.10.2"},"message":"Invoking method: getUserById","ecs":{"version":"8.11"}}`

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	plugin := &ECSLogsPlugin{}
	plugin.ProcessLine(logLine)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that output contains dim grey color code for TRACE
	// Dim grey is \033[2;37m (dim + white)
	if !strings.Contains(output, "\033[2;37m") {
		t.Errorf("TRACE log should contain dim grey color code (\\033[2;37m), got: %s", output)
	}

	if !strings.Contains(output, "TRACE") {
		t.Errorf("Output should contain 'TRACE', got: %s", output)
	}
}
