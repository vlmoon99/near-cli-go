package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// ---------------- TEST COMMANDS ---------------- //

func handleTestPackage(args []string) {
	fmt.Println("Running unit tests for the package...")

	if err := testSmartContract("tinygo", "test", "./"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Package tests complete!")
}

func handleTestProject(args []string) {
	fmt.Println("Running unit tests for the project...")

	if err := testSmartContract("tinygo", "test", "./..."); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Project tests complete!")
}

func testSmartContract(name string, args ...string) error {
	var stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil && strings.Contains(stderr.String(), "error") {
		cmd = exec.Command(name, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("test failed after retry: %v", err)
		}
	}

	fmt.Printf("Output: %v\n", stderr.String())
	return nil
}
