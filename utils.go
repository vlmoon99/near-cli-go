package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vlmoon99/near-cli-go/bindata"
)

func InitEmbeddedBins() {
	nearCliPath := filepath.Join(os.TempDir(), "near")
	if _, err := os.Stat(nearCliPath); err == nil {
		return
	}
	if err := os.WriteFile(nearCliPath, bindata.NearCli, 0755); err != nil {
		panic("failed to write near-cli: " + err.Error())
	}
}

func CheckDependencies() {
	programs := map[string]string{
		"go":     "Go programming language",
		"tinygo": "TinyGo compiler",
	}
	missing := []string{}
	for prog := range programs {
		if _, err := exec.LookPath(prog); err != nil {
			missing = append(missing, prog)
		}
	}
	if len(missing) > 0 {
		fmt.Printf("Missing dependencies: %s\n", strings.Join(missing, ", "))
		os.Exit(1)
	}
}

func WriteToFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func CreateFolderAndNavigate(name string) error {
	if err := os.MkdirAll(name, os.ModePerm); err != nil {
		return err
	}
	return os.Chdir(name)
}

func ExecuteCommand(name string, args ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errStr := stderr.String()
		if strings.Contains(errStr, "network is unreachable") || strings.Contains(errStr, "no route to host") {
			return nil, fmt.Errorf("%s", ErrNetworkUnreachable)
		}
		return nil, fmt.Errorf("%s: %v: %s", ErrRunningCmd, err, errStr)
	}
	return stdout.Bytes(), nil
}

func ExecuteWithRetry(name string, args []string, dir string, retries int, debug bool) error {
	var lastErr error
	for i := range retries {
		cmd := exec.Command(name, args...)

		// Set the working directory if provided
		if dir != "" {
			cmd.Dir = dir
		}

		output, err := cmd.CombinedOutput()
		if err == nil {
			if debug {
				fmt.Println(string(output))
			}
			return nil
		}
		lastErr = err
		if debug || i == retries-1 {
			fmt.Printf("Attempt %d failed: %s\nOutput: %s\n", i+1, err, string(output))
		}
	}
	return fmt.Errorf("%s: %v", ErrBuildFailed, lastErr)
}
