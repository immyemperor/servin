//go:build !linux

package cmd

import (
	"fmt"
	"os"
)

func setupContainerEnvironment() error {
	// This tool only works on Linux
	fmt.Fprintf(os.Stderr, "Error: This containerization tool only works on Linux\n")
	return fmt.Errorf("unsupported platform")
}
