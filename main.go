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
		Name:  "near-go",
		Usage: "CLI tool for managing projects on Near Blockchain",
		Commands: []cli.Command{
			{
				Name:  "create",
				Usage: "Create a new project",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "project-name, p", Required: true},
					&cli.StringFlag{Name: "module-name, m", Required: true},
					&cli.StringFlag{Name: "project-type, t", Required: true},
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
				Usage: "Build the project",
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
						Usage: "Keep the generated intermediate Go file (generated_build.go) after build",
					},
				},
				Action: func(c *cli.Context) error {
					return HandleBuild(c.String("source"), c.String("output"), c.Bool("keep-generated"))
				},
			},
			{
				Name:  "test",
				Usage: "Run tests",
				Subcommands: []cli.Command{
					{
						Name:   "project",
						Action: func(c *cli.Context) error { return HandleTests("project") },
					},
					{
						Name:   "package",
						Action: func(c *cli.Context) error { return HandleTests("package") },
					},
				},
			},
			{
				Name:  "account",
				Usage: "Manage blockchain accounts",
				Subcommands: []cli.Command{
					{
						Name: "create",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "network, n", Required: true},
							&cli.StringFlag{Name: "account-name, a"},
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
						Name: "import",
						Action: func(c *cli.Context) error {
							return HandleImportAccount()
						},
					},
				},
			},
			{
				Name:  "deploy",
				Usage: "Deploy the project",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "contract-id, id", Required: true},
					&cli.StringFlag{Name: "network, n", Required: true},
					&cli.StringFlag{Name: "file, f", Usage: "WASM file to deploy", Value: "main.wasm"},
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
				Usage: "Call a smart contract function",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "signer, from", Required: true},
					&cli.StringFlag{Name: "contract, to", Required: true},
					&cli.StringFlag{Name: "method, function", Required: true},
					&cli.StringFlag{Name: "args", Value: "{}"},
					&cli.StringFlag{Name: "gas", Value: "100 Tgas"},
					&cli.StringFlag{Name: "deposit", Value: "0 NEAR"},
					&cli.StringFlag{Name: "network", Required: true},
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
