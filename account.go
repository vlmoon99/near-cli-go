package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"
// )

// // // ---------------- IMPORT MAINNET ACCOUNT ---------------- //

// // func handleImportMainnetAccount(args []string) {
// // 	scanner := bufio.NewScanner(os.Stdin)

// // 	fmt.Print("Enter your seed phrase (12 words): ")
// // 	scanner.Scan()
// // 	seedPhrase := strings.TrimSpace(scanner.Text())

// // 	fmt.Print("Enter your account ID: ")
// // 	scanner.Scan()
// // 	accountID := strings.TrimSpace(scanner.Text())

// // 	fmt.Println("Importing mainnet account...")
// // 	importCmd := []string{
// // 		"account", "import-account", "using-seed-phrase", seedPhrase,
// // 		"--seed-phrase-hd-path", "m/44'/397'/0'",
// // 		"network-config", "mainnet",
// // 		accountID,
// // 	}
// // 	runCommand("near", importCmd...)

// // 	fmt.Println("Mainnet account imported successfully!")
// // }

// // ---------------- CREATE DEV ACCOUNT ---------------- //

// func handleCreateDevAccount(args []string) {
// 	scanner := bufio.NewScanner(os.Stdin)

// 	fmt.Print("Enter your desired account ID (without the .testnet postfix): ")
// 	scanner.Scan()
// 	accountID := strings.TrimSpace(scanner.Text())

// 	accountID = accountID + ".testnet"

// 	fmt.Println("Creating developer account...")

// 	createCmd := []string{
// 		"account", "create-account", "sponsor-by-faucet-service", accountID,
// 		"autogenerate-new-keypair", "save-to-legacy-keychain", "network-config", "testnet", "create",
// 	}

// 	output, err := runCommand("near", createCmd...)

// 	if err != nil {
// 		log.Fatalf("Error while creating developer account: %v", err)
// 	}

// 	fmt.Println(string(output))

// 	fmt.Println("Developer account created successfully!")
// }
