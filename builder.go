package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// HandleBuild orchestrates the build process
func HandleBuild(sourceDir, outputName string) error {
	// 1. Resolve Absolute Paths to avoid confusion when changing directories
	absSourceDir, err := filepath.Abs(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of source: %w", err)
	}

	// Default output name handling
	if outputName == "" {
		outputName = "main.wasm"
	}
	if !strings.HasSuffix(outputName, ".wasm") {
		outputName += ".wasm"
	}

	absOutputName, err := filepath.Abs(outputName)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of output: %w", err)
	}

	fmt.Printf("DEBUG: HandleBuild context\n  Source: %s\n  Output: %s\n", absSourceDir, absOutputName)

	fmt.Printf("🔍 Scanning project in: %s\n", absSourceDir)
	generatedCode, err := GenerateCode(absSourceDir)
	if err != nil {
		fmt.Printf("DEBUG: GenerateCode returned error: %v\n", err)
		return fmt.Errorf("code generation failed: %w", err)
	}

	// The temporary file will be created INSIDE the source directory
	tmpFileName := "generated_build.go"
	tmpFilePath := filepath.Join(absSourceDir, tmpFileName)

	fmt.Println("📝 Writing intermediate build file...")
	if err := WriteToFile(tmpFilePath, generatedCode); err != nil {
		return fmt.Errorf("failed to write generated file '%s': %w", tmpFilePath, err)
	}

	// // --- Temporary File Cleanup ---
	// // COMMENT OUT this block to keep the generated file for debugging/analysis
	// defer func() {
	// 	fmt.Printf("🧹 Cleaning up temporary file: %s\n", tmpFilePath)
	// 	if err := os.Remove(tmpFilePath); err != nil {
	// 		fmt.Printf("⚠️ Warning: Failed to clean up temporary file '%s': %v\n", tmpFilePath, err)
	// 	}
	// }()
	// // ------------------------------

	// --- TinyGo Compilation ---
	args := []string{
		"build", "-size", "short", "-no-debug",
		"-o", absOutputName, // Use absolute path for output so it writes to the correct place
		"-target", "wasm-unknown",
		tmpFileName, // Just the filename, because we will run command INSIDE absSourceDir
	}

	fmt.Printf("🔨 Compiling to %s...\n", outputName)

	// Execute TinyGo inside the contract directory
	// passing absSourceDir as the 'dir' argument
	if err := ExecuteWithRetry("tinygo", args, absSourceDir, 2, os.Getenv("DEBUG") != ""); err != nil {
		fmt.Printf("DEBUG: TinyGo compilation failed: %v\n", err)
		return err
	}

	// --- Verification ---
	if _, err := os.Stat(absOutputName); os.IsNotExist(err) {
		return fmt.Errorf("%s: output file '%s' not found after build", ErrWasmNotFound, absOutputName)
	}

	fmt.Printf("✅ Build completed successfully: %s\n", outputName)
	return nil
}

// HandleTests runs the tests for the smart contract.
func HandleTests(testType string) error {
	target := "./..."
	if testType == "package" {
		target = "./"
	} else if testType != "project" {
		return fmt.Errorf("invalid test type provided: '%s'. Use 'project' or 'package'.", testType)
	}

	fmt.Printf("🧪 Running %s tests...\n", testType)

	// Pass "" as dir to run in current working directory
	if err := ExecuteWithRetry("tinygo", append([]string{"test"}, target), "", 2, true); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	fmt.Println("✅ Tests passed!")
	return nil
}
