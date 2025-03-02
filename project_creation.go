package main

import (
	"fmt"
	"log"
	"os"
)

func handleCreate(args []string) {
	projectName, moduleName := parseCreateFlags(args)
	createProject(projectName, moduleName)
}

func parseCreateFlags(args []string) (string, string) {
	var (
		projectName string
		moduleName  string
	)

	for i := 0; i < len(args)-1; i++ {
		switch args[i] {
		case "-p":
			projectName = args[i+1]
		case "-m":
			moduleName = args[i+1]
		}
	}

	if projectName == "" || moduleName == "" {
		log.Fatal("Error: Project name (-p) and module name (-m) are required")
	}
	return projectName, moduleName
}

func createProject(projectName, moduleName string) {
	fmt.Println("Creating project directory...")
	if err := os.Mkdir(projectName, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.Chdir(projectName); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initializing Go module...")
	runCommand("go", "mod", "init", moduleName)

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		log.Fatal("Error: go.mod is missing.")
	}

	fmt.Println("Installing dependencies...")
	runCommand("go", "get", "github.com/vlmoon99/near-sdk-go@v0.0.8")

	if _, err := os.Stat("go.sum"); os.IsNotExist(err) {
		log.Fatal("Error: Cannot compile. go.sum is missing.")
	}

	fmt.Println("Creating main.go file...")
	code :=
		`package main

	import (
		"github.com/vlmoon99/near-sdk-go/env"
	)

	//go:export InitContract
	func InitContract() {
		env.LogString("Init Smart Contract")
	}`

	writeToFile("main.go", code)

	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		log.Fatal("Error: Cannot compile. main.go is missing.")
	}

	fmt.Println("Project created successfully!")
}
