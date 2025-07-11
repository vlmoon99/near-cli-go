# NEAR CLI Go

A simple CLI tool to manage NEAR smart contract projects. Supports creating, building, deploying contracts, and managing accounts on the NEAR blockchain. Compatible with Linux, macOS, and Windows.

---

## 🚨 **IMPORTANT PREREQUISITES** 🚨

Before using the `near-go` CLI, ensure you have the following tools installed on your PC:

1. **[Go](https://go.dev/doc/install)** – _Required for Go-based development._  
2. **[TinyGo](https://tinygo.org/getting-started/install/)** – _Required for building smart contracts._  

⚠️ **If any of these tools are missing, you won't be able to use the full functionality of the `near-go` CLI.**  
This CLI acts as a proxy to `TinyGo`.  

---

## Installation

```bash
curl -LO https://github.com/vlmoon99/near-cli-go/releases/latest/download/install.sh && bash install.sh && rm install.sh
```

This script installs the `near-go` binary to your `~/bin` directory and adds it to your `PATH`.  
Alternatively, [download from GitHub Releases](https://github.com/vlmoon99/near-cli-go/releases/tag/v1.0.0) and move it manually to your bin.

---

## Uninstall

1. Get path of `near-go` binary:

   ```bash
   which near-go
   ```

2. Remove it:

   ```bash
   rm -rf /home/your-user/bin/near-go
   ```

---

## Usage

Once installed, use the following commands to manage NEAR smart contracts.

### Available Commands

<details>
<summary><strong>1. Create a new project</strong></summary>

```bash
near-go create -p <projectName> -m <moduleName> -t <type of project>

near-go create -p "test1" -m "test1" -t "smart-contract-empty"
```
</details>

<details>
<summary><strong>2. Build the project</strong></summary>

```bash
near-go build
```
Generates a `main.wasm` using TinyGo.
</details>

<details>
<summary><strong>3. Run tests</strong></summary>

```bash
near-go test
```
Runs smart contract tests using TinyGo.
</details>

<details>
<summary><strong>4. Manage blockchain accounts</strong></summary>

```bash
near-go account <subcommand>
```

Examples:

```bash
near-go account create -n "testnet" -a "accountid.testnet"
near-go account import
```
</details>

<details>
<summary><strong>5. Deploy the smart contract</strong></summary>

```bash
near-go deploy -id "accountid.testnet" -n "testnet"
```
</details>

<details>
<summary><strong>6. Call smart contract functions</strong></summary>

```bash
near-go call --from <signer> --to <contract> --function <method> [--args <json>] [--gas <amount>] [--deposit <amount>] --network <network>
```

**Examples:**

Call a view or change method with default args and values:

```bash
near-go call \
  --from neargocli.testnet \
  --to neargocli.testnet \
  --function ReadIncommingTxData \
  --network testnet
```

Call with all parameters set explicitly:

```bash
near-go call \
  --from neargocli.testnet \
  --to neargocli.testnet \
  --function ReadIncommingTxData \
  --args '{}' \
  --gas '100 Tgas' \
  --deposit '0 NEAR' \
  --network testnet
```

Call a method with arguments:

```bash
near-go call \
  --from neargocli.testnet \
  --to neargocli.testnet \
  --function WriteData \
  --args '{"key": "testKey", "data": "test1"}' \
  --gas '100 Tgas' \
  --deposit '0 NEAR' \
  --network testnet
```
</details>

<details>
<summary><strong>7. View help</strong></summary>

```bash
near-go help
```
Displays list of all available commands and help for individual commands.
</details>

---

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
   call     Call a smart contract function
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

---

# Development Setup & Building from Source

If you want to contribute to this project, follow these steps to set up your environment and build the CLI from source.

---

## 1. Setup Repository for Development

Run the following script to download and install all required internal binaries:

```bash
bash setup.sh
```

This will fetch the necessary NEAR CLI binaries and place them in:

```
/near-cli-go/bindata/tools
```

---

## 2. Build near-go from Source

To build the CLI for all supported platforms and architectures, run:

```bash
bash build.sh
```

This script will:

- Build the CLI for:
  - **Linux** (`amd64` and `arm64`)
  - **macOS** (`amd64` and `arm64`)
- Output binaries named according to their platform and architecture (e.g., `near-cli-linux-amd64`).
- Zip each binary for easy distribution (e.g., `near-cli-linux-amd64.zip`).

After running the script, you will find the zipped binaries in the project directory, ready for distribution or testing.

---

**Now you are ready to contribute!**