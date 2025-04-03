### 🚨 **IMPORTANT PREREQUISITES** 🚨

**Before using the `near-go` CLI, make sure you have the following tools installed on your PC:**  

1. **[Node.js and npm](https://nodejs.org/)** - **_Required for managing JavaScript dependencies._**  
2. **[Yarn](https://yarnpkg.com/getting-started/install)** - **_Alternative package manager for JavaScript dependencies._**  
3. **[Go](https://go.dev/doc/install)** - **_Required for Go-based development._**  
4. **[Rust](https://www.rust-lang.org/tools/install)** - **_Required for Rust-based development and integration testing._**  
5. **[near-cli-rs](https://github.com/near/near-cli-rs)** - **_Required for interacting with the NEAR network._**  
6. **[TinyGo](https://tinygo.org/getting-started/install/)** - **_Required for building smart contracts._**  
7. **[Near CLI](https://github.com/near/near-cli)** - **_Required for interacting with the NEAR network._**  

### ⚠️ **Ensure these tools are installed to avoid errors!** ⚠️


# NEAR CLI Go

This is a simple CLI tool to manage NEAR smart contract projects. It provides functionality to create, build, deploy, and manage developer accounts on the NEAR network. The tool is compatible with Linux, macOS, and Windows.

## Installation
```bash
curl -LO https://github.com/vlmoon99/near-cli-go/releases/latest/download/install.sh && bash install.sh
```
This script will download the necessary binary and install it into your local `~/bin` directory. It will also update your `PATH` to ensure the `near-go` binary is accessible globally.

You can also donwload raw bin from the [Gtihub Releases](https://github.com/vlmoon99/near-cli-go/releases/tag/v1.0.0) and move it into your bin directory.

## Usage

Once the CLI is installed, you can use the following commands to manage projects on the NEAR blockchain.

### Available Commands

1. **Create a new project:**

   ```bash
   near-go create -p <projectName> -m <moduleName>
   ```

   This creates a new project directory and initializes it with a Go module.

2. **Build the project:**

   ```bash
   near-go build
   ```

   This compiles the smart contract using TinyGo and generates the `main.wasm` file.

3. **Run tests:**

   ```bash
   near-go test
   ```

   This runs the tests for the smart contract package and project using the TinyGo test command.

4. **Manage blockchain accounts:**

   ```bash
   near-go account <subcommand>
   ```

   This command provides account management functionalities on the NEAR blockchain.

5. **Deploy the smart contract:**

   ```bash
   near-go deploy [--prod]
   ```

   This deploys the compiled contract to the NEAR network. If `--prod` is specified, it deploys to the mainnet; otherwise, it deploys to the testnet.

6. **View CLI help:**

   ```bash
   near-go help
   ```

   Displays a list of available commands or detailed help for a specific command.

## Example Workflow

1. **Verify installation:**

   ```bash
   go run main.go
   ```

   If all necessary programs are installed, the CLI should display available commands.

2. **Create a new project:**

   ```bash
   near-go create -p my_project -m github.com/myGithubName/myProject
   ```

3. **Build the smart contract:**

   ```bash
   near-go build
   ```

4. **Manage blockchain accounts:**

   ```bash
   near-go account create-dev
   ```

   This command creates a new dev account on the testnet.

5. **Deploy the smart contract to testnet:**

   ```bash
   near-go deploy
   ```

6. **Run smart contract tests:**

   ```bash
   near-go test
   ```

7. **Deploy the smart contract to mainnet (optional):**

   ```bash
   near-go deploy --prod
   ```

   This step deploys your smart contract to the mainnet, using the mainnet account you imported earlier.

### CLI Help Output Example

```bash
(base) test1@test1:~/dev/near-cli-go$ go run main.go
All necessary programs are installed.
NAME:
   near-go - CLI tool for managing projects on Near Blockchain

USAGE:
    [global options] command [command options] [arguments...]

COMMANDS:
   create   Create a new project
   build    Build the project
   test     Run tests
   account  Manage blockchain accounts
   deploy   Deploy the project to production
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```