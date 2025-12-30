package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	InitEmbeddedBins()
	CheckDependencies()

	app := &cli.App{
		Name:    "near-go",
		Usage:   "CLI tool for managing projects on Near Blockchain",
		Version: NearSdkGoVersion,
		Authors: []cli.Author{
			{Name: "Github : vlmoon99, Telegram : @vlmoon99"},
		},
		Description: "A comprehensive toolchain for scaffolding, building, testing, and deploying NEAR smart contracts written in Go. It utilizes TinyGo for WASM compilation and an annotation-based code generator for boilerplate reduction.",
		Commands: []cli.Command{
			{
				Name:  "create",
				Usage: "Scaffold a new smart contract project",
				Description: "Creates a standard directory structure, initializes a go.mod file, " +
					"and downloads the required NEAR Go SDK dependencies.",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-name, p", Required: true, Usage: "Name of the project folder to create"},
					&cli.StringFlag{Name: "module-name, m", Required: true, Usage: "Go module name (e.g., github.com/user/project)"},
					&cli.StringFlag{Name: "project-type, t", Required: true, Usage: "Type of project (e.g., smart-contract-empty)"},
				},
				Action: func(c *cli.Context) error {
					if c.String("project-name") == "" || c.String("module-name") == "" {
						return errors.New(ErrProvidedProjectNameModuleNameType)
					}
					return HandleCreateProject(c.String("project-name"), c.String("project-type"), c.String("module-name"))
				},
			},
			{
				Name:  "build",
				Usage: "Compile the smart contract to WASM",
				Description: `Executes the full build pipeline using Comment Directives:
   
   1. Scans for @contract annotations:
      - @contract:state: Identifies the main state struct (Only 1 allowed).
      - @contract:init: Marks the initialization method (Only 1 allowed).
      - @contract:view: Read-only method. Compatible with promise_callback.
      - @contract:mutating: Modifies state. Compatible with payable and promise_callback.
      - @contract:payable: Accepts attached NEAR.
      - @contract:promise_callback: Handles async promise results. Must be combined with 'view' or 'mutating'.

   2. Generates 'generated_build.go' with JSON logic and SDK glue code.
   3. Compiles using TinyGo to WASM.`,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "source, s",
						Usage: "Source directory to scan for smart contract code",
						Value: "./",
					},
					&cli.StringFlag{
						Name:  "output, o",
						Usage: "Output filename for the WASM binary",
						Value: "main.wasm",
					},
					&cli.BoolFlag{
						Name:  "keep-generated, k",
						Usage: "Keep the intermediate 'generated_build.go' file for inspection/debugging",
					},
				},
				Action: func(c *cli.Context) error {
					return HandleBuild(c.String("source"), c.String("output"), c.Bool("keep-generated"))
				},
			},
			{
				Name:  "test",
				Usage: "Run contract unit tests using TinyGo",
				Description: "Executes Go tests using the TinyGo compiler to simulate the WASM environment constraints. " +
					"This ensures dependencies and logic are compatible with the strict requirements of the NEAR runtime.",
				Subcommands: []cli.Command{
					{
						Name:   "project",
						Usage:  "Run tests recursively for the entire project (./...)",
						Action: func(c *cli.Context) error { return HandleTests("project") },
					},
					{
						Name:   "package",
						Usage:  "Run tests only for the current directory (./)",
						Action: func(c *cli.Context) error { return HandleTests("package") },
					},
				},
			},
			{
				Name:  "account",
				Usage: "Manage NEAR blockchain accounts",
				Subcommands: []cli.Command{
					{
						Name:  "create",
						Usage: "Create a new account (dev or testnet)",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "network, n", Required: true, Usage: "Network ID (testnet, mainnet, dev)"},
							&cli.StringFlag{Name: "account-name, a", Usage: "Desired account name (required for non-dev networks)"},
						},
						Action: func(c *cli.Context) error {
							net, name := c.String("network"), c.String("account-name")
							if net == "" {
								return errors.New(ErrProvidedNetwork)
							}
							if net == "dev" && name == "" {
								return errors.New(ErrProvidedNetworkAndAccountName)
							}
							return HandleCreateAccount(net, name)
						},
					},
					{
						Name:  "import",
						Usage: "Import an existing account via private key",
						Action: func(c *cli.Context) error {
							return HandleImportAccount()
						},
					},
				},
			},
			{
				Name:  "deploy",
				Usage: "Deploy a compiled WASM contract",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contract-id, id", Required: true, Usage: "Account ID to deploy the contract to"},
					&cli.StringFlag{Name: "network, n", Required: true, Usage: "Network ID (testnet, mainnet)"},
					&cli.StringFlag{Name: "file, f", Usage: "Path to WASM file", Value: "main.wasm"},
				},
				Action: func(c *cli.Context) error {
					id, net := c.String("contract-id"), c.String("network")
					if id == "" || net == "" {
						return errors.New(ErrProvidedNetworkAndContractId)
					}
					return HandleDeployContract(id, net)
				},
			},
			{
				Name:  "call",
				Usage: "Invoke a method on a smart contract",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "signer, from", Required: true, Usage: "Account ID signing the transaction"},
					&cli.StringFlag{Name: "contract, to", Required: true, Usage: "Contract Account ID"},
					&cli.StringFlag{Name: "method, function", Required: true, Usage: "Method name to invoke"},
					&cli.StringFlag{Name: "args", Value: "{}", Usage: "JSON arguments string"},
					&cli.StringFlag{Name: "gas", Value: "100 Tgas", Usage: "Prepaid gas"},
					&cli.StringFlag{Name: "deposit", Value: "0 NEAR", Usage: "Attached deposit"},
					&cli.StringFlag{Name: "network", Required: true, Usage: "Network ID"},
				},
				Action: func(c *cli.Context) error {
					return HandleCallFunction(
						c.String("signer"), c.String("contract"),
						c.String("method"), c.String("args"),
						c.String("gas"), c.String("deposit"), c.String("network"),
					)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
