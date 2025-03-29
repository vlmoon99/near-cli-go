package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

//Constants

const (
	SmartContractTypeProject        = "smart-contract"
	FullStackTypeProjectReactNodeJs = "full-stack-react-nodejs"
)

const (
	SmartContractProjectFolder                 = "contract"
	SmartContractProjectIntegrationTestsFolder = "integration_tests"
	ClientProjectFolder                        = "client"
	BackendProjectFolder                       = "backend"
)

const (
	appJsxPath     = "../../template/client/App.jsx.template"
	appJsxFileName = "./src/App.jsx"

	blockchainDataInfoJsxPath     = "../../template/client/BlockchainDataInfo.jsx.template"
	blockchainDataInfoJsxFileName = "./src/BlockchainDataInfo.jsx"

	mainJsxPath     = "../../template/client/main.jsx.template"
	mainJsxFileName = "./src/main.jsx"

	viteConfigPath     = "../../template/client/vite.config.js.template"
	viteConfigFileName = "./vite.config.js"
)

const (
	mainGoPath        = "../../template/contract/main.go.template"
	mainGoFileName    = "./main.go"
	mainRsPath        = "../../../template/contract/main.rs.template"
	mainRsFileName    = "./src/main.rs"
	cargoTomlPath     = "../../../template/contract/Cargo.toml.template"
	cargoTomlFileName = "./Cargo.toml"
)

const (
	tsConfigJsonPath     = "../../template/backend/tsconfig.json.template"
	tsConfigJsonFileName = "./tsconfig.json"
	indexTsPath          = "../../template/backend/index.ts.template"
	indexTsFileName      = "./src/index.ts"
	gitIgnorePath        = "../../template/backend/gitignore.tempalte"
	gitIgnoreFileName    = "./.gitignore"
	dotEnvPath           = "../../template/backend/env.template"
	dotEnvFileName       = "./.env"
)

// Errors

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
)

//Utils

func NearCLIWrapper(args ...string) error {
	cmd := exec.Command("near", args...)
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
	fmt.Printf("Running command: %s %v\n", name, args)

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
		log.Fatalf("%s: %v", ErrNavPrevDir, err)
	}

	fmt.Println("Changed directory to previous.")
}

//Utils

// Project

func HandleCreateProject(projectName, projectType, moduleName string) {
	if projectType == SmartContractTypeProject {
		fmt.Println("Creating project directory...")
		CreateFolderAndNavigateThere(projectName)
		CreateSmartContractProject(moduleName)
		fmt.Println("Project created successfully!")

	} else if projectType == FullStackTypeProjectReactNodeJs {
		fmt.Println("Creating project directory...")
		CreateFolderAndNavigateThere(projectName)
		CreateSmartContractProject(moduleName)
		GoBackToThePrevDirectory()
		CreateReactClientProject()
		GoBackToThePrevDirectory()
		CreateNodeJsBackendProject()
		fmt.Println("Project created successfully!")
	} else {
		log.Fatal(ErrIncorrectType)
	}
}

func CreateSmartContractProject(moduleName string) {
	CreateFolderAndNavigateThere(SmartContractProjectFolder)

	fmt.Println("Initializing Go module...")
	RunCommand("go", "mod", "init", moduleName)

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		log.Fatal(ErrGoProjectModFileIsMissing)
	}

	fmt.Println("Installing dependencies...")
	RunCommand("go", "get", "github.com/vlmoon99/near-sdk-go@v0.0.8")

	if _, err := os.Stat("go.sum"); os.IsNotExist(err) {
		log.Fatal(ErrGoProjectSumFileIsMissing)
	}

	fmt.Println("Creating main.go file...")

	mainGoFileContent, err := ioutil.ReadFile(mainGoPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	WriteToFile(mainGoFileName, string(mainGoFileContent))

	CreateSmartContractIntegrationTests()
	GoBackToThePrevDirectory()
}

