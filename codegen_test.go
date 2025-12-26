package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestProject(t *testing.T, content string) string {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "main.go")
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	return dir
}

func TestGenerateCode_StateSerialization(t *testing.T) {
	contractCode := `
package main

// @contract:state
type Contract struct {
	Count int
}
`
	dir := setupTestProject(t, contractCode)

	generated, err := GenerateCode(dir)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if !strings.Contains(generated, "encodingJson \"encoding/json\"") {
		t.Errorf("Expected import encoding/json, but not found")
	}

	if !strings.Contains(generated, "encodingJson.Unmarshal(val, &state)") {
		t.Errorf("Expected JSON Unmarshal for state in getState()")
	}

	if !strings.Contains(generated, "encodingJson.Marshal(state)") {
		t.Errorf("Expected JSON Marshal for state in setState()")
	}

	if strings.Contains(generated, "borsh.Deserialize") {
		t.Errorf("Found forbidden 'borsh.Deserialize' - State should use JSON")
	}
}

func TestGenerateCode_ParameterCapitalization(t *testing.T) {
	contractCode := `
package main

type Contract struct {}

// @contract:mutating
func (c *Contract) SendMessage(newMessage string, userId int, isActive bool) {
	// logic
}
`
	dir := setupTestProject(t, contractCode)

	generated, err := GenerateCode(dir)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	expectedStructField := "NewMessage string `json:\"newMessage\"`"
	if !strings.Contains(generated, expectedStructField) {
		t.Errorf("Expected capitalized struct field '%s', got code:\n%s", expectedStructField, generated)
	}

	expectedStructField2 := "UserId int `json:\"userId\"`"
	if !strings.Contains(generated, expectedStructField2) {
		t.Errorf("Expected capitalized struct field '%s'", expectedStructField2)
	}

	if !strings.Contains(generated, "params.NewMessage") {
		t.Errorf("Expected usage 'params.NewMessage', not found")
	}

	if !strings.Contains(generated, "params.UserId") {
		t.Errorf("Expected usage 'params.UserId', not found")
	}
}

func TestGenerateCode_InitAndPayable(t *testing.T) {
	contractCode := `
package main

type Contract struct {}

// @contract:init
// @contract:payable min_deposit=1NEAR
func (c *Contract) InitContract(startMsg string) {}
`
	dir := setupTestProject(t, contractCode)

	generated, err := GenerateCode(dir)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if !strings.Contains(generated, "existingVal, _ := env.StateRead()") {
		t.Errorf("Init method should check for existing state")
	}
	if !strings.Contains(generated, "state := defaultInit()") {
		t.Errorf("Init method should start with defaultInit()")
	}

	if !strings.Contains(generated, "validatePayment(\"1000000000000000000000000\")") {
		t.Errorf("Expected validation of payment '1NEAR'")
	}
}

func TestGenerateCode_ComplexTypes(t *testing.T) {
	contractCode := `
package main
type Contract struct {}

// @contract:view
func (c *Contract) UpdateList(items []string, data map[string]int) {}
`
	dir := setupTestProject(t, contractCode)
	generated, err := GenerateCode(dir)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if !strings.Contains(generated, "Items []string `json:\"items\"`") {
		t.Errorf("Failed to generate capitalized struct for []string array")
	}
	if !strings.Contains(generated, "Data map[string]int `json:\"data\"`") {
		t.Errorf("Failed to generate capitalized struct for map")
	}
}

func TestGenerateCode_PrivateMethods(t *testing.T) {
	contractCode := `
package main
type Contract struct {}

// @contract:public
func (c *Contract) PublicMethod() {}

// @contract:private
func (c *Contract) InternalHelper() {}

// @contract:view
func (c *Contract) ViewMethod() {}
`
	dir := setupTestProject(t, contractCode)
	generated, err := GenerateCode(dir)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if !strings.Contains(generated, "func public_method()") {
		t.Errorf("Public method failed to export")
	}
	if !strings.Contains(generated, "func view_method()") {
		t.Errorf("View method failed to export")
	}

	if strings.Contains(generated, "func internal_helper()") {
		t.Errorf("SECURITY RISK: Private method was exported!")
	}
}

