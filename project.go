package main

import (
	"fmt"
	"os"
)

func HandleCreateProject(projectName, projectType, moduleName string) error {
	if projectType != SmartContractTypeProject {
		return fmt.Errorf("%s", ErrIncorrectType)
	}

	fmt.Printf("🚀 Creating project '%s'...\n", projectName)
	if err := CreateFolderAndNavigate(projectName); err != nil {
		return err
	}

	if err := initializeSmartContract(moduleName); err != nil {
		os.Chdir("..")
		os.RemoveAll(projectName)
		return err
	}

	return nil
}

func initializeSmartContract(moduleName string) error {
	if err := CreateFolderAndNavigate(SmartContractProjectFolder); err != nil {
		return err
	}

	fmt.Println("📦 Initializing Go module...")
	if _, err := ExecuteCommand("go", "mod", "init", moduleName); err != nil {
		return err
	}

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("%s", ErrGoProjectModFileIsMissing)
	}

	fmt.Println("📥 Downloading dependencies...")
	if _, err := ExecuteCommand("go", "get", fmt.Sprintf("github.com/vlmoon99/near-sdk-go@%s", NearSdkGoVersion)); err != nil {
		return err
	}

	if _, err := os.Stat("go.sum"); os.IsNotExist(err) {
		return fmt.Errorf("%s", ErrGoProjectSumFileIsMissing)
	}

	fmt.Println("📝 Creating template...")
	content, err := templates.ReadFile(ContractMainGoPath)
	if err != nil {
		return fmt.Errorf("%s %v", ErrToReadFile, err)
	}

	if err := WriteToFile(ContractMainGoFileName, string(content)); err != nil {
		return err
	}

	fmt.Println("✅ Project created successfully!")
	return nil
}