package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/vlmoon99/near-cli-go/bindata"
)

//go:embed template/**/*
var templates embed.FS

// 1.в CLI оставить сreate,build functions,deploy + call-function (from cli)
const (
	CreateCommand = "create"
	BuildCommand  = "build"
	DeployCommand = "deploy"
	CallFunction  = "call"
)

const (
	SmartContractTypeProject = "smart-contract-empty"
)

const (
	SmartContractProjectFolder = "contract"
)

const (
	ContractMainGoPath     = "template/contract/main.go.template"
	ContractMainGoFileName = "./main.go"
)

const (
	ErrProvidedNetwork                   = "(USER_INPUT_ERROR): Missing 'network'"
	ErrProvidedNetworkAndAccountName     = "(USER_INPUT_ERROR): Missing both 'network' and 'account-name'"
	ErrProvidedNetworkAndContractId      = "(USER_INPUT_ERROR): Missing both 'network' and 'contract-id'"
	ErrProvidedProjectNameModuleNameType = "(USER_INPUT_ERROR): Missing 'project-name', 'module-name', or 'type'"
	ErrIncorrectType                     = "(USER_INPUT_ERROR): Invalid project type"
	ErrRunningNearCLI                    = "(INTERNAL_UTILS): Failed to execute Near CLI"
	ErrRunningCmd                        = "(INTERNAL_UTILS): Failed to start command"
	ErrNavPrevDir                        = "(INTERNAL_UTILS): Failed to navigate to previous directory"
	ErrInitReactVite                     = "(INTERNAL_PROJECT_CLIENT_REACT): Failed to initialize React project with Vite"
	ErrGoProjectModFileIsMissing         = "(INTERNAL_PROJECT_CONTRACT): Missing 'go.mod' file"
	ErrGoProjectSumFileIsMissing         = "(INTERNAL_PROJECT_CONTRACT): Missing 'go.sum' file"
	ErrGoProjectMainGoFileIsMissing      = "(INTERNAL_PROJECT_CONTRACT): Missing 'main.go' file"
	ErrGettingCurrentDir                 = "(INTERNAL_PROJECT): Error getting current directory:"
	ErrToReadFile                        = "(INTERNAL_PROJECT): Failed to read file"
)

func WriteBinaryIfNotExists(path string, data []byte) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	return os.WriteFile(path, data, 0755)
}

func InitEmbeddedBins() {
	tempDir := os.TempDir()

	nearCliPath := filepath.Join(tempDir, "near")

	if err := WriteBinaryIfNotExists(nearCliPath, bindata.NearCli); err != nil {
		panic("failed to write near-cli: " + err.Error())
	}

}

func NearCLIWrapper(args ...string) error {
	nearCliPath := filepath.Join(os.TempDir(), "near")

	cmd := exec.Command(nearCliPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", ErrRunningNearCLI, err)
	}
	return nil
}

