package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

var projectName, moduleName string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli create -p <projectName> -m <moduleName>")
		return
	}

	if os.Args[1] == "create" {
		parseFlags()
		createProject()
	} else {
		fmt.Println("Unknown command:", os.Args[1])
	}
}

func parseFlags() {
	for i := 2; i < len(os.Args)-1; i++ {
		switch os.Args[i] {
		case "-p":
			projectName = os.Args[i+1]
		case "-m":
			moduleName = os.Args[i+1]
		}
	}

	if projectName == "" || moduleName == "" {
		log.Fatal("Error: Project name (-p) and module name (-m) are required")
	}
}

func createProject() {
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
	if err := runCommand("go", "get", "github.com/vlmoon99/near-sdk-go@v0.0.5"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Creating main.go file...")
	code := `package main

import (
	"github.com/vlmoon99/near-sdk-go/env"
)

//go:export InitContract
func InitContract() {
	env.LogString("Init Smart Contract")
}`
	if err := writeToFile("main.go", code); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Smart contract project created successfully!")
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
