package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type MethodInfo struct {
	Name         string
	ReceiverType string
	Params       []Param
	Returns      []string
	IsPublic     bool
	IsPrivate    bool
	IsView       bool
	IsMutating   bool
	IsPayable    bool
	IsInit       bool
	MinDeposit   string
	FilePath     string
	RelativePath string
	SourceCode   string
}

type Param struct {
	Name string
	Type string
}

type StateInfo struct {
	Name         string
	Fields       []FieldInfo
	FilePath     string
	RelativePath string
	SourceCode   string
}

type FieldInfo struct {
	Name string
	Type string
}

type FileContent struct {
	FilePath     string
	RelativePath string
	Declarations []string
	Imports      []string
	IsStateFile  bool
}

func GenerateCode(rootDir string) (string, error) {
	fmt.Printf("DEBUG: CodeGen scanning directory: %s\n", rootDir)

	allMethods, stateStructs, fileContents, err := parseAllFilesRecursive(rootDir)
	if err != nil {
		return "", err
	}

	if len(allMethods) == 0 && len(stateStructs) == 0 {
		return "", fmt.Errorf("no methods or state structs with @contract annotations found")
	}

	fmt.Printf("DEBUG: Found %d State Structs and %d Public Methods\n", len(stateStructs), countPublicMethods(allMethods))

	return generateCode(allMethods, stateStructs, fileContents), nil
}

func countPublicMethods(methods []*MethodInfo) int {
	count := 0
	for _, m := range methods {
		if (m.IsPublic || m.IsInit) && !m.IsPrivate {
			count++
		}
	}
	return count
}

func parseAllFilesRecursive(rootDir string) ([]*MethodInfo, []*StateInfo, []*FileContent, error) {
	var allMethods []*MethodInfo
	var stateStructs []*StateInfo
	var fileContents []*FileContent

	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") && name != "." && name != "./" {
				return filepath.SkipDir
			}
			if name == "vendor" || name == "node_modules" || name == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if strings.HasPrefix(filepath.Base(path), "generated_") {
			return nil
		}

		relPath, _ := filepath.Rel(rootDir, path)

		methods, states, content, err := parseContract(path, relPath)
		if err != nil {
			fmt.Printf("DEBUG: Failed to parse %s: %v\n", relPath, err)
			return nil
		}

		if content == nil {
			return nil
		}

		if len(methods) > 0 {
			allMethods = append(allMethods, methods...)
		}

		if len(states) > 0 {
			stateStructs = append(stateStructs, states...)
			content.IsStateFile = true
		}

		fileContents = append(fileContents, content)
		return nil
	})

	return allMethods, stateStructs, fileContents, err
}

func parseContract(filePath string, relativePath string) ([]*MethodInfo, []*StateInfo, *FileContent, error) {
	fileContentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, nil, err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, nil, err
	}

	if file.Name.Name != "main" {
		return nil, nil, nil, nil
	}

	var methods []*MethodInfo
	var stateStructs []*StateInfo
	content := &FileContent{
		FilePath:     filePath,
		RelativePath: relativePath,
		Declarations: []string{},
		Imports:      []string{},
	}

	for _, imp := range file.Imports {
		startPos := fset.Position(imp.Pos()).Offset
		endPos := fset.Position(imp.End()).Offset
		content.Imports = append(content.Imports, string(fileContentBytes[startPos:endPos]))
	}

	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.IMPORT {
				continue
			}

			// Start by assuming the declaration position
			startPos := fset.Position(d.Pos()).Offset

			// FIX: If comments exist (like //go:embed), start from the comment position
			if d.Doc != nil {
				startPos = fset.Position(d.Doc.Pos()).Offset
			}

			endPos := fset.Position(d.End()).Offset
			declCode := string(fileContentBytes[startPos:endPos])

			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					isState := false
					if d.Doc != nil && hasStateAnnotation(d.Doc) {
						isState = true
					} else if typeSpec.Doc != nil && hasStateAnnotation(typeSpec.Doc) {
						isState = true
					}

					if isState {
						structType, ok := typeSpec.Type.(*ast.StructType)
						if !ok {
							continue
						}
						state := extractStateInfo(typeSpec, structType, fset, fileContentBytes)
						state.FilePath = filePath
						state.RelativePath = relativePath
						stateStructs = append(stateStructs, state)
					}
				}
			}
			content.Declarations = append(content.Declarations, declCode)

		case *ast.FuncDecl:
			// Start by assuming the function declaration position
			startPos := fset.Position(d.Pos()).Offset

			// FIX: If comments exist (like //go:export), start from the comment position
			if d.Doc != nil {
				startPos = fset.Position(d.Doc.Pos()).Offset
			}

			endPos := fset.Position(d.End()).Offset
			declCode := string(fileContentBytes[startPos:endPos])

			if d.Recv != nil && len(d.Recv.List) > 0 {
				method := extractMethodWithSource(d, fset, fileContentBytes)
				method.FilePath = filePath
				method.RelativePath = relativePath

				if method.IsPublic || method.IsView || method.IsMutating || method.IsPayable || method.IsPrivate || method.IsInit {
					methods = append(methods, method)
				}
			}
			content.Declarations = append(content.Declarations, declCode)
		}
	}

	return methods, stateStructs, content, nil
}

