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
	"strings"

	"github.com/urfave/cli"
	"github.com/vlmoon99/near-cli-go/bindata"
)

//go:embed template/**/*
var templates embed.FS

const (
	CreateCommand  = "create"
	TestCommand    = "test"
	BuildCommand   = "build"
	AccountCommand = "account"
	DeployCommand  = "deploy"
	CallFunction   = "call"
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
	NearSdkGoVersion = "v0.0.13"
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
	ErrBuildFailed                       = "(BUILD_ERROR): Build failed after retries"
	ErrWasmNotFound                      = "(BUILD_ERROR): WASM file not found after build"
	ErrNetworkUnreachable                = "(NETWORK_ERROR): Unable to download dependencies. Please check your internet connection"
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

func TinygoRunWithRetryWrapper(args []string, entityType string) error {
	maxRetries := 2
	var lastError error

	fmt.Printf("🔨 Building %s...\n", entityType)

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			fmt.Printf("⚡ Retrying build (attempt %d/%d)...\n", i+1, maxRetries)
		}

		cmd := exec.Command("tinygo", args...)
		output, err := cmd.CombinedOutput()

		if err == nil {
			fmt.Println(string(output))
			return nil
		}

		lastError = err

		if i == maxRetries-1 || os.Getenv("DEBUG") != "" {
			fmt.Printf("🔍 Build output:\n%s\n", string(output))
		}
	}

	return fmt.Errorf("%s: %v", ErrBuildFailed, lastError)
}

func RunCommand(name string, args ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("%s: %v", ErrRunningCmd, err)
	}

	if err := cmd.Wait(); err != nil {
		if strings.Contains(stderr.String(), "network is unreachable") ||
			strings.Contains(stderr.String(), "no route to host") ||
			strings.Contains(stderr.String(), "dial tcp") {
			return nil, fmt.Errorf("%s", ErrNetworkUnreachable)
		}
		return nil, fmt.Errorf("%v: %s", err, stderr.String())
	}

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

func HandleCreateProject(projectName, projectType, moduleName string) error {
	if projectType == SmartContractTypeProject {
		fmt.Printf("🚀 Creating new smart contract project '%s'...\n", projectName)
		CreateFolderAndNavigateThere(projectName)
		if err := CreateSmartContractProject(moduleName); err != nil {
			fmt.Printf("❌ Project creation failed: %v\n", err)
			if err := os.Chdir(".."); err == nil {
				os.RemoveAll(projectName)
			}
			return err
		}
		return nil
	}
	return fmt.Errorf("%s", ErrIncorrectType)
}

func CreateSmartContractProject(moduleName string) error {
	CreateFolderAndNavigateThere(SmartContractProjectFolder)

	fmt.Println("📦 Initializing Go module...")
	if _, err := RunCommand("go", "mod", "init", moduleName); err != nil {
		return fmt.Errorf("Failed to initialize Go module: %v", err)
	}

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("%s", ErrGoProjectModFileIsMissing)
	}

	fmt.Println("📥 Downloading dependencies...")
	if _, err := RunCommand("go", "get", fmt.Sprintf("github.com/vlmoon99/near-sdk-go@%s", NearSdkGoVersion)); err != nil {
		return err
	}

	if _, err := os.Stat("go.sum"); os.IsNotExist(err) {
		return fmt.Errorf("%s", ErrGoProjectSumFileIsMissing)
	}

	fmt.Println("📝 Creating contract template...")
	mainGoFileContent, err := templates.ReadFile(ContractMainGoPath)
	if err != nil {
		return fmt.Errorf("%s %v", ErrToReadFile, err)
	}

	WriteToFile(ContractMainGoFileName, string(mainGoFileContent))
	fmt.Println("✅ Smart contract project created successfully!")
	return nil
}

// Project

// Build

func HandleBuild() error {
	return BuildContract()
}

func BuildContract() error {
	buildArgs := []string{
		"build", "-size", "short", "-no-debug", "-o", "main.wasm",
		"-target", "wasm-unknown", "./",
	}

	if err := TinygoRunWithRetryWrapper(buildArgs, "smart contract"); err != nil {
		fmt.Printf("❌ %v\n", err)
		return err
	}

	if _, err := os.Stat("main.wasm"); os.IsNotExist(err) {
		return fmt.Errorf("%s", ErrWasmNotFound)
	}

	listCmd := exec.Command("ls", "-lh", "main.wasm")
	listOutput, err := listCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("⚠️ Warning: Could not get WASM file info: %s\n", err)
	} else {
		fmt.Printf("📦 Generated WASM file: %s", string(listOutput))
	}

	fmt.Println("✅ Build completed successfully!")
	return nil
}

// Build

// Test

func HandleTests(testType string) error {
	switch testType {
	case "project":
		return ProjectTest()
	case "package":
		return PackageTest()
	default:
		return fmt.Errorf("Invalid test type! Use 'project' or 'package'.")
	}
}

func ProjectTest() error {
	if err := TinygoRunWithRetryWrapper([]string{"test", "./..."}, "project tests"); err != nil {
		return err
	}
	fmt.Println("✅ Project tests completed successfully!")
	return nil
}

