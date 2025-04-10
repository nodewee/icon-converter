// Package main provides the entry point for the icon-converter application.
package main

import (
	"icon-converter/cmd"
	"os"
)

// main is the entry point of the application.
// It executes the root command and exits with code 1 if an error occurs.
func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
