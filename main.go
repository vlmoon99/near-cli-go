package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CommandHandler func(args []string)

var commands = map[string]CommandHandler{
	"create":                 handleCreate,
	"build":                  handleBuild,
	"deploy":                 handleDeploy,
	"create-dev-account":     handleCreateDevAccount,
	"import-mainnet-account": handleImportMainnetAccount,
}

var (
	projectName string
	moduleName  string
)

func main() {
	if !checkTinyGo() {
		fmt.Println("TinyGo is not installed. Please visit this link for installation instructions: https://tinygo.org/getting-started/install/")
	}

	if !checkNearRsCli() {
		fmt.Println("NEAR CLI RS is not installed. Please visit this link for installation instructions: https://github.com/near/near-cli-rs")
	}

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	handler, exists := commands[command]
	if !exists {
		fmt.Println("Unknown command:", command)
		printUsage()
		return
	}

	handler(os.Args[2:])
}

// Print available commands
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  cli create -p <projectName> -m <moduleName>")
	fmt.Println("  cli build")
	fmt.Println("  cli deploy [--prod]")
	fmt.Println("  cli create-dev-account")
	fmt.Println("  cli import-mainnet-account")
}

// ---------------- CREATE COMMAND ---------------- //

func handleCreate(args []string) {
	parseCreateFlags(args)
	createProject()
}

func parseCreateFlags(args []string) {
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
}

func createProject() {
	fmt.Println("Creating project directory...")
	if err := os.Mkdir(projectName, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.Chdir(projectName); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initializing Go module...")
	runCommand("go", "mod", "init", moduleName)

	fmt.Println("Installing dependencies...")
	runCommand("go", "get", "github.com/vlmoon99/near-sdk-go@v0.0.5")

	fmt.Println("Creating main.go file...")
	code := `package main

import (
	"github.com/vlmoon99/near-sdk-go/env"
)

//go:export InitContract
func InitContract() {
	env.LogString("Init Smart Contract")
}`
	writeToFile("main.go", code)

	fmt.Println("Smart contract project created successfully!")
}

// ---------------- BUILD COMMAND ---------------- //

func handleBuild(args []string) {

	// Check if main.go exists
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		log.Fatal("Error: Cannot compile. main.go is missing.")
	}

	fmt.Println("Building smart contract...")

	// Build command for TinyGo
	buildCmd := []string{
		"build", "-size", "short", "-no-debug", "-panic=trap",
		"-scheduler=none", "-gc=leaking", "-o", "main.wasm", "-target", "wasm-unknown", "./",
	}

	runCommand("tinygo", buildCmd...)

	fmt.Println("Build complete! Generated main.wasm")
}

// ---------------- DEPLOY COMMAND ---------------- //

func handleDeploy(args []string) {
	var smartContractID string
	isProd := len(args) > 0 && args[0] == "--prod"
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter your smart contract account ID: ")
	scanner.Scan()
	smartContractID = strings.TrimSpace(scanner.Text())

	fmt.Println("Building for deployment...")
	buildCmd := []string{
		"build", "-size", "short", "-no-debug", "-panic=trap",
		"-scheduler=none", "-gc=leaking", "-o", "main.wasm", "-target", "wasm-unknown", "./",
	}
	runCommand("tinygo", buildCmd...)

	fmt.Println("Verifying build...")
	runCommand("ls", "-lh", "main.wasm")

	network := "testnet"
	if isProd {
		network = "mainnet"
	}

	fmt.Println("Deploying contract...")
	deployCmd := []string{
		"contract", "deploy", smartContractID, "use-file", "./main.wasm",
		"with-init-call", "InitContract", "json-args", "{}",
		"prepaid-gas", "100.0 Tgas", "attached-deposit", "0 NEAR",
		"network-config", network, "sign-with-legacy-keychain", "send",
	}
	runCommand("near", deployCmd...)

	fmt.Println("Deployment complete!")
}

// ---------------- CREATE DEV ACCOUNT ---------------- //

func handleCreateDevAccount(args []string) {
	scanner := bufio.NewScanner(os.Stdin)

	// Prompt the user for the account ID
	fmt.Print("Enter your desired account ID (without the .testnet postfix): ")
	scanner.Scan()
	accountID := strings.TrimSpace(scanner.Text())

	// Append the .testnet postfix to the account ID
	accountID = accountID + ".testnet"

	fmt.Println("Creating developer account...")

	createCmd := []string{
		"account", "create-account", "sponsor-by-faucet-service", accountID,
		"autogenerate-new-keypair", "save-to-legacy-keychain", "network-config", "testnet", "create",
	}

	// Run the command
	runCommand("near", createCmd...)

	fmt.Println("Developer account created successfully!")
}

// ---------------- IMPORT MAINNET ACCOUNT ---------------- //

func handleImportMainnetAccount(args []string) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter your seed phrase (12 words): ")
	scanner.Scan()
	seedPhrase := strings.TrimSpace(scanner.Text())

	fmt.Println("Importing mainnet account...")
	importCmd := []string{
		"account", "import-account", "using-seed-phrase", seedPhrase,
		"--seed-phrase-hd-path", "m/44'/397'/0'",
		"network-config", "mainnet",
	}
	runCommand("near", importCmd...)

	fmt.Println("Mainnet account imported successfully!")
}

// ---------------- UTILITY FUNCTIONS ---------------- //

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running command: %s %v\n%v", name, args, err)
	}
}

func writeToFile(filename, content string) {
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

func checkTinyGo() bool {
	cmd := exec.Command("tinygo", "version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running tinygo:", err)
		return false
	}
	fmt.Println("TinyGo Version:", string(output))
	return true
}

func checkNearRsCli() bool {
	cmd := exec.Command("near", "--version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running near:", err)
		return false
	}
	fmt.Println("Near CLI RS Version:", string(output))
	return true
}