func TestGenerateCode_ReturnValues(t *testing.T) {
	contractCode := `
package main
type Contract struct {}

// @contract:view
func (c *Contract) GetStatus() (bool, string) {
	return true, "ok"
}

// @contract:view
func (c *Contract) GetMap() map[string]int {
	return nil
}
`
	dir := setupTestProject(t, contractCode)
	generated, err := GenerateCode(dir)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if !strings.Contains(generated, "resultJSON, err := encodingJson.Marshal(result)") {
		t.Errorf("Expected result marshaling for return values")
	}
	if !strings.Contains(generated, "contractBuilder.ReturnValue(string(resultJSON))") {
		t.Errorf("Expected contractBuilder.ReturnValue usage")
	}
}

func TestGenerateCode_PreserveImports(t *testing.T) {
	contractCode := `
package main

import (
	"math/big"
	"fmt"
)

type Contract struct {}

// @contract:mutating
func (c *Contract) BigMath(num *big.Int) {
	fmt.Println("test")
}
`
	dir := setupTestProject(t, contractCode)
	generated, err := GenerateCode(dir)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if !strings.Contains(generated, "\"math/big\"") {
		t.Errorf("Failed to preserve 'math/big' import")
	}
	if !strings.Contains(generated, "\"fmt\"") {
		t.Errorf("Failed to preserve 'fmt' import")
	}
}

func TestGenerateCode_ReceiverNaming(t *testing.T) {
	contractCode := `
package main
type Contract struct {
	Val int
}

// @contract:mutating
func (self *Contract) Increment() {
	self.Val++
}
`
	dir := setupTestProject(t, contractCode)
	generated, err := GenerateCode(dir)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if !strings.Contains(generated, "state.Increment()") {
		t.Errorf("Method call generation failed for receiver named 'self'")
	}
}

func TestGenerateCode_DangerousNames(t *testing.T) {
	contractCode := `
package main
type Contract struct {}

// @contract:mutating
func (c *Contract) EdgeCaseNames(error string, input int, state bool, make []int, new string) {}
`
	dir := setupTestProject(t, contractCode)
	generated, err := GenerateCode(dir)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if !strings.Contains(generated, "Error string `json:\"error\"`") {
		t.Errorf("Failed to generate field for 'error' param")
	}
	if !strings.Contains(generated, "Input int `json:\"input\"`") {
		t.Errorf("Failed to generate field for 'input' param")
	}
	if !strings.Contains(generated, "State bool `json:\"state\"`") {
		t.Errorf("Failed to generate field for 'state' param")
	}
	if !strings.Contains(generated, "Make []int `json:\"make\"`") {
		t.Errorf("Failed to generate field for 'make' param")
	}

	if !strings.Contains(generated, "params.Error") {
		t.Errorf("Method call should use params.Error")
	}
	if !strings.Contains(generated, "params.Input") {
		t.Errorf("Method call should use params.Input")
	}
}

func TestGenerateCode_PointerParameters(t *testing.T) {
	contractCode := `
package main
type Contract struct {}

// @contract:mutating
func (c *Contract) Update(val *int, text *string) {}
`
	dir := setupTestProject(t, contractCode)
	generated, err := GenerateCode(dir)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if !strings.Contains(generated, "Val *int `json:\"val\"`") {
		t.Errorf("Failed to generate pointer field for *int")
	}
	if !strings.Contains(generated, "params.Val") {
		t.Errorf("Failed to pass pointer param")
	}
}

func TestGenerateCode_ValueReceiver(t *testing.T) {
	contractCode := `
package main
type Contract struct {}

// @contract:view
func (c Contract) JustRead() {}
`
	dir := setupTestProject(t, contractCode)
	generated, err := GenerateCode(dir)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if !strings.Contains(generated, "func just_read()") {
		t.Errorf("Failed to generate export for value receiver method")
	}
}

func TestGenerateCode_NoAnnotations(t *testing.T) {
	contractCode := `
package main
type Contract struct {}
func (c *Contract) RegularMethod() {}
`
	dir := setupTestProject(t, contractCode)

	_, err := GenerateCode(dir)
	if err == nil {
		t.Errorf("Expected error when no @contract annotations are present, got nil")
	}

	if !strings.Contains(err.Error(), "no methods or state structs") {
		t.Errorf("Expected specific error message about missing annotations, got: %v", err)
	}
}
