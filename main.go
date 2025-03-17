package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

func NearCLIWrapper(args ...string) error {
	cmd := exec.Command("near", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running Near CLI: %w", err)
	}
	return nil
}

// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("Please provide a Near CLI command to execute.")
// 		fmt.Println("Example: ./go-cli-near account")
// 		return
// 	}

// 	args := os.Args[1:]
// if err := NearCLIWrapper(args...); err != nil {
// 	fmt.Printf("Error: %s\n", err)
// }
// }

func main() {
	app := &cli.App{
		Name:  "near-go",
		Usage: "CLI tool for managing projects on Near Blockchain",
		Commands: []cli.Command{
			{
				Name:  "create",
				Usage: "Create a new project",
				Action: func(c *cli.Context) error {
					fmt.Println("Project created successfully!")
					return nil
				},
			},
			{
				Name:  "build",
				Usage: "Build the project",
				Action: func(c *cli.Context) error {
					fmt.Println("Project build completed!")
					return nil
				},
			},
			{
				Name:  "deploy",
				Usage: "Deploy the project to production",
				Action: func(c *cli.Context) error {
					fmt.Println("Project deployed to production!")
					return nil
				},
			},
			{
				Name:  "account",
				Usage: "Manage blockchain accounts",
				Subcommands: []cli.Command{
					{
						Name:  "create",
						Usage: "Create a new account for development",
						Action: func(c *cli.Context) error {
							if err := NearCLIWrapper("account"); err != nil {
								fmt.Printf("Error: %s\n", err)
							}
							fmt.Println("Dev account created!")
							return nil
						},
					},
					{
						Name:  "import",
						Usage: "Import an account",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "type",
								Value: "dev",
								Usage: "Specify account type (e.g., dev, prod)",
							},
						},
						Action: func(c *cli.Context) error {
							fmt.Printf("Account imported for: %s\n", c.String("type"))
							return nil
						},
					},
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
							return nil
						},
					},
					{
						Name:  "package",
						Usage: "Test the package",
						Action: func(c *cli.Context) error {
							fmt.Println("Package test completed!")
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
