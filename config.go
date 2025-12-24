package main

import "embed"

//go:embed template/**/*
var templates embed.FS

const (
	NearSdkGoVersion = "v0.0.13"

	SmartContractTypeProject   = "smart-contract-empty"
	SmartContractProjectFolder = "contract"

	ContractMainGoPath     = "template/contract/main.go.template"
	ContractMainGoFileName = "./main.go"

	ErrProvidedNetwork                   = "(USER_INPUT_ERROR): Missing 'network'"
	ErrProvidedNetworkAndAccountName     = "(USER_INPUT_ERROR): Missing both 'network' and 'account-name'"
	ErrProvidedNetworkAndContractId      = "(USER_INPUT_ERROR): Missing both 'network' and 'contract-id'"
	ErrProvidedProjectNameModuleNameType = "(USER_INPUT_ERROR): Missing 'project-name', 'module-name', or 'type'"
	ErrIncorrectType                     = "(USER_INPUT_ERROR): Invalid project type"
	ErrRunningNearCLI                    = "(INTERNAL_UTILS): Failed to execute Near CLI"
	ErrRunningCmd                        = "(INTERNAL_UTILS): Failed to start command"
	ErrGoProjectModFileIsMissing         = "(INTERNAL_PROJECT_CONTRACT): Missing 'go.mod' file"
	ErrGoProjectSumFileIsMissing         = "(INTERNAL_PROJECT_CONTRACT): Missing 'go.sum' file"
	ErrToReadFile                        = "(INTERNAL_PROJECT): Failed to read file"
	ErrBuildFailed                       = "(BUILD_ERROR): Build failed after retries"
	ErrWasmNotFound                      = "(BUILD_ERROR): WASM file not found after build"
	ErrNetworkUnreachable                = "(NETWORK_ERROR): Unable to download dependencies"
)