func hasStateAnnotation(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}
	for _, comment := range doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(comment.Text), "//"))
		if strings.Contains(text, "@contract:state") {
			return true
		}
	}
	return false
}

func extractStateInfo(typeSpec *ast.TypeSpec, structType *ast.StructType, fset *token.FileSet, fileContent []byte) *StateInfo {
	state := &StateInfo{Name: typeSpec.Name.Name}
	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			fieldType := typeToString(field.Type)
			for _, name := range field.Names {
				state.Fields = append(state.Fields, FieldInfo{Name: name.Name, Type: fieldType})
			}
		}
	}
	startPos := fset.Position(typeSpec.Pos()).Offset
	endPos := fset.Position(structType.End()).Offset
	state.SourceCode = string(fileContent[startPos:endPos])
	return state
}

func extractMethodWithSource(fn *ast.FuncDecl, fset *token.FileSet, fileContent []byte) *MethodInfo {
	method := &MethodInfo{Name: fn.Name.Name}
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		method.ReceiverType = extractReceiverType(fn.Recv.List[0].Type)
	}
	if fn.Doc != nil {
		for _, comment := range fn.Doc.List {
			parseAnnotation(comment.Text, method)
		}
	}
	if fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			typeName := typeToString(field.Type)
			for _, name := range field.Names {
				method.Params = append(method.Params, Param{Name: name.Name, Type: typeName})
			}
		}
	}
	if fn.Type.Results != nil {
		for _, field := range fn.Type.Results.List {
			method.Returns = append(method.Returns, typeToString(field.Type))
		}
	}
	startPos := fset.Position(fn.Pos()).Offset
	endPos := fset.Position(fn.End()).Offset
	method.SourceCode = string(fileContent[startPos:endPos])
	return method
}

func extractReceiverType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return extractReceiverType(t.X)
	default:
		return "Unknown"
	}
}

func parseAnnotation(text string, method *MethodInfo) {
	text = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(text), "//"))
	if !strings.HasPrefix(text, "@contract:") {
		return
	}
	parts := strings.Fields(strings.TrimPrefix(text, "@contract:"))
	if len(parts) == 0 {
		return
	}
	switch parts[0] {
	case "init":
		method.IsInit = true
		method.IsMutating = true
		method.IsPublic = true
	case "public":
		method.IsPublic = true
	case "private":
		method.IsPrivate = true
	case "view":
		method.IsView = true
		method.IsPublic = true
	case "mutating":
		method.IsMutating = true
		method.IsPublic = true
	case "payable":
		method.IsPayable = true
		method.IsPublic = true
		for _, part := range parts[1:] {
			if strings.HasPrefix(part, "min_deposit=") {
				method.MinDeposit = strings.TrimPrefix(part, "min_deposit=")
			}
		}
	}
}

func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	case *ast.MapType:
		return "map[" + typeToString(t.Key) + "]" + typeToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return "unknown"
	}
}

func generateCode(methods []*MethodInfo, stateStructs []*StateInfo, fileContents []*FileContent) string {
	var sb strings.Builder

	sb.WriteString("// Code generated by NEAR contract generator. DO NOT EDIT.\n")
	sb.WriteString("// This file uses encoding/json for both state serialization and parameter parsing.\n\n")
	sb.WriteString("package main\n\n")

	importMap := make(map[string]bool)
	importMap["contractBuilder \"github.com/vlmoon99/near-sdk-go/contract\""] = true
	importMap["\"github.com/vlmoon99/near-sdk-go/env\""] = true
	importMap["\"github.com/vlmoon99/near-sdk-go/types\""] = true
	importMap["encodingJson \"encoding/json\""] = true
	importMap["\"strconv\""] = true

	for _, content := range fileContents {
		for _, imp := range content.Imports {
			cleanImp := strings.TrimSpace(imp)
			if cleanImp != "" {
				importMap[cleanImp] = true
			}
		}
	}

	sb.WriteString("import (\n")
	for imp := range importMap {
		sb.WriteString("\t" + imp + "\n")
	}
	sb.WriteString(")\n\n")

	for _, content := range fileContents {
		sb.WriteString(fmt.Sprintf("// ===== From: %s =====\n", content.RelativePath))
		for _, decl := range content.Declarations {
			sb.WriteString(decl)
			sb.WriteString("\n\n")
		}
	}

	if len(stateStructs) > 0 {
		state := stateStructs[0]
		sb.WriteString(generateDefaultInit(state))
		sb.WriteString("\n")
		sb.WriteString(generateGetState(state))
		sb.WriteString("\n")
		sb.WriteString(generateSetState(state))
		sb.WriteString("\n")
	}

	sb.WriteString("// ===== Generated Exports =====\n")
	for _, m := range methods {
		if m.IsPrivate {
			continue
		}
		if !m.IsPublic && !m.IsInit {
			continue
		}
		sb.WriteString(generateExportFunction(m))
		sb.WriteString("\n")
	}

	sb.WriteString("// ===== Helper Functions =====\n")
	sb.WriteString(generateValidatePayment())
	sb.WriteString("\n")

	return sb.String()
}

