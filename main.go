package main

import (
	"log"
	"os"

	"github.com/ahmedYasserM/eval/cmd"
)

func main() {
	// Configure logger for the CLI
	logger := log.New(os.Stderr, "[eval] ", log.LstdFlags|log.Lshortfile)

	if err := cmd.Execute(); err != nil {
		logger.Fatalf("Root Command Faile: %v", err)
	}

}