func PackageTest() error {
	if err := TinygoRunWithRetryWrapper([]string{"test", "./"}, "package tests"); err != nil {
		return err
	}
	fmt.Println("✅ Package tests completed successfully!")
	return nil
}

//Test

//Deploy

func HandleDeployContract(smartContractID string, network string) (bool, error) {
	deployCmd := []string{
		"contract", "deploy", smartContractID, "use-file", "./main.wasm",
		"without-init-call",
		"network-config", network,
		"sign-with-legacy-keychain",
		"send",
	}

	if err := NearCLIWrapper(deployCmd...); err != nil {
		fmt.Printf("Error: %s\n", err)
		return true, err
	}
	return false, nil
}

//Deploy

// Account

func HandleCreateAccount(network string, name string) (bool, error) {
	if network == "prod" {
		if err := NearCLIWrapper("account", "create-account", "fund-later", "use-auto-generation", "save-to-folder", "./"); err != nil {
			fmt.Printf("Error: %s\n", err)
			return true, err
		}
	} else {
		if err := NearCLIWrapper("account", "create-account", "sponsor-by-faucet-service", name, "autogenerate-new-keypair", "save-to-legacy-keychain", "network-config", "testnet", "create"); err != nil {
			fmt.Printf("Error: %s\n", err)
			return true, err
		}
	}
	return false, nil
}

func HandleImportAccount() (bool, error) {
	if err := NearCLIWrapper("account", "import-account"); err != nil {
		fmt.Printf("Error: %s\n", err)
		return true, err
	}
	return false, nil
}

//Account

// Call Smart Contract
func HandleCallFunction(
	signer, contract, method, args, gas, deposit, network string,
) error {
	callCmd := []string{
		"contract", "call-function",
		"as-transaction", contract, method,
		"json-args", args,
		"prepaid-gas", gas,
		"attached-deposit", deposit,
		"sign-as", signer,
		"network-config", network,
		"sign-with-keychain",
		"send",
	}

	fmt.Printf("📞 Calling function `%s` on contract `%s` from `%s`...\n", method, contract, signer)
	if err := NearCLIWrapper(callCmd...); err != nil {
		fmt.Printf("❌ Error calling function: %s\n", err)
		return err
	}

	fmt.Println("✅ Smart contract function call completed.")
	return nil
}

//Call Smart Contract

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
				Name:  CreateCommand,
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
						return fmt.Errorf("%s", ErrProvidedProjectNameModuleNameType)
					}

					if err := HandleCreateProject(projectName, projectType, moduleName); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  BuildCommand,
				Usage: "Build the project",
				Action: func(c *cli.Context) error {
					err := HandleBuild()
					if err != nil {
						fmt.Printf("❌ Build failed: %v\n", err)
						return err
					}
					return nil
				},
			},
			{
				Name:  TestCommand,
				Usage: "Run tests",
				Subcommands: []cli.Command{
					{
						Name:  "project",
						Usage: "Test the project",
						Action: func(c *cli.Context) error {
							if err := HandleTests("project"); err != nil {
								fmt.Printf("❌ Project tests failed: %v\n", err)
								return err
							}
							return nil
						},
					},
					{
						Name:  "package",
						Usage: "Test the package",
						Action: func(c *cli.Context) error {
							if err := HandleTests("package"); err != nil {
								fmt.Printf("❌ Package tests failed: %v\n", err)
								return err
							}
							return nil
						},
					},
				},
			},
			{
				Name:  AccountCommand,
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

							shouldReturn, err := HandleCreateAccount(network, name)
							if shouldReturn {
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
							shouldReturn, err := HandleImportAccount()
							if shouldReturn {
								return err
							}

							fmt.Println("Account imported successfully!")
							return nil
						},
					},
				},
			},
			{
				Name:  DeployCommand,
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
					shouldReturn, err := HandleDeployContract(smartContractID, network)
					if shouldReturn {
						return err
					}

					fmt.Println("Smart Contract deployed !")
					return nil
				},
			},
			{
				Name:  "call",
				Usage: "Call a smart contract function",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "signer, from",
						Usage:    "Account that signs the transaction",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "contract, to",
						Usage:    "Target contract account",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "method, function",
						Usage:    "Function name to call",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "args",
						Usage: "JSON-formatted arguments (optional, default: '{}')",
					},
					&cli.StringFlag{
						Name:  "gas",
						Usage: "Amount of gas to attach (default: 100 Tgas)",
					},
					&cli.StringFlag{
						Name:  "deposit",
						Usage: "Amount of NEAR to attach (default: 0 NEAR)",
					},
					&cli.StringFlag{
						Name:     "network",
						Usage:    "Specify the network",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					signer := c.String("signer")
					contract := c.String("contract")
					method := c.String("method")
					args := c.String("args")
					gas := c.String("gas")
					deposit := c.String("deposit")
					network := c.String("network")
					print(args)
					if args == "" {
						args = "{}"
					}
					if gas == "" {
						gas = "100 Tgas"
					}
					if deposit == "" {
						deposit = "0 NEAR"
					}

					return HandleCallFunction(signer, contract, method, args, gas, deposit, network)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