func TinygoRunWithRetryWrapper(command string, args []string, entityType string) {
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
	// fmt.Printf("Running command: %s %v\n", name, args)

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("Command start error: %v\n", err)

		return nil, fmt.Errorf("%s: %v", ErrRunningCmd, err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Command execution error: %v\n", err)
		fmt.Printf("Stderr: %s\n", stderr.String())
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

func CreateFolderAndNavigateThere(name string) {
	if err := os.Mkdir(name, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.Chdir(name); err != nil {
		log.Fatal(err)
	}
}

func GoBackToThePrevDirectory() {
	if err := os.Chdir(".."); err != nil {
		log.Fatalf("%s %v", ErrNavPrevDir, err)
	}
}

//Utils

// Project

func HandleCreateProject(projectName, projectType, moduleName string) {
	if projectType == SmartContractTypeProject {
		CreateFolderAndNavigateThere(projectName)
		CreateSmartContractProject(moduleName)
	} else {
		log.Fatal(ErrIncorrectType)
	}

}

func CreateSmartContractProject(moduleName string) {
	CreateFolderAndNavigateThere(SmartContractProjectFolder)
	RunCommand("go", "mod", "init", moduleName)

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		log.Fatal(ErrGoProjectModFileIsMissing)
	}

	RunCommand("go", "get", "github.com/vlmoon99/near-sdk-go@v0.0.8")

	if _, err := os.Stat("go.sum"); os.IsNotExist(err) {
		log.Fatal(ErrGoProjectSumFileIsMissing)
	}

	mainGoFileContent, err := templates.ReadFile(ContractMainGoPath)
	if err != nil {
		log.Fatalf("%s %v", ErrToReadFile, err)
	}

	WriteToFile(ContractMainGoFileName, string(mainGoFileContent))

	// CreateSmartContractIntegrationTests()
	// GoBackToThePrevDirectory()
}

// Project

// Build

func HandleBuild() {
	BuildContract()
}

func BuildContract() {
	TinygoRunWithRetryWrapper("tinygo", []string{
		"build", "-size", "short", "-no-debug", "-panic=trap",
		"-scheduler=none", "-gc=leaking", "-o", "main.wasm",
		"-target", "wasm-unknown", "./",
	}, "build")

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
	TinygoRunWithRetryWrapper("tinygo", []string{"test", "./..."}, "project")
}

func PackageTest() {
	TinygoRunWithRetryWrapper("tinygo", []string{"test", "./"}, "package")
}

//Test

// Check internal deps
func CheckDependencies(programs map[string]string) {
	missing := []string{}

	for program, helpMsg := range programs {
		if !IsInstalled(program) {
			missing = append(missing, fmt.Sprintf("%s - %s", program, helpMsg))
		}
	}

	if len(missing) > 0 {
		fmt.Println("The following required programs are missing:")
		for _, msg := range missing {
			fmt.Println(" -", msg)
		}
		ExitWithHelp()
	} else {
		fmt.Println("All necessary programs are installed.")
	}
}

func IsInstalled(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func ExitWithHelp() {
	fmt.Println("Please install the missing programs and try again.")
	exec.Command("exit", "1").Run()
}

// Check internal deps

func main() {
	InitEmbeddedBins()
	programs := map[string]string{
		"go":     "Go programming language (Install from: https://go.dev/dl/)",
		"tinygo": "TinyGo compiler for WebAssembly (Install from: https://tinygo.org/getting-started/)",
	}

	CheckDependencies(programs)
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
						Usage:    "Specify the module name for Go Smart Contract project , it can be your github for example 'https://github.com/{accountId}/{prjectName}'",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "project-type, t",
						Usage:    "Specify the type of the project, it can be 'smart-contract-empty'",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					projectName := c.String("project-name")
					projectType := c.String("project-type")
					moduleName := c.String("module-name")

					if projectName == "" || projectType == "" || moduleName == "" {
						return errors.New(ErrProvidedProjectNameModuleNameType)
					}

					HandleCreateProject(projectName, projectType, moduleName)
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
								Required: false,
							},
						},
						Action: func(c *cli.Context) error {
							network := c.String("network")
							name := c.String("account-name")

							if network == "" {
								return errors.New(ErrProvidedNetwork)
							}

							if network == "dev" && name == "" {
								return errors.New(ErrProvidedNetworkAndAccountName)

							}

							if network == "prod" {
								if err := NearCLIWrapper("account", "create-account", "fund-later", "use-auto-generation", "save-to-folder", "./"); err != nil {
									fmt.Printf("Error: %s\n", err)
									return err
								}
							} else {
								if err := NearCLIWrapper("account", "create-account", "sponsor-by-faucet-service", name, "autogenerate-new-keypair", "save-to-legacy-keychain", "network-config", "testnet", "create"); err != nil {
									fmt.Printf("Error: %s\n", err)
									return err
								}
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
						return errors.New(ErrProvidedNetworkAndContractId)

					}
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

					fmt.Println("Smart Contract deployed !")
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
