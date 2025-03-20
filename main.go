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

func RunWithRetry(command string, args []string, entityType string) {
	fmt.Printf("Running tests for the %s...\n", entityType)

	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("First %s test attempt failed: %s\n", entityType, string(output))

		fmt.Printf("Retrying %s tests...\n", entityType)
		output, err = cmd.CombinedOutput()
		fmt.Println(string(output))

		if err != nil {
			fmt.Printf("Second %s test attempt failed: %s\n", entityType, string(output))
			return
		}
	}

	fmt.Println(string(output))
	fmt.Printf("%s tests completed successfully!\n", entityType)
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

func HandleCreateProject(projectName, moduleName string) {
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

func HandleBuild() {
	BuildContract()
}
func BuildContract() {
	RunWithRetry("tinygo", []string{
		"build", "-size", "short", "-no-debug", "-panic=trap",
		"-scheduler=none", "-gc=leaking", "-o", "main.wasm",
		"-target", "wasm-unknown", "./",
	}, "build")

	// After successful build, list the output file
	listCmd := exec.Command("ls", "-lh", "main.wasm")
	listOutput, err := listCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error listing the file: %s\n", string(listOutput))
		return
	}

	fmt.Println(string(listOutput))
	fmt.Println("Project build completed!")
}

// Build

// Test

func HandleTests(testType string) {
	switch testType {
	case "project":
		ProjectTest()
	case "package":
		PackageTest()
	default:
		fmt.Println("Invalid test type! Use 'project' or 'package'.")
	}
}

func ProjectTest() {
	RunWithRetry("tinygo", []string{"test", "./..."}, "project")
}

func PackageTest() {
	RunWithRetry("tinygo", []string{"test", "./"}, "package")
}

func FullTest() {
	ProjectTest()
	//Integration tests, etc
}

//Test

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

					HandleCreateProject(projectName, moduleName)
					return nil
				},
			},
			{
				Name:  "build",
				Usage: "Build the project",
				Action: func(c *cli.Context) error {
					HandleBuild()
					return nil
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
							HandleTests("project")
							return nil
						},
					},
					{
						Name:  "package",
						Usage: "Test the package",
						Action: func(c *cli.Context) error {
							HandleTests("package")
							fmt.Println("Package test completed!")
							return nil
						},
					},
					// {
					// 	Name:  "integration",
					// 	Usage: "Integration Tests",
					// 	Action: func(c *cli.Context) error {
					// 		fmt.Println("Package test completed!")
					// 		return nil
					// 	},
					// },
				},
			},
			{
				Name:  "account",
				Usage: "Manage blockchain accounts",
				Subcommands: []cli.Command{
					{
						Name:  "create",
						Usage: "Create a new account on the Near Blockchain",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "network, n",
								Usage:    "Specify the netowrk of the account",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "account-name, a",
								Usage:    "Specify the account name",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							network := c.String("network")
							name := c.String("account-name")

							if network == "" || name == "" {
								return fmt.Errorf("Both 'network' and 'account-name' must be provided.")
							}

							if err := NearCLIWrapper(fmt.Sprintf("account create-account sponsor-by-faucet-service %s autogenerate-new-keypair save-to-legacy-keychain network-config %s create", name, network)); err != nil {
								fmt.Printf("Error: %s\n", err)
								return err
							}

							fmt.Println("Development account created successfully!")
							return nil
						},
					},
					{
						Name:  "import",
						Usage: "Import an account",
						Action: func(c *cli.Context) error {
							if err := NearCLIWrapper("account", "import-account"); err != nil {
								fmt.Printf("Error: %s\n", err)
								return err
							}

							fmt.Println("Account imported successfully!")
							return nil
						},
					},
				},
			},
			{
				Name:  "deploy",
				Usage: "Deploy the project to production",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "contract-id, id",
						Usage:    "Specify the smart contract id",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "network, n",
						Usage:    "Specify the netowrk of the account",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					smartContractID := c.String("contract-id")
					network := c.String("network")

					if network == "" || smartContractID == "" {
						return fmt.Errorf("Both 'network' and 'contract-id' must be provided.")
					}
					//1.Contract dev,prod
					deployCmd := []string{
						"contract", "deploy", smartContractID, "use-file", "./main.wasm",
						"without-init-call",
						"network-config", network,
						"sign-with-legacy-keychain",
						"send",
					}

					if err := NearCLIWrapper(deployCmd...); err != nil {
						fmt.Printf("Error: %s\n", err)
						return err
					}

					fmt.Println("Project deployed to production!")
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
