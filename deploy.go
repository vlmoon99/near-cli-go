package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"
// )

// // ---------------- DEPLOY COMMAND ---------------- //

// func handleDeploy(args []string) {
// 	var smartContractID string
// 	isProd := len(args) > 0 && args[0] == "--prod"
// 	scanner := bufio.NewScanner(os.Stdin)

// 	if !isProd {
// 		fmt.Print("Enter your smart contract account ID (without testnet prefix): ")
// 	}
// 	if isProd {
// 		fmt.Print("Enter your smart contract account ID (with all necessary prefixes): ")
// 	}
// 	scanner.Scan()
// 	smartContractID = strings.TrimSpace(scanner.Text())

// 	if !isProd {
// 		smartContractID = smartContractID + ".testnet"
// 	}

// 	fmt.Println("Building for deployment...")

// 	if err := buildSmartContract(); err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Verifying build...")
// 	runCommand("ls", "-lh", "main.wasm")

// 	network := "testnet"
// 	if isProd {
// 		network = "mainnet"
// 	}

// 	fmt.Println("Deploying contract...")
// deployCmd := []string{
// 	"contract", "deploy", smartContractID, "use-file", "./main.wasm",
// 	"without-init-call",
// 	"network-config", network,
// 	"sign-with-legacy-keychain",
// 	"send",
// }
// 	output, err := runCommand("near", deployCmd...)
// 	if err != nil {
// 		log.Fatalf("Error deploying contract: %v\n%s", err, string(output))
// 	}

// 	fmt.Println("Deployment complete!")
// 	fmt.Printf("Deploy command output: %v", string(output))
// }
