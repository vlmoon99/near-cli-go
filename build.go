package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func handleBuild(args []string) {
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		log.Fatal("Error: Cannot compile. main.go is missing.")
	}

	fmt.Println("Building smart contract...")

	if err := buildSmartContract(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Build complete! Generated main.wasm")
}

func buildSmartContract() error {
	buildCmd := []string{
		"build", "-size", "short", "-no-debug", "-panic=trap",
		"-scheduler=none", "-gc=leaking", "-o", "main.wasm", "-target", "wasm-unknown", "./",
	}

	output, err := runCommand("tinygo", buildCmd...)
	if err != nil && strings.Contains(string(output), "unsupported parameter type") {
		output, err = runCommand("tinygo", buildCmd...)
		if err != nil || strings.Contains(string(output), "error") || strings.Contains(string(output), "Error") {
			return fmt.Errorf("build failed after retry: %v", err)
		}

		fmt.Printf("Build Output : %v", string(output))

	}

	fmt.Printf("Build Output : %v", string(output))

	return nil
}
