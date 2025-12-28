package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func HandleBuild(sourceDir, outputName string, keepGenerated bool) error {
	absSourceDir, err := filepath.Abs(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of source: %w", err)
	}

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

	tmpFileName := "generated_build.go"
	tmpFilePath := filepath.Join(absSourceDir, tmpFileName)

	fmt.Println("📝 Writing intermediate build file...")
	if err := WriteToFile(tmpFilePath, generatedCode); err != nil {
		return fmt.Errorf("failed to write generated file '%s': %w", tmpFilePath, err)
	}

	defer func() {
		if !keepGenerated {
			fmt.Printf("🧹 Cleaning up temporary file: %s\n", tmpFilePath)
			if err := os.Remove(tmpFilePath); err != nil {
				fmt.Printf("⚠️ Warning: Failed to clean up temporary file '%s': %v\n", tmpFilePath, err)
			}
		} else {
			fmt.Printf("💾 Kept generated file: %s\n", tmpFilePath)
		}
	}()

	args := []string{
		"build", "-size", "short", "-no-debug",
		"-o", absOutputName,
		"-target", "wasm-unknown",
		tmpFileName,
	}

	fmt.Printf("🔨 Compiling to %s...\n", outputName)

	if err := ExecuteWithRetry(GetTinyGoPath(), args, absSourceDir, 2, os.Getenv("DEBUG") != ""); err != nil {
		fmt.Printf("DEBUG: TinyGo compilation failed: %v\n", err)
		return err
	}

	if _, err := os.Stat(absOutputName); os.IsNotExist(err) {
		return fmt.Errorf("%s: output file '%s' not found after build", ErrWasmNotFound, absOutputName)
	}

	fmt.Printf("✅ Build completed successfully: %s\n", outputName)
	return nil
}

func HandleTests(testType string) error {
	target := "./..."
	if testType == "package" {
		target = "./"
	} else if testType != "project" {
		return fmt.Errorf("invalid test type provided: '%s'. Use 'project' or 'package'.", testType)
	}

	fmt.Printf("🧪 Running %s tests...\n", testType)

	if err := ExecuteWithRetry(GetTinyGoPath(), append([]string{"test"}, target), "", 2, true); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	fmt.Println("✅ Tests passed!")
	return nil
}