func generateDefaultInit(state *StateInfo) string {
	return fmt.Sprintf(`func defaultInit() *%s {
	return &%s{}
}`, state.Name, state.Name)
}

func generateGetState(state *StateInfo) string {
	return fmt.Sprintf(`func getState() *%s {
	val, err := env.StateRead()
	if err != nil || len(val) == 0 {
		return defaultInit()
	}
	var state %s
	err = encodingJson.Unmarshal(val, &state)
	if err != nil {
		env.PanicStr("Failed to deserialize state")
	}
	return &state
}`, state.Name, state.Name)
}

func generateSetState(state *StateInfo) string {
	return fmt.Sprintf(`func setState(state *%s) {
	val, err := encodingJson.Marshal(state)
	if err != nil {
		env.PanicStr("Failed to serialize state")
	}
	err = env.StateWrite(val)
	if err != nil {
		env.PanicStr("Failed to write state")
	}
}`, state.Name)
}

func generateValidatePayment() string {
	return `func validatePayment(minDeposit string) bool {
	minAmount, err := strconv.ParseFloat(minDeposit, 64)
	if err != nil {
		env.LogString("Invalid min deposit amount: " + minDeposit)
		return false
	}
	minYocto := minAmount * 1e24
	minYoctoStr := strconv.FormatFloat(minYocto, 'f', 0, 64)
	minRequired, err := types.U128FromString(minYoctoStr)
	if err != nil {
		env.LogString("Failed to create Uint128: " + err.Error())
		return false
	}
	attachedDeposit, err := env.GetAttachedDeposit()
	if attachedDeposit.Cmp(minRequired) < 0 {
		env.LogString("Insufficient payment")
		return false
	}
	return true
}`
}

func generateExportFunction(m *MethodInfo) string {
	var sb strings.Builder

	exportName := toSnakeCase(m.Name)
	sb.WriteString(fmt.Sprintf("// Export: %s (from %s)\n", exportName, m.RelativePath))
	sb.WriteString(fmt.Sprintf("//go:export %s\n", exportName))
	sb.WriteString(fmt.Sprintf("func %s() {\n", exportName))
	sb.WriteString("\tcontractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {\n")

	if m.IsInit {
		sb.WriteString("\t\t// Initialization: Check if already initialized\n")
		sb.WriteString("\t\texistingVal, _ := env.StateRead()\n")
		sb.WriteString("\t\tif len(existingVal) > 0 {\n")
		sb.WriteString("\t\t\tenv.PanicStr(\"Contract already initialized\")\n")
		sb.WriteString("\t\t}\n")
		sb.WriteString("\t\tstate := defaultInit()\n\n")
	} else {
		sb.WriteString("\t\tstate := getState()\n\n")
	}

	if m.IsPayable {
		sb.WriteString(fmt.Sprintf("\t\tif !validatePayment(\"%s\") {\n", m.MinDeposit))
		sb.WriteString("\t\t\tenv.PanicStr(\"Insufficient payment\")\n")
		sb.WriteString("\t\t}\n\n")
	}

	sb.WriteString(generateParamParser(m))
	sb.WriteString("\n")

	returnsError := false
	if len(m.Returns) > 0 && m.Returns[len(m.Returns)-1] == "error" {
		returnsError = true
	}

	hasDataResult := false
	if returnsError {
		if len(m.Returns) > 1 {
			hasDataResult = true
		}
	} else {
		if len(m.Returns) > 0 {
			hasDataResult = true
		}
	}

	sb.WriteString("\t\t// Call method\n")
	sb.WriteString("\t\t")

	if hasDataResult && returnsError {
		sb.WriteString("result, callErr := ")
	} else if hasDataResult {
		sb.WriteString("result := ")
	} else if returnsError {
		sb.WriteString("callErr := ")
	}

	sb.WriteString("state.")
	sb.WriteString(m.Name)
	sb.WriteString("(")

	if len(m.Params) > 0 {
		if len(m.Params) == 1 && !isBasicType(m.Params[0].Type) {
			sb.WriteString("params")
		} else {
			for i, p := range m.Params {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString("params.")
				sb.WriteString(capitalizeFirst(p.Name))
			}
		}
	}
	sb.WriteString(")\n\n")

	if returnsError {
		sb.WriteString("\t\tif callErr != nil {\n")
		sb.WriteString("\t\t\tenv.PanicStr(callErr.Error())\n")
		sb.WriteString("\t\t}\n\n")
	}

	if m.IsMutating {
		sb.WriteString("\t\tsetState(state)\n\n")
	}

	if hasDataResult {
		sb.WriteString("\t\tresultJSON, err := encodingJson.Marshal(result)\n")
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString("\t\t\tenv.PanicStr(\"Failed to marshal result to JSON\")\n")
		sb.WriteString("\t\t}\n")
		sb.WriteString("\t\tcontractBuilder.ReturnValue(string(resultJSON))\n")
	}

	sb.WriteString("\t\treturn nil\n")
	sb.WriteString("\t})\n")
	sb.WriteString("}\n")

	return sb.String()
}

