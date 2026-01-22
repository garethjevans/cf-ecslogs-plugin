package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/garethjevans/cf-ecslogs-plugin/pkg/ecslogs"
)

func main() {
	plugin := &ecslogs.ECSLogsPlugin{}
	var scanner *bufio.Scanner

	// If a file argument is provided, read from file; otherwise read from stdin
	if len(os.Args) >= 2 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		// Read from stdin for streaming tests
		scanner = bufio.NewScanner(os.Stdin)
	}

	for scanner.Scan() {
		line := scanner.Text()
		plugin.ProcessLine(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}
