package main

import (
	"bytes"
	"errors"
	"fmt"
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
	SmartContractProjectFolder = "contract"
	ClientProjectFolder        = "client"
	BackendProjectFolder       = "backend"
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
	code :=
		`
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/json"
	"github.com/vlmoon99/near-sdk-go/types"
)

//go:export InitContract
func InitContract() {
	env.LogString("Init Smart Contract")
}

//go:export WriteData
func WriteData() {
	options := types.ContractInputOptions{IsRawBytes: true}
	input, detectedType, err := env.ContractInput(options)
	if err != nil {
		env.PanicStr("Failed to get contract input: " + err.Error())
	}

	env.LogString("Contract input (JSON): " + string(input))
	env.LogString("Detected input type: " + detectedType)

	if detectedType != "object" {
		env.ContractValueReturn([]byte("Error : Incorrect type"))
	}
	parser := json.NewParser(input)

	keyResult, err := parser.GetString("key")
	if err != nil {
		env.ContractValueReturn([]byte("Error : Incorrect key"))
	}

	dataResult, err := parser.GetRawBytes("data")
	if err != nil {
		env.ContractValueReturn([]byte("Error : Incorrect data"))
	}

	env.StorageWrite([]byte(keyResult), dataResult)
	env.LogString("WriteData was successful")
}

//go:export ReadData
func ReadData() {
	options := types.ContractInputOptions{IsRawBytes: true}
	input, detectedType, err := env.ContractInput(options)
	if err != nil {
		env.PanicStr("Failed to get contract input: " + err.Error())
	}

	env.LogString("Contract input (JSON): " + string(input))
	env.LogString("Detected input type: " + detectedType)

	if detectedType != "object" {
		env.ContractValueReturn([]byte("Error : Incorrect type"))
	}
	parser := json.NewParser(input)

	keyResult, err := parser.GetString("key")
	if err != nil {
		env.ContractValueReturn([]byte("Error : Incorrect key"))
	}

	data, err := env.StorageRead([]byte(keyResult))
	if err != nil {
		env.ContractValueReturn([]byte("Error : Incorrect read from the storage by that key"))
	}
	env.LogString("ReadData was successful")

	env.ContractValueReturn(data)
}

//go:export AcceptPayment
func AcceptPayment() {
	attachedDeposit, err := env.GetAttachedDepoist()
	if err != nil {
		env.PanicStr("Failed to get attached deposit: " + err.Error())
	}
	env.LogString("Attachet Deposit :" + attachedDeposit.String())
	promiseIdx := env.PromiseBatchCreate([]byte("neargocli.testnet"))
	env.PromiseBatchActionTransfer(promiseIdx, attachedDeposit)
	env.PromiseReturn(promiseIdx)
	//neargocli.testnet
	env.LogString("AcceptPayment")
}

//go:export ReadIncommingTxData
func ReadIncommingTxData() {

	options := types.ContractInputOptions{IsRawBytes: true}
	input, detectedType, err := env.ContractInput(options)
	if err != nil {
		env.PanicStr("Failed to get contract input: " + err.Error())
	}
	env.LogString("Contract input (raw bytes): " + string(input))
	env.LogString("Detected input type: " + detectedType)

	attachedDeposit, err := env.GetAttachedDepoist()
	if err != nil {
		env.PanicStr("Failed to get attached deposit: " + err.Error())
	}
	env.LogString(fmt.Sprintf("Attached deposit: %s", attachedDeposit.String()))

	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	env.LogString("Current account ID: " + accountId)

	signerId, err := env.GetSignerAccountID()
	if err != nil || signerId == "" {
		env.PanicStr("Failed to get signer account ID: " + err.Error())
	}
	env.LogString("Signer account ID: " + signerId)

	signerPK, err := env.GetSignerAccountPK()
	if err != nil || signerPK == nil {
		env.PanicStr("Failed to get signer account PK: " + err.Error())
	}
	env.LogString("Signer account PK: " + hex.EncodeToString(signerPK))

	predecessorId, err := env.GetPredecessorAccountID()
	if err != nil || predecessorId == "" {
		env.PanicStr("Failed to get predecessor account ID: " + err.Error())
	}
	env.LogString("Predecessor account ID: " + predecessorId)

	blockHeight := env.GetCurrentBlockHeight()
	env.LogString("Current block height: " + fmt.Sprintf("%d", blockHeight))

	blockTimeMs := env.GetBlockTimeMs()
	env.LogString("Block time in ms: " + fmt.Sprintf("%d", blockTimeMs))

	epochHeight := env.GetEpochHeight()
	env.LogString("Epoch height: " + fmt.Sprintf("%d", epochHeight))

	storageUsage := env.GetStorageUsage()
	env.LogString("Storage usage: " + fmt.Sprintf("%d", storageUsage))

	accountBalance, err := env.GetAccountBalance()
	if err != nil {
		env.PanicStr("Failed to get account balance: " + err.Error())
	}
	env.LogString(fmt.Sprintf("Account balance: %s", accountBalance.String()))

	lockedBalance, err := env.GetAccountLockedBalance()
	if err != nil {
		env.PanicStr("Failed to get account locked balance: " + err.Error())
	}
	env.LogString(fmt.Sprintf("Account locked balance: %s", lockedBalance.String()))

	prepaidGas := env.GetPrepaidGas()
	env.LogString(fmt.Sprintf("Prepaid gas: %ds", prepaidGas.Inner))

	usedGas := env.GetUsedGas()
	env.LogString(fmt.Sprintf("Used gas: %d", usedGas.Inner))

	env.LogString("ReadIncommingTxData")
}

//go:export ReadBlockchainData
func ReadBlockchainData() {
	//neargocli.testnet
	env.LogString("ReadBlockchainData")
}

		`

	WriteToFile("main.go", code)

	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		log.Fatal(ErrGoProjectMainGoFileIsMissing)
	}

	CreateSmartContractIntegrationTests()
	GoBackToThePrevDirectory()
}