func generateParamParser(method *MethodInfo) string {
	if len(method.Params) == 0 {
		return "\t\t// No parameters to parse\n"
	}

	var sb strings.Builder
	sb.WriteString("\t\t// Parse input parameters from JSON\n")

	if len(method.Params) == 1 && !isBasicType(method.Params[0].Type) {
		p := method.Params[0]
		sb.WriteString(fmt.Sprintf("\t\tvar params %s\n", p.Type))
		sb.WriteString("\t\terr := encodingJson.Unmarshal(input.Data, &params)\n")
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString(fmt.Sprintf("\t\t\tenv.LogString(\"JSON unmarshal error for %s: \" + err.Error())\n", p.Type))
		sb.WriteString(fmt.Sprintf("\t\t\tenv.PanicStr(\"Failed to parse %s parameter\")\n", p.Type))
		sb.WriteString("\t\t}\n")
	} else {
		sb.WriteString("\t\tvar params struct {\n")
		for _, p := range method.Params {
			fieldName := capitalizeFirst(p.Name)
			jsonTag := fmt.Sprintf("`json:\"%s\"`", p.Name)
			sb.WriteString(fmt.Sprintf("\t\t\t%s %s %s\n", fieldName, p.Type, jsonTag))
		}
		sb.WriteString("\t\t}\n")

		sb.WriteString("\t\terr := encodingJson.Unmarshal(input.Data, &params)\n")
		sb.WriteString("\t\tif err != nil {\n")
		sb.WriteString("\t\t\tenv.LogString(\"JSON unmarshal error: \" + err.Error())\n")
		sb.WriteString("\t\t\tenv.PanicStr(\"Failed to parse input parameters\")\n")
		sb.WriteString("\t\t}\n")
	}

	return sb.String()
}

func isBasicType(typeStr string) bool {
	basicTypes := map[string]bool{
		"string":                 true,
		"bool":                   true,
		"int":                    true,
		"int8":                   true,
		"int16":                  true,
		"int32":                  true,
		"int64":                  true,
		"uint":                   true,
		"uint8":                  true,
		"uint16":                 true,
		"uint32":                 true,
		"uint64":                 true,
		"uintptr":                true,
		"byte":                   true,
		"rune":                   true,
		"float32":                true,
		"float64":                true,
		"complex64":              true,
		"complex128":             true,
		"[]byte":                 true,
		"[]string":               true,
		"[]int":                  true,
		"[]int64":                true,
		"[]float64":              true,
		"[]bool":                 true,
		"map[string]string":      true,
		"map[string]interface{}": true,
	}

	if basicTypes[typeStr] {
		return true
	}
	if strings.HasPrefix(typeStr, "*") {
		baseType := strings.TrimPrefix(typeStr, "*")
		return basicTypes[baseType]
	}
	if strings.HasPrefix(typeStr, "[]") {
		baseType := strings.TrimPrefix(typeStr, "[]")
		return basicTypes[baseType] || isBasicType(baseType)
	}
	if strings.HasPrefix(typeStr, "map[") {
		return true
	}
	return false
}

func toSnakeCase(s string) string {
	var result strings.Builder
	runes := []rune(s)
	length := len(runes)

	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			if (unicode.IsLower(prev) || unicode.IsDigit(prev)) ||
				(unicode.IsUpper(prev) && i+1 < length && unicode.IsLower(runes[i+1])) {
				result.WriteRune('_')
			}
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func capitalizeFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
