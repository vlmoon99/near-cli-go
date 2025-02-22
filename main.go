package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CommandHandler func(args []string)

var commands = map[string]CommandHandler{
	"create":             handleCreate,
	"build":              handleBuild,
	"deploy":             handleDeploy,
	"create-dev-account": handleCreateDevAccount,
	// "import-mainnet-account": handleImportMainnetAccount,
	"test-package": handleTestPackage,
	"test-project": handleTestProject,
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
	// fmt.Println("  cli import-mainnet-account")
	fmt.Println("  cli test-package")
	fmt.Println("  cli test-project")
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
	runCommand("go", "get", "github.com/vlmoon99/near-sdk-go@v0.0.8")

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

	fmt.Println("Smart contract project created successfully!")
}

// ---------------- BUILD COMMAND ---------------- //

func handleBuild(args []string) {
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		log.Fatal("Error: Cannot compile. main.go is missing.")
	}

	fmt.Println("Building smart contract...")

	if err := buildSmartContract(); err != nil {
		log.Fatal(err)
	}

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

	if err := buildSmartContract(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Verifying build...")
	runCommand("ls", "-lh", "main.wasm")

	network := "testnet"
	if isProd {
		network = "mainnet"
	}

	fmt.Println("Deploying contract...")
	deployCmd := []string{
		"contract", "deploy", smartContractID, "use-file", "./main.wasm",
		"without-init-call",
		"network-config", network,
		"sign-with-legacy-keychain",
		"send",
	}
	output, err := runCommand("near", deployCmd...)
	if err != nil {
		log.Fatalf("Error deploying contract: %v\n%s", err, string(output))
	}

	fmt.Println("Deployment complete!")
	fmt.Printf("Deploy command output: %v", string(output))
}

// ---------------- CREATE DEV ACCOUNT ---------------- //

func handleCreateDevAccount(args []string) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter your desired account ID (without the .testnet postfix): ")
	scanner.Scan()
	accountID := strings.TrimSpace(scanner.Text())

	accountID = accountID + ".testnet"

	fmt.Println("Creating developer account...")

	createCmd := []string{
		"account", "create-account", "sponsor-by-faucet-service", accountID,
		"autogenerate-new-keypair", "save-to-legacy-keychain", "network-config", "testnet", "create",
	}

	runCommand("near", createCmd...)

	fmt.Println("Developer account created successfully!")
}

// // ---------------- IMPORT MAINNET ACCOUNT ---------------- //

// func handleImportMainnetAccount(args []string) {
// 	scanner := bufio.NewScanner(os.Stdin)

// 	fmt.Print("Enter your seed phrase (12 words): ")
// 	scanner.Scan()
// 	seedPhrase := strings.TrimSpace(scanner.Text())

// 	fmt.Print("Enter your account ID: ")
// 	scanner.Scan()
// 	accountID := strings.TrimSpace(scanner.Text())

// 	fmt.Println("Importing mainnet account...")
// 	importCmd := []string{
// 		"account", "import-account", "using-seed-phrase", seedPhrase,
// 		"--seed-phrase-hd-path", "m/44'/397'/0'",
// 		"network-config", "mainnet",
// 		accountID,
// 	}
// 	runCommand("near", importCmd...)

// 	fmt.Println("Mainnet account imported successfully!")
// }

// ---------------- TEST COMMANDS ---------------- //

func handleTestPackage(args []string) {
	fmt.Println("Running unit tests for the package...")

	if err := testSmartContract("tinygo", "test", "./"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Package tests complete!")
}

func handleTestProject(args []string) {
	fmt.Println("Running unit tests for the project...")

	if err := testSmartContract("tinygo", "test", "./..."); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Project tests complete!")
}

// ---------------- UTILITY FUNCTIONS ---------------- //

func testSmartContract(name string, args ...string) error {
	var stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil && strings.Contains(stderr.String(), "error") {
		cmd = exec.Command(name, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("test failed after retry: %v", err)
		}
	}

	fmt.Printf("Output: %v\n", stderr.String())
	return nil
}

func buildSmartContract() error {
	buildCmd := []string{
		"build", "-size", "short", "-no-debug", "-panic=trap",
		"-scheduler=none", "-gc=leaking", "-o", "main.wasm", "-target", "wasm-unknown", "./",
	}

	output, err := runCommand("tinygo", buildCmd...)
	if err != nil && strings.Contains(string(output), "unsupported parameter type") {
		output, err = runCommand("tinygo", buildCmd...)
		if err != nil {
			return fmt.Errorf("build failed after retry: %v", err)
		}

		fmt.Printf("Build Output : %v", string(output))

	}

	fmt.Printf("Build Output : %v", string(output))

	return nil
}

func runCommand(name string, args ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%v: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
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
