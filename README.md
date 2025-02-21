### 🚨 **IMPORTANT PREREQUISITES** 🚨

**Before using the `near-go` CLI, make sure you have the following tools installed on your PC:**

1. **[TinyGo](https://tinygo.org/getting-started/install/)** - **_Required for building smart contracts._**
2. **[near-cli-rs](https://github.com/near/near-cli-rs)** - **_Required for interacting with the NEAR network._**

### ⚠️ **Ensure these tools are installed to avoid errors!** ⚠️

Once these tools are installed, you can proceed with the installation of the `near-go` CLI.

# NEAR CLI Go

This is a simple CLI tool to manage NEAR smart contract projects. It provides functionality to create, build, deploy, and manage developer accounts on the NEAR network. The tool is compatible with Linux, macOS, and Windows.

## Installation

Absolutely! Let's make it clear and prominent to ensure users see it:

```bash
curl -LO https://github.com/vlmoon99/near-cli-go/releases/latest/download/install.sh && bash install.sh
```

This script will download the necessary binary and install it into your local `~/bin` directory. It will also update your `PATH` to ensure the `near-go` binary is accessible globally.

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

1. **Create a new project:**

   ```bash
   near-go create -p my_project -m github.com/myGithubName/myProject
   ```

2. **Build the smart contract:**

   ```bash
   near-go build
   ```

3. **Create a development account (for testnet):**

   ```bash
   near-go create-dev-account
   ```

   This command creates a new dev account on the testnet. You may need to follow any additional instructions that appear after running this command (e.g., setting up credentials).

4. **Deploy the smart contract to testnet:**

   ```bash
   near-go deploy
   ```

5. **Import your mainnet account (before deploying to production):**

   ```bash
   near-go import-mainnet-account
   ```

   This will import your mainnet account for deploying to the production environment. Make sure to have your mainnet credentials available for this step.

6. **Deploy the smart contract to mainnet (optional):**

   ```bash
   near-go deploy --prod
   ```

   This step deploys your smart contract to the mainnet, using the mainnet account you imported earlier.
