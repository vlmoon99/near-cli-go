# NEAR CLI Go

This is a simple CLI tool to manage NEAR smart contract projects. It provides functionality to create, build, deploy, and manage developer accounts on the NEAR network. The tool is compatible with Linux, macOS, and Windows.

## Installation

### Linux / macOS

Run the following command to download and install the CLI tool:

```bash
curl -LO https://github.com/vlmoon99/near-cli-go/blob/main/install.sh && bash install.sh
```

This script will download the necessary binary and install it into your local `~/bin` directory. It will also update your `PATH` to ensure the `near-go` binary is accessible globally.

### Windows

Run the following PowerShell script to install the CLI tool:

```powershell
Invoke-WebRequest -Uri "https://github.com/vlmoon99/near-cli-go/raw/main/install.ps1" -OutFile "install.ps1"; .\install.ps1
```

This will download the `near-go.exe` binary and install it into the `~/bin` directory. It will also update the user-specific `PATH` variable.

## Usage

Once the CLI is installed, you can use the following commands:

### Available Commands

1. **Create a new project:**

   ```bash
   near-go create -p <projectName> -m <moduleName>
   ```

   This creates a new project directory and initializes it with a Go module.

2. **Build the smart contract:**

   ```bash
   near-go build
   ```

   This builds the smart contract using TinyGo and generates the `main.wasm` file.

3. **Deploy the smart contract:**

   ```bash
   near-go deploy [--prod]
   ```

   This deploys the compiled contract to the NEAR network. If `--prod` is specified, it deploys to the mainnet; otherwise, it deploys to the testnet.

4. **Create a developer account:**

   ```bash
   near-go create-dev-account
   ```

   This creates a developer account on the NEAR testnet. You will need to provide an account ID when prompted.

5. **Import a mainnet account:**

   ```bash
   near-go import-mainnet-account
   ```

   This allows you to import an existing NEAR account on the mainnet using a seed phrase.

## Example Workflow

1. Create a new project:
   
   ```bash
   near-go create -p my_project -m my_module
   ```

2. Build the smart contract:
   
   ```bash
   near-go build
   ```

3. Deploy the smart contract to testnet:
   
   ```bash
   near-go deploy
   ```

4. Deploy the smart contract to mainnet (optional):

   ```bash
   near-go deploy --prod
   ```
