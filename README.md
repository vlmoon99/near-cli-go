### 🚨 **IMPORTANT PREREQUISITES** 🚨  

Before using the `near-go` CLI, ensure you have the following tools installed on your PC:

1. **[Go](https://go.dev/doc/install)** – _Required for Go-based development._  
2. **[TinyGo](https://tinygo.org/getting-started/install/)** – _Required for building smart contracts._  

⚠️ **If any of these tools are missing, you won't be able to use the full functionality of the `near-go` CLI.**  
This CLI acts as a proxy to `TinyGo`.  

---

# Setup repo for development
```bash
bash setup.sh
```

This script installs all internal bins to the  ```/near-cli-go/bindata/tools```



# NEAR CLI Go

A simple CLI tool to manage NEAR smart contract projects. Supports creating, building, deploying contracts, and managing accounts on the NEAR blockchain. Compatible with Linux, macOS, and Windows.

## Installation

```bash
curl -LO https://github.com/vlmoon99/near-cli-go/releases/latest/download/install.sh && bash install.sh
```

This script installs the `near-go` binary to your `~/bin` directory and adds it to your `PATH`.  
Alternatively, [download from GitHub Releases](https://github.com/vlmoon99/near-cli-go/releases/tag/v1.0.0) and move it manually to your bin.

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

#### 1. **Create a new project**

```bash
near-go create -p <projectName> -m <moduleName> -t <type of project>

near-go create -p "test1" -m "test1" -t "smart-contract-empty"
```

#### 2. **Build the project**

```bash
near-go build
```

Generates a `main.wasm` using TinyGo.

#### 3. **Run tests**

```bash
near-go test
```

Runs smart contract tests using TinyGo.

#### 4. **Manage blockchain accounts**

```bash
near-go account <subcommand>
```

Examples:

```bash
near-go account create -n "testnet" -a "accountid.testnet"
near-go account import
```

#### 5. **Deploy the smart contract**

```bash
near-go deploy -id "accountid.testnet" -n "testnet"
```

#### 6. **Call smart contract functions**

```bash
near-go call --from <signer> --to <contract> --function <method> [--args <json>] [--gas <amount>] [--deposit <amount>] --network <network>
```

**Examples:**

Call a view or change method with default args and values:

```bash
go run main.go call \
  --from neargocli.testnet \
  --to neargocli.testnet \
  --function ReadIncommingTxData \
  --network testnet
```

Call with all parameters set explicitly:

```bash
go run main.go call \
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
go run main.go call \
  --from neargocli.testnet \
  --to neargocli.testnet \
  --function WriteData \
  --args '{"key": "testKey", "data": "test1"}' \
  --gas '100 Tgas' \
  --deposit '0 NEAR' \
  --network testnet
```

#### 7. **View help**

```bash
near-go help
```

Displays list of all available commands and help for individual commands.

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