func CreateSmartContractIntegrationTests() {
	fmt.Println("Creating 'integration_tests' folder...")
	CreateFolderAndNavigateThere("integration_tests")

	fmt.Println("Initializing Cargo project...")
	RunCommand("cargo", "init", "--bin")

	cargoTomlContent := `
[package]
name = "integration_tests"
version = "0.1.0"
edition = "2021"

[dependencies]
anyhow = "1.0.93"
json-patch = "3.0.1"
near-workspaces = "0.15.0"
serde = "1.0.215"
serde_json = "1.0.133"
tokio = "1.41.1"
near-gas = "0.3.0"
`

	fmt.Println("Writing Cargo.toml file...")
	WriteToFile("Cargo.toml", cargoTomlContent)

	integrationTestCode := `
use near_gas::NearGas;
use near_workspaces::types::NearToken;
use serde_json::json;

async fn deploy_contract(
    worker: &near_workspaces::Worker<near_workspaces::network::Sandbox>,
) -> anyhow::Result<near_workspaces::Contract> {
    const WASM_FILEPATH: &str = "../main.wasm";
    let wasm = std::fs::read(WASM_FILEPATH)?;
    let contract = worker.dev_deploy(&wasm).await?;
    Ok(contract)
}

async fn call_integration_test_function(
    contract: &near_workspaces::Contract,
    function_name: &str,
    args: serde_json::Value,
    deposit: NearToken,
    gas: NearGas,
) -> anyhow::Result<()> {
    let outcome = contract
        .call(function_name)
        .args_json(args)
        .deposit(deposit)
        .gas(gas)
        .transact()
        .await;

    match outcome {
        Ok(result) => {
            println!("result.is_success: {:#?}", result.clone().is_success());
            println!("Functions Logs: {:#?}", result.logs());
            Ok(())
        }
        Err(err) => {
            println!(
                "{} result: Test failed with error: {:#?}",
                function_name, err
            );
            Err(err.into())
        }
    }
}


#[tokio::main]
async fn main() -> Result<()> {
    let worker = sandbox().await?;
    let contract = deploy_contract(&worker).await?;
    let standard_deposit = NearToken::from_near(3);
    let standard_gas = NearGas::from_tgas(300);
    println!("Dev Account ID: {}", contract.id());

    let success_results = vec![
        call_integration_test_function(
            &contract,
            "InitContract",
            json!({}),
            standard_deposit,
            standard_gas,
        ).await,
        call_integration_test_function(
            &contract,
            "WriteData",
            json!({ "key": "testKey", "data": "testData" }),
            standard_deposit,
            standard_gas,
        ).await,
        call_integration_test_function(
            &contract,
            "ReadData",
            json!({ "key": "testKey" }),
            standard_deposit,
            standard_gas,
        ).await,
        call_integration_test_function(
            &contract,
            "AcceptPayment",
            json!({}),
            standard_deposit,
            standard_gas,
        ).await,
        call_integration_test_function(
            &contract,
            "ReadIncommingTxData",
            json!({}),
            standard_deposit,
            standard_gas,
        ).await,
        call_integration_test_function(
            &contract,
            "ReadBlockchainData",
            json!({}),
            standard_deposit,
            standard_gas,
        ).await,
    ];

    for result in success_results {
        if let Err(e) = result {
            eprintln!("Error: {:?}", e);
        }
    }

    Ok(())
}
	`

	fmt.Println("Writing boilerplate integration test code...")
	WriteToFile("src/main.rs", integrationTestCode)

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
	RunCommand("npm", "install")

	fmt.Println("React client setup complete!")
}

func CreateNodeJsBackendProject() {
	CreateFolderAndNavigateThere(BackendProjectFolder)

	fmt.Println("Initializing Node.js project...")
	RunCommand("npm", "init", "-y")

	fmt.Println("Creating simple server file...")
	code := `
	console.log("Hello World!")
	`
	fmt.Println("React client setup complete!")
	WriteToFile("index.js", code)

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
