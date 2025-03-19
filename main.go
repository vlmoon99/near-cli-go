package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

//Utils

func NearCLIWrapper(args ...string) error {
	cmd := exec.Command("near", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running Near CLI: %w", err)
	}
	return nil
}

func RunCommand(name string, args ...string) ([]byte, error) {
	fmt.Printf("Running command: %s %v\n", name, args) // Log the command being executed

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("Command start error: %v\n", err)
		return nil, fmt.Errorf("start error: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Command execution error: %v\n", err)
		fmt.Printf("Stderr: %s\n", stderr.String()) // Log stderr for debugging
		return nil, fmt.Errorf("%v: %s", err, stderr.String())
	}

	//TODO : Getting know why near cli rs gives us error
	if stdout.String() == "" {
		return stderr.Bytes(), nil
	}

	return stdout.Bytes(), nil
}

func WriteToFile(filename, content string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
}

//Utils

// Project

func CreateProject(projectName, moduleName string) {
	fmt.Println("Creating project directory...")
	if err := os.Mkdir(projectName, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.Chdir(projectName); err != nil {
		log.Fatal(err)
	}

	CreateSmartContractProject(moduleName)

	fmt.Println("Project created successfully!")
}

func CreateSmartContractProject(moduleName string) {
	fmt.Println("Initializing Go module...")
	RunCommand("go", "mod", "init", moduleName)

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		log.Fatal("Error: go.mod is missing.")
	}

	fmt.Println("Installing dependencies...")
	RunCommand("go", "get", "github.com/vlmoon99/near-sdk-go@v0.0.8")

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

	WriteToFile("main.go", code)

	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		log.Fatal("Error: Cannot compile. main.go is missing.")
	}
}

// Project

// Build
func BuildContract() (bool, error) {
	cmd := exec.Command("tinygo", "build", "-size", "short", "-no-debug", "-panic=trap", "-scheduler=none", "-gc=leaking", "-o", "main.wasm", "-target", "wasm-unknown", "./")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("First build attempt failed: %s\n", string(output))

		fmt.Println("Retrying build...")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Second build attempt failed: %s\n", string(output))
			return true, err
		}
	}

	fmt.Println("Project build completed!")

	listCmd := exec.Command("ls", "-lh", "main.wasm")
	listOutput, err := listCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error listing the file: %s\n", string(listOutput))
		return true, err
	}

	fmt.Println(string(listOutput))
	return false, nil
}

// Build
func main() {
	app := &cli.App{
		Name:  "near-go",
		Usage: "CLI tool for managing projects on Near Blockchain",
		Commands: []cli.Command{
			{
				Name:  "create",
				Usage: "Create a new project",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "project-name, p",
						Usage:    "Specify the name of the project",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "module-name, m",
						Usage:    "Specify the module name for Go project",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					projectName := c.String("project-name")
					moduleName := c.String("module-name")

					if projectName == "" || moduleName == "" {
						return fmt.Errorf("both project-name and module-name must be provided")
					}

					CreateProject(projectName, moduleName)
					return nil
				},
			},
			{
				Name:  "build",
				Usage: "Build the project",
				Action: func(c *cli.Context) error {
					shouldReturn, err := BuildContract()
					if shouldReturn {
						return err
					}
					return nil
				},
			},
			{
				Name:  "deploy",
				Usage: "Deploy the project to production",
				Action: func(c *cli.Context) error {
					//1.Contract dev,prod
					//2.client ?
					fmt.Println("Project deployed to production!")
					return nil
				},
			},
			{
				Name:  "account",
				Usage: "Manage blockchain accounts",
				Subcommands: []cli.Command{
					{
						Name:  "create",
						Usage: "Create a new account for development",
						Action: func(c *cli.Context) error {
							//1.dev prod
							if err := NearCLIWrapper("account"); err != nil {
								fmt.Printf("Error: %s\n", err)
							}
							fmt.Println("Dev account created!")
							return nil
						},
					},
					{
						Name:  "import",
						Usage: "Import an account",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "type",
								Value: "dev",
								Usage: "Specify account type (e.g., dev, prod)",
							},
						},
						Action: func(c *cli.Context) error {
							//1. dev prod
							fmt.Printf("Account imported for: %s\n", c.String("type"))
							return nil
						},
					},
				},
			},
			{
				Name:  "test",
				Usage: "Run tests",
				Subcommands: []cli.Command{
					{
						Name:  "project",
						Usage: "Test the project",
						Action: func(c *cli.Context) error {
							fmt.Println("Project test completed!")
							return nil
						},
					},
					{
						Name:  "package",
						Usage: "Test the package",
						Action: func(c *cli.Context) error {
							fmt.Println("Package test completed!")
							return nil
						},
					},
					{
						Name:  "integration",
						Usage: "Integration Tests",
						Action: func(c *cli.Context) error {
							fmt.Println("Package test completed!")
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