func CreateSmartContractIntegrationTests() {
	fmt.Println("Creating 'integration_tests' folder...")
	CreateFolderAndNavigateThere(SmartContractProjectIntegrationTestsFolder)

	mainRsFileContent, err := ioutil.ReadFile(mainRsPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	cargoTomlFileContent, err := ioutil.ReadFile(cargoTomlPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	fmt.Println("Initializing Cargo project...")
	RunCommand("cargo", "init", "--bin")

	fmt.Println("Writing Cargo.toml file...")
	WriteToFile(cargoTomlFileName, string(cargoTomlFileContent))

	fmt.Println("Writing boilerplate integration test code...")
	WriteToFile(mainRsFileName, string(mainRsFileContent))

	fmt.Println("Integration tests setup completed successfully!")
}

func CreateReactClientProject() {
	CreateFolderAndNavigateThere(ClientProjectFolder)

	fmt.Println("Initializing React project with Vite...")

	command := "echo 'y' | npx create-vite@latest . --template react"
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("%s: %v", ErrInitReactVite, err)
	}

	fmt.Println("Installing dependencies...")
	RunCommand("yarn", "install")

	RunCommand("yarn", "add", "near-api-js")
	RunCommand("yarn", "add", "@near-wallet-selector/core")
	RunCommand("yarn", "add", "@near-wallet-selector/modal-ui")

	RunCommand("yarn", "add", "@near-wallet-selector/my-near-wallet")
	RunCommand("yarn", "add", "@near-wallet-selector/sender")
	RunCommand("yarn", "add", "@near-wallet-selector/nearfi")
	RunCommand("yarn", "add", "@near-wallet-selector/here-wallet")
	RunCommand("yarn", "add", "@near-wallet-selector/math-wallet")
	RunCommand("yarn", "add", "@near-wallet-selector/nightly")
	RunCommand("yarn", "add", "@near-wallet-selector/meteor-wallet")
	RunCommand("yarn", "add", "@near-wallet-selector/ledger")
	RunCommand("yarn", "add", "@near-wallet-selector/wallet-connect")
	RunCommand("yarn", "add", "@near-wallet-selector/default-wallets")
	RunCommand("yarn", "add", "@near-wallet-selector/coin98-wallet")
	RunCommand("yarn", "add", "@near-wallet-selector/react-hook")
	RunCommand("yarn", "add", "--dev", "vite-plugin-node-polyfills")

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	fmt.Println("Current working directory:", dir)

	appJsxFileContent, err := ioutil.ReadFile(appJsxPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	blockchainDataInfoJsxFileContent, err := ioutil.ReadFile(blockchainDataInfoJsxPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	mainJsxFileContent, err := ioutil.ReadFile(mainJsxPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	viteConfigFileContent, err := ioutil.ReadFile(viteConfigPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	fmt.Println("Writing main.jsx file...")
	WriteToFile(viteConfigFileName, string(viteConfigFileContent))

	fmt.Println("Writing main.jsx file...")
	WriteToFile(mainJsxFileName, string(mainJsxFileContent))

	fmt.Println("Writing App.jsx file...")
	WriteToFile(appJsxFileName, string(appJsxFileContent))

	fmt.Println("Writing BlockchainDataInfo.jsx file...")
	WriteToFile(blockchainDataInfoJsxFileName, string(blockchainDataInfoJsxFileContent))

	fmt.Println("React client setup complete!")
}

func CreateNodeJsBackendProject() {
	CreateFolderAndNavigateThere(BackendProjectFolder)
	// yarn init -y
	// yarn add express cors dotenv near-api-js near-lake-framework near-seed-phrase
	// yarn add -D typescript ts-node @types/node @types/express
	RunCommand("yarn", "init", "-y")
	RunCommand("yarn", "add", "express", "cors", "dotenv", "near-api-js", "near-lake-framework", "near-seed-phrase")
	RunCommand("yarn", "add", "-D", "typescript", "ts-node", "@types/node", "@types/express")

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Current directory:", dir)

	tsConfigJsonFileContent, err := ioutil.ReadFile(tsConfigJsonPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	WriteToFile(tsConfigJsonFileName, string(tsConfigJsonFileContent))

	gitIgnoreFileContent, err := ioutil.ReadFile(gitIgnorePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	WriteToFile(gitIgnoreFileName, string(gitIgnoreFileContent))

	dotEnvFileContent, err := ioutil.ReadFile(dotEnvPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	WriteToFile(dotEnvFileName, string(dotEnvFileContent))

	err = os.Mkdir("src", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating folder:", err)
	} else {
		fmt.Println("Folder 'src' created successfully!")
	}

	indexTsFileContent, err := ioutil.ReadFile(indexTsPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	WriteToFile(indexTsFileName, string(indexTsFileContent))

	fmt.Println("Node.js server setup complete!")

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

func FullTest() {
	ProjectTest()
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
						Usage:    "Specify the module name for Go Smart Contract project , it can be your github for example 'https://github.com/{accountId}/{prjectName}'",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "project-type, t",
						Usage:    "Specify the type of the project, it can be 'smart-contract', 'full-stack'",
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
