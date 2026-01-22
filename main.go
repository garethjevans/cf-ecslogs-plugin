package main

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/garethjevans/cf-ecslogs-plugin/pkg/ecslogs"
)

func main() {
	plugin.Start(&ecslogs.ECSLogsPlugin{})
}

func init() {
	// Handle any panics gracefully
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", r)
			os.Exit(1)
		}
	}()
}
