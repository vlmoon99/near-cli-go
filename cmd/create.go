package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var projectName, moduleName string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new smart contract project",
	Run: func(cmd *cobra.Command, args []string) {
		if projectName == "" || moduleName == "" {
			log.Fatal("Error: Project name and module name are required")
		}

		fmt.Println("Creating project directory...")
		if err := runCommand("mkdir", projectName); err != nil {
			log.Fatal(err)
		}
		if err := os.Chdir(projectName); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Initializing Go module...")
		if err := runCommand("go", "mod", "init", moduleName); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Installing dependencies...")
		if err := runCommand("go", "get", "github.com/vlmoon99/near-sdk-go@v0.0.2"); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Creating main.go file...")
		code := `package main

import (
	"github.com/vlmoon99/near-sdk-go/sdk"
)

//go:export InitContract
func InitContract() {
	sdk.LogString("Init Smart Contract")
}`
		if err := writeToFile("main.go", code); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Smart contract project created successfully!")
	},
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func writeToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&projectName, "project", "p", "", "Project name")
	createCmd.Flags().StringVarP(&moduleName, "module", "m", "", "Go module name")
}
