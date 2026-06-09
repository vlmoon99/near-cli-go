package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func HandleCreateProject(projectName, projectType, moduleName string) error {
	if projectType != SmartContractTypeProject {
		return fmt.Errorf("%s", ErrIncorrectType)
	}

	fmt.Printf("🚀 Creating project '%s'...\n", projectName)

	projectDir, err := filepath.Abs(projectName)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	contractDir := filepath.Join(projectDir, SmartContractProjectFolder)
	if err := os.MkdirAll(contractDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	if err := initializeSmartContract(contractDir, moduleName); err != nil {
		return err
	}

	return nil
}

func initializeSmartContract(contractDir, moduleName string) error {
	fmt.Println("📝 Creating template...")
	content, err := templates.ReadFile(ContractMainGoPath)
	if err != nil {
		return fmt.Errorf("%s %v", ErrToReadFile, err)
	}

	mainGoPath := filepath.Join(contractDir, "main.go")
	if err := WriteToFile(mainGoPath, string(content)); err != nil {
		return err
	}

	fmt.Println("📦 Initializing Go module...")
	if _, err := ExecuteCommandInDir(contractDir, "go", "mod", "init", moduleName); err != nil {
		return fmt.Errorf("failed to init go module: %w", err)
	}

	goModPath := filepath.Join(contractDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("%s", ErrGoProjectModFileIsMissing)
	}

	fmt.Println("📥 Downloading dependencies...")
	if _, err := ExecuteCommandInDir(contractDir, "go", "get", fmt.Sprintf("github.com/vlmoon99/near-sdk-go@%s", NearSdkGoVersion)); err != nil {
		fmt.Printf("⚠️ Warning: Failed to download dependencies: %v\n", err)
		fmt.Println("   Please run 'go get ./...' manually inside the contract folder.")
	} else {
		goSumPath := filepath.Join(contractDir, "go.sum")
		if _, err := os.Stat(goSumPath); os.IsNotExist(err) {
			fmt.Println("⚠️ Warning: go.sum was not generated.")
		}
	}

	fmt.Println("✅ Project created successfully!")
	return nil
}
