package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vlmoon99/near-cli-go/bindata"
)

// Constants for paths
const (
	ToolDirName = ".near-go"
)

func getToolHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to temp if home is not available
		return os.TempDir()
	}
	return filepath.Join(home, ToolDirName)
}

func InitEmbeddedBins() {
	toolHome := getToolHome()

	// Create ~/.near-go if it doesn't exist
	if err := os.MkdirAll(toolHome, 0755); err != nil {
		panic("failed to create tool home directory: " + err.Error())
	}

	nearCliPath := filepath.Join(toolHome, "near")
	if _, err := os.Stat(nearCliPath); err != nil {
		if err := os.WriteFile(nearCliPath, bindata.NearCli, 0755); err != nil {
			panic("failed to write near-cli: " + err.Error())
		}
	}

	tinyGoDir := filepath.Join(toolHome, "tinygo")
	tinyGoBinDir := filepath.Join(tinyGoDir, "bin")
	tinyGoMainBin := filepath.Join(tinyGoBinDir, "tinygo")

	if _, err := os.Stat(tinyGoMainBin); err == nil {
		return
	}

	fmt.Println("Extracting embedded TinyGo... (this happens once)")
	if err := Unzip(bindata.TinyGoZip, toolHome); err != nil {
		panic("failed to extract tinygo: " + err.Error())
	}

	entries, err := os.ReadDir(tinyGoBinDir)
	if err != nil {
		fmt.Printf("Warning: could not read tinygo bin directory to set permissions: %v\n", err)
		return
	}

	for _, entry := range entries {
		binPath := filepath.Join(tinyGoBinDir, entry.Name())

		if err := os.Chmod(binPath, 0755); err != nil {
			fmt.Printf("Warning: failed to chmod %s: %v\n", entry.Name(), err)
		}
	}
}

func CheckDependencies() {
	programs := map[string]string{
		"go": "Go programming language",
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

func GetTinyGoPath() string {
	return filepath.Join(getToolHome(), "tinygo", "bin", "tinygo")
}

func Unzip(src []byte, dest string) error {
	r, err := zip.NewReader(bytes.NewReader(src), int64(len(src)))
	if err != nil {
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
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
