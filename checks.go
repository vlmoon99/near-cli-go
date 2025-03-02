package main

import (
	"fmt"
	"os/exec"
)

func checkTinyGo() bool {
	cmd := exec.Command("tinygo", "version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running tinygo:", err)
		return false
	}
	fmt.Println("TinyGo Version:", string(output))
	return true
}

func checkNearRsCli() bool {
	cmd := exec.Command("near", "--version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running near:", err)
		return false
	}
	fmt.Println("Near CLI RS Version:", string(output))
	return true
}
