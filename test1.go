package main

// import (
// 	"fmt"
// 	"os"
// )

// type CommandHandler func(args []string)

// var commands = map[string]CommandHandler{
// 	"create":             handleCreate,
// 	"build":              handleBuild,
// 	"deploy":             handleDeploy,
// 	"create-dev-account": handleCreateDevAccount,
// 	// "import-mainnet-account": handleImportMainnetAccount,
// 	"test-package": handleTestPackage,
// 	"test-project": handleTestProject,
// }

// func main() {
// 	if !checkTinyGo() {
// 		fmt.Println("TinyGo is not installed. Please visit this link for installation instructions: https://tinygo.org/getting-started/install/")
// 	}

// 	if !checkNearRsCli() {
// 		fmt.Println("NEAR CLI RS is not installed. Please visit this link for installation instructions: https://github.com/near/near-cli-rs")
// 	}

// 	if len(os.Args) < 2 {
// 		printUsage()
// 		return
// 	}

// 	command := os.Args[1]
// 	handler, exists := commands[command]
// 	if !exists {
// 		fmt.Println("Unknown command:", command)
// 		printUsage()
// 		return
// 	}

// 	handler(os.Args[2:])
// }

// func printUsage() {
// 	fmt.Println("Usage:")
// 	fmt.Println("  cli create -p <projectName> -m <moduleName>")
// 	fmt.Println("  cli build")
// 	fmt.Println("  cli deploy [--prod]")
// 	fmt.Println("  cli create-dev-account")
// 	fmt.Println("  cli test-package")
// 	fmt.Println("  cli test-project")
// 	// fmt.Println("  cli import-mainnet-account")
// }
