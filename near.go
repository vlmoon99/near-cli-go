package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func runNearCLI(args ...string) error {
	cmd := exec.Command(filepath.Join(os.TempDir(), "near"), args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", ErrRunningNearCLI, err)
	}
	return nil
}

func HandleDeployContract(id, network string) error {
	args := []string{
		"contract", "deploy", id, "use-file", "./main.wasm",
		"without-init-call", "network-config", network,
		"sign-with-legacy-keychain", "send",
	}
	return runNearCLI(args...)
}

func HandleCreateAccount(network, name string) error {
	if network == "prod" {
		return runNearCLI("account", "create-account", "fund-later", "use-auto-generation", "save-to-folder", "./")
	}
	return runNearCLI("account", "create-account", "sponsor-by-faucet-service", name, 
		"autogenerate-new-keypair", "save-to-legacy-keychain", 
		"network-config", "testnet", "create")
}

func HandleImportAccount() error {
	return runNearCLI("account", "import-account")
}

func HandleCallFunction(signer, contract, method, args, gas, deposit, network string) error {
	cmd := []string{
		"contract", "call-function", "as-transaction", contract, method,
		"json-args", args, "prepaid-gas", gas, "attached-deposit", deposit,
		"sign-as", signer, "network-config", network,
		"sign-with-keychain", "send",
	}
	fmt.Printf("📞 Calling %s on %s...\n", method, contract)
	return runNearCLI(cmd...)
}