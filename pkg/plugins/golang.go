// go_extractor.go - Go (Golang) Code Extractor
package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// GoExtractor extracts Go (Golang) code
type GoExtractor struct {
	patterns map[string]*regexp.Regexp
}

// NewGoExtractor creates a new Go extractor
func NewGoExtractor() *GoExtractor {
	return &GoExtractor{}
}

// Name returns the extractor name
func (e *GoExtractor) Name() string {
	return "go"
}

// Extensions returns supported file extensions
func (e *GoExtractor) Extensions() []string {
	return []string{".go", ".mod", ".sum", ".work"}
}

// Initialize sets up regex patterns
func (e *GoExtractor) Initialize() error {
	e.patterns = make(map[string]*regexp.Regexp)

	patterns := map[string]string{
		// Package declaration
		"package": `^package\s+(\w+)`,

		// Import patterns
		"import_single":   `^import\s+"([^"]+)"`,
		"import_multiple": `^import\s*\(([^)]+)\)`,
		"import_group":    `import\s*\(\s*(?:[^)]+\n)+\s*\)`,
		"import_alias":    `import\s+(\w+)\s+"([^"]+)"`,

		// Type declarations
		"struct":     `type\s+(\w+)\s+struct\s*\{`,
		"interface":  `type\s+(\w+)\s+interface\s*\{`,
		"type_alias": `type\s+(\w+)\s+(\w+)`,
		"type_func":  `type\s+(\w+)\s+func\([^)]*\)`,

		// Function patterns
		"function":    `func\s+(\w+)\s*\([^)]*\)(?:\s+\([^)]*\))?\s*\{`,
		"method":      `func\s+\([^)]+\)\s+(\w+)\s*\([^)]*\)(?:\s+\([^)]*\))?\s*\{`,
		"constructor": `func\s+New(\w+)\s*\([^)]*\)`,

		// Receiver patterns
		"pointer_receiver": `func\s+\(\s*\*\s*(\w+)\s*\)`,
		"value_receiver":   `func\s+\(\s*(\w+)\s*\)`,

		// Variable declarations
		"var_single":     `var\s+(\w+)\s+(\w+)`,
		"var_multiple":   `var\s*\(\s*(?:[^)]+\n)+\s*\)`,
		"const_single":   `const\s+(\w+)\s*=`,
		"const_multiple": `const\s*\(\s*(?:[^)]+\n)+\s*\)`,

		// Error patterns
		"error_return": `error\)`,
		"error_check":  `if\s+err\s*!=`,
		"error_wrap":   `fmt\.Errorf|errors\.Wrap`,

		// Go routine patterns
		"goroutine": `go\s+(\w+)\s*\(`,
		"channel":   `chan\s+(?:\*?\s*\w+|struct\s*\{\})`,
		"make_chan": `make\s*\(\s*chan`,
		"select":    `select\s*\{`,

		// Interface implementation patterns
		"interface_impl": `type\s+(\w+)\s+interface\s*\{[^}]*(\w+)\s*\([^)]*\)`,

		// Test patterns
		"test_function":      `func\s+(Test\w+)\s*\(\s*\*\s*testing\.T\s*\)`,
		"benchmark_function": `func\s+(Benchmark\w+)\s*\(\s*\*\s*testing\.B\s*\)`,
		"example_function":   `func\s+(Example\w+)\s*\(\s*\)`,

		// HTTP patterns
		"http_handler": `func\s+(\w+)\s*\(\s*http\.ResponseWriter,\s*\*http\.Request\s*\)`,
		"http_route":   `\.(?:Handle|HandleFunc|Get|Post|Put|Delete|Patch)\s*\(\s*"([^"]+)"`,
		"middleware":   `func\s+(\w+)\s*\(\s*http\.Handler\s*\)\s*http\.Handler`,

		// Database patterns
		"sql_query":  `\.(?:Query|QueryRow|Exec)\s*\(`,
		"gorm_model": `type\s+(\w+)\s+struct\s*\{[^}]*gorm\.Model`,

		// Configuration files
		"go_mod":         `module\s+([^\s]+)`,
		"go_mod_require": `require\s+([^\s]+)\s+([^\s]+)`,
		"go_mod_replace": `replace\s+([^\s]+)\s+=>`,
		"go_work":        `go\s+([\d.]+)`,

		// Comment patterns
		"go_doc":            `^//\s+(\w+)\s`,
		"godoc_package":     `^//\s+Package\s+(\w+)`,
		"godoc_function":    `^//\s+(\w+)\s`,
		"multiline_comment": `/\*[^*]*\*+(?:[^/*][^*]*\*+)*/`,

		// Build constraints
		"build_constraint": `^//\s*\+build\s+`,
		"go_build":         `^//go:build\s+`,

		// Error handling patterns
		"defer":   `defer\s+`,
		"panic":   `panic\s*\(`,
		"recover": `recover\s*\(`,

		// Concurrency patterns
		"sync_mutex":       `sync\.(?:Mutex|RWMutex)`,
		"sync_waitgroup":   `sync\.WaitGroup`,
		"sync_once":        `sync\.Once`,
		"atomic_operation": `atomic\.`,

		// Context patterns
		"context_param": `context\.Context`,
		"with_cancel":   `context\.WithCancel`,
		"with_timeout":  `context\.WithTimeout`,
		"with_deadline": `context\.WithDeadline`,

		// Reflection patterns
		"reflect_type":  `reflect\.TypeOf`,
		"reflect_value": `reflect\.ValueOf`,
		"tag":           `` + "`" + `([^` + "`" + `]+)` + "`",

		// I/O patterns
		"io_reader":     `io\.Reader`,
		"io_writer":     `io\.Writer`,
		"bufio_scanner": `bufio\.Scanner`,

		// JSON patterns
		"json_marshal":   `json\.Marshal`,
		"json_unmarshal": `json\.Unmarshal`,
		"json_tag":       `json:"([^"]+)"`,

		// YAML/XML patterns
		"yaml_tag": `yaml:"([^"]+)"`,
		"xml_tag":  `xml:"([^"]+)"`,

		// Embed patterns
		"embed": `//go:embed\s+`,

		// Generic patterns (Go 1.18+)
		"generic_type": `\[[A-Z]\s+(?:\w+\s*,\s*)*\w+\]`,
		"generic_func": `func\s+\w+\[`,
		"any_type":     `any\b`,
		"comparable":   `comparable\b`,

		// Wire patterns (Dependency Injection)
		"wire_inject":   `wire\.Build`,
		"wire_provider": `func\s+Provide`,

		// Cobra patterns (CLI)
		"cobra_command": `cobra\.Command`,
		"cobra_run":     `Run:\s*func`,

		// Gin patterns (Web Framework)
		"gin_route":   `\.(?:GET|POST|PUT|DELETE)\s*\(\s*"([^"]+)"`,
		"gin_context": `\*gin\.Context`,

		// Echo patterns (Web Framework)
		"echo_route":   `\.(?:GET|POST|PUT|DELETE)\s*\(\s*"([^"]+)"`,
		"echo_context": `echo\.Context`,

		// gRPC patterns
		"grpc_service":     `service\s+(\w+)\s*\{`,
		"grpc_rpc":         `rpc\s+(\w+)\s*\(`,
		"protobuf_message": `message\s+(\w+)\s*\{`,
	}

	for name, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("failed to compile pattern %s: %w", name, err)
		}
		e.patterns[name] = re
	}

	return nil
}

// Cleanup releases resources
func (e *GoExtractor) Cleanup() {
	e.patterns = nil
}

// ShouldProcess checks if this extractor should handle the file
func (e *GoExtractor) ShouldProcess(filename string) bool {
	lowerName := strings.ToLower(filename)

	// Check by extension
	ext := filepath.Ext(lowerName)
	if ext == ".go" || ext == ".mod" || ext == ".sum" || ext == ".work" {
		return true
	}

	// Check by filename
	if filename == "go.mod" || filename == "go.sum" ||
		filename == "go.work" || filename == ".golangci.yml" ||
		filename == "Makefile" || strings.Contains(lowerName, "go.") {
		return true
	}

	// Check for Go build files
	if strings.HasPrefix(filename, "go_") || strings.HasSuffix(filename, "_test.go") {
		return true
	}

	return false
}

// Extract extracts Go code blocks from content
func (e *GoExtractor) Extract(content string) []CodeBlock {
	var blocks []CodeBlock

	// Extract package name
	packageName := e.extractPackageName(content)

	// Extract imports
	imports := e.extractImports(content)

	// Extract Go-specific blocks
	blocks = append(blocks, e.extractStructs(content, packageName, imports)...)
	blocks = append(blocks, e.extractInterfaces(content, packageName, imports)...)
	blocks = append(blocks, e.extractFunctions(content, packageName, imports)...)
	blocks = append(blocks, e.extractMethods(content, packageName, imports)...)
	blocks = append(blocks, e.extractTypes(content, packageName, imports)...)
	blocks = append(blocks, e.extractVariables(content, packageName, imports)...)

	// Extract tests
	blocks = append(blocks, e.extractTests(content, packageName, imports)...)

	// Extract configuration files
	blocks = append(blocks, e.extractGoMod(content)...)
	blocks = append(blocks, e.extractConfigFiles(content)...)

	// Extract web framework code
	blocks = append(blocks, e.extractWebFrameworkCode(content, packageName, imports)...)

	// Extract gRPC/protobuf code
	blocks = append(blocks, e.extractProtobufCode(content, imports)...)

	return blocks
}

// extractPackageName extracts package name from content
func (e *GoExtractor) extractPackageName(content string) string {
	if match := e.patterns["package"].FindStringSubmatch(content); match != nil && len(match) > 1 {
		return match[1]
	}
	return "main"
}

// extractImports extracts import statements
func (e *GoExtractor) extractImports(content string) []string {
	var imports []string

	// Extract single imports
	for _, match := range e.patterns["import_single"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			imports = append(imports, match[1])
		}
	}

	// Extract multiple imports from import block
	if match := e.patterns["import_multiple"].FindStringSubmatch(content); match != nil && len(match) > 1 {
		importBlock := match[1]
		// Split by newlines and extract import paths
		lines := strings.Split(importBlock, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			}
			// Remove quotes and trailing comments
			if strings.Contains(line, `"`) {
				parts := strings.Split(line, `"`)
				if len(parts) >= 2 {
					importPath := parts[1]
					imports = append(imports, importPath)
				}
			}
		}
	}

	// Extract imports with aliases
	for _, match := range e.patterns["import_alias"].FindAllStringSubmatch(content, -1) {
		if len(match) > 2 {
			imports = append(imports, fmt.Sprintf("%s %s", match[1], match[2]))
		}
	}

	return imports
}

// extractStructs extracts Go structs
func (e *GoExtractor) extractStructs(content, packageName string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	for _, match := range e.patterns["struct"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			structName := match[1]
			structContent := e.extractStructBody(content, structName)

			// Check for embedded fields
			embeddedFields := e.extractEmbeddedFields(structContent)

			// Extract tags
			tags := e.extractTags(structContent)

			blocks = append(blocks, CodeBlock{
				Content:     structContent,
				Type:        "struct",
				Package:     packageName,
				Filename:    e.structToFilename(structName),
				Language:    "go",
				Imports:     imports,
				Annotations: append(tags, embeddedFields...),
			})
		}
	}

	// Extract GORM models
	for _, match := range e.patterns["gorm_model"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			modelName := match[1]
			modelContent := e.extractGormModel(content, modelName)

			blocks = append(blocks, CodeBlock{
				Content:     modelContent,
				Type:        "gorm_model",
				Package:     packageName,
				Filename:    e.structToFilename(modelName),
				Language:    "go",
				Imports:     append(imports, "gorm.io/gorm"),
				Annotations: []string{"gorm.Model"},
			})
		}
	}

	return blocks
}

// extractInterfaces extracts Go interfaces
func (e *GoExtractor) extractInterfaces(content, packageName string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	for _, match := range e.patterns["interface"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			interfaceName := match[1]
			interfaceContent := e.extractInterfaceBody(content, interfaceName)

			// Extract method signatures
			methods := e.extractInterfaceMethods(interfaceContent)

			blocks = append(blocks, CodeBlock{
				Content:     interfaceContent,
				Type:        "interface",
				Package:     packageName,
				Filename:    e.interfaceToFilename(interfaceName),
				Language:    "go",
				Imports:     imports,
				Annotations: methods,
			})
		}
	}

	return blocks
}

// extractFunctions extracts Go functions
func (e *GoExtractor) extractFunctions(content, packageName string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	// Extract regular functions
	for _, match := range e.patterns["function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]

			// Skip if it's a method (has receiver)
			funcLine := e.getLineWithPattern(content, `func\s+`+regexp.QuoteMeta(funcName))
			if strings.Contains(funcLine, "(") && strings.Contains(funcLine, ")") &&
				strings.Contains(funcLine, "func") && strings.Index(funcLine, "(") < strings.Index(funcLine, "func") {
				// This is a method, skip it
				continue
			}

			funcContent := e.extractFunctionBody(content, funcName)

			// Check for error returns
			returnsError := e.patterns["error_return"].MatchString(funcContent)

			blocks = append(blocks, CodeBlock{
				Content:  funcContent,
				Type:     "function",
				Package:  packageName,
				Filename: e.functionToFilename(funcName),
				Language: "go",
				Imports:  imports,
				Modifiers: func() []string {
					if returnsError {
						return []string{"error"}
					}
					return nil
				}(),
			})
		}
	}

	// Extract constructors
	for _, match := range e.patterns["constructor"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			typeName := match[1]
			constructorContent := e.extractConstructorBody(content, typeName)

			blocks = append(blocks, CodeBlock{
				Content:  constructorContent,
				Type:     "constructor",
				Package:  packageName,
				Filename: e.functionToFilename("New" + typeName),
				Language: "go",
				Imports:  imports,
			})
		}
	}

	// Extract HTTP handlers
	for _, match := range e.patterns["http_handler"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			handlerName := match[1]
			handlerContent := e.extractHttpHandlerBody(content, handlerName)

			blocks = append(blocks, CodeBlock{
				Content:  handlerContent,
				Type:     "http_handler",
				Package:  packageName,
				Filename: e.functionToFilename(handlerName),
				Language: "go",
				Imports:  append(imports, "net/http"),
			})
		}
	}

	// Extract middleware
	for _, match := range e.patterns["middleware"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			middlewareName := match[1]
			middlewareContent := e.extractMiddlewareBody(content, middlewareName)

			blocks = append(blocks, CodeBlock{
				Content:  middlewareContent,
				Type:     "middleware",
				Package:  packageName,
				Filename: e.functionToFilename(middlewareName),
				Language: "go",
				Imports:  append(imports, "net/http"),
			})
		}
	}

	return blocks
}

// extractMethods extracts Go methods
func (e *GoExtractor) extractMethods(content, packageName string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	// Extract methods with receivers
	for _, match := range e.patterns["method"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			methodName := match[1]
			methodContent := e.extractMethodBody(content, methodName)

			// Check receiver type
			isPointerReceiver := e.patterns["pointer_receiver"].MatchString(methodContent)

			blocks = append(blocks, CodeBlock{
				Content:  methodContent,
				Type:     "method",
				Package:  packageName,
				Filename: e.methodToFilename(methodName),
				Language: "go",
				Imports:  imports,
				Modifiers: func() []string {
					if isPointerReceiver {
						return []string{"pointer_receiver"}
					}
					return []string{"value_receiver"}
				}(),
			})
		}
	}

	return blocks
}

// extractTypes extracts Go type declarations
func (e *GoExtractor) extractTypes(content, packageName string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	// Extract type aliases
	for _, match := range e.patterns["type_alias"].FindAllStringSubmatch(content, -1) {
		if len(match) > 2 {
			typeName := match[1]
			underlyingType := match[2]
			typeContent := fmt.Sprintf("type %s %s", typeName, underlyingType)

			// Skip if it's actually a struct or interface (handled separately)
			if underlyingType == "struct" || underlyingType == "interface" {
				continue
			}

			blocks = append(blocks, CodeBlock{
				Content:  typeContent,
				Type:     "type_alias",
				Package:  packageName,
				Filename: e.typeToFilename(typeName),
				Language: "go",
				Imports:  imports,
			})
		}
	}

	// Extract function types
	for _, match := range e.patterns["type_func"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			typeName := match[1]
			typeContent := e.extractFunctionTypeBody(content, typeName)

			blocks = append(blocks, CodeBlock{
				Content:  typeContent,
				Type:     "func_type",
				Package:  packageName,
				Filename: e.typeToFilename(typeName),
				Language: "go",
				Imports:  imports,
			})
		}
	}

	return blocks
}

// extractVariables extracts Go variables and constants
func (e *GoExtractor) extractVariables(content, packageName string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	// Extract single variables
	for _, match := range e.patterns["var_single"].FindAllStringSubmatch(content, -1) {
		if len(match) > 2 {
			varName := match[1]
			varType := match[2]
			varContent := fmt.Sprintf("var %s %s", varName, varType)

			blocks = append(blocks, CodeBlock{
				Content:  varContent,
				Type:     "variable",
				Package:  packageName,
				Filename: e.variableToFilename(varName),
				Language: "go",
				Imports:  imports,
			})
		}
	}

	// Extract single constants
	for _, match := range e.patterns["const_single"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			constName := match[1]
			constContent := e.extractConstantBody(content, constName)

			blocks = append(blocks, CodeBlock{
				Content:  constContent,
				Type:     "constant",
				Package:  packageName,
				Filename: e.constantToFilename(constName),
				Language: "go",
				Imports:  imports,
			})
		}
	}

	return blocks
}

// extractTests extracts Go tests
func (e *GoExtractor) extractTests(content, packageName string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	// Extract test functions
	for _, match := range e.patterns["test_function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			testName := match[1]
			testContent := e.extractTestBody(content, testName)

			blocks = append(blocks, CodeBlock{
				Content:  testContent,
				Type:     "test",
				Package:  packageName,
				Filename: e.testToFilename(testName),
				Language: "go",
				Imports:  append(imports, "testing"),
			})
		}
	}

	// Extract benchmark functions
	for _, match := range e.patterns["benchmark_function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			benchmarkName := match[1]
			benchmarkContent := e.extractBenchmarkBody(content, benchmarkName)

			blocks = append(blocks, CodeBlock{
				Content:  benchmarkContent,
				Type:     "benchmark",
				Package:  packageName,
				Filename: e.benchmarkToFilename(benchmarkName),
				Language: "go",
				Imports:  append(imports, "testing"),
			})
		}
	}

	// Extract example functions
	for _, match := range e.patterns["example_function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			exampleName := match[1]
			exampleContent := e.extractExampleBody(content, exampleName)

			blocks = append(blocks, CodeBlock{
				Content:  exampleContent,
				Type:     "example",
				Package:  packageName,
				Filename: e.exampleToFilename(exampleName),
				Language: "go",
				Imports:  imports,
			})
		}
	}

	return blocks
}

// extractGoMod extracts go.mod file
func (e *GoExtractor) extractGoMod(content string) []CodeBlock {
	var blocks []CodeBlock

	if e.patterns["go_mod"].MatchString(content) {
		moduleName := ""
		if match := e.patterns["go_mod"].FindStringSubmatch(content); match != nil && len(match) > 1 {
			moduleName = match[1]
		}

		// Extract requirements
		var requires []string
		for _, match := range e.patterns["go_mod_require"].FindAllStringSubmatch(content, -1) {
			if len(match) > 2 {
				requires = append(requires, fmt.Sprintf("%s %s", match[1], match[2]))
			}
		}

		// Extract replacements
		var replacements []string
		for _, match := range e.patterns["go_mod_replace"].FindAllStringSubmatch(content, -1) {
			if len(match) > 1 {
				replacements = append(replacements, match[1])
			}
		}

		goModContent := fmt.Sprintf("module %s\n\ngo 1.21\n", moduleName)

		if len(requires) > 0 {
			goModContent += "\nrequire (\n"
			for _, req := range requires {
				goModContent += fmt.Sprintf("\t%s\n", req)
			}
			goModContent += ")\n"
		}

		if len(replacements) > 0 {
			goModContent += "\nreplace (\n"
			for _, rep := range replacements {
				goModContent += fmt.Sprintf("\t%s => ./local\n", rep)
			}
			goModContent += ")\n"
		}

		blocks = append(blocks, CodeBlock{
			Content:  goModContent,
			Type:     "go_mod",
			Package:  "",
			Filename: "go.mod",
			Language: "mod",
		})
	}

	return blocks
}

// extractConfigFiles extracts other Go configuration files
func (e *GoExtractor) extractConfigFiles(content string) []CodeBlock {
	var blocks []CodeBlock

	// Extract go.work file
	if e.patterns["go_work"].MatchString(content) {
		// Extract Go version
		goVersion := "1.21"
		if match := e.patterns["go_work"].FindStringSubmatch(content); match != nil && len(match) > 1 {
			goVersion = match[1]
		}

		workContent := fmt.Sprintf("go %s\n\nuse (\n\t.\n)", goVersion)

		blocks = append(blocks, CodeBlock{
			Content:  workContent,
			Type:     "go_work",
			Package:  "",
			Filename: "go.work",
			Language: "work",
		})
	}

	return blocks
}

// extractWebFrameworkCode extracts web framework specific code
func (e *GoExtractor) extractWebFrameworkCode(content, packageName string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	// Extract Gin routes
	for _, match := range e.patterns["gin_route"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			routePath := match[1]
			routeContent := e.extractGinRoute(content, routePath)

			blocks = append(blocks, CodeBlock{
				Content:  routeContent,
				Type:     "gin_route",
				Package:  packageName,
				Filename: "routes.go",
				Language: "go",
				Imports:  append(imports, "github.com/gin-gonic/gin"),
			})
		}
	}

	// Extract Echo routes
	for _, match := range e.patterns["echo_route"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			routePath := match[1]
			routeContent := e.extractEchoRoute(content, routePath)

			blocks = append(blocks, CodeBlock{
				Content:  routeContent,
				Type:     "echo_route",
				Package:  packageName,
				Filename: "routes.go",
				Language: "go",
				Imports:  append(imports, "github.com/labstack/echo/v4"),
			})
		}
	}

	// Extract HTTP routes
	for _, match := range e.patterns["http_route"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			routePath := match[1]
			routeContent := e.extractHttpRoute(content, routePath)

			blocks = append(blocks, CodeBlock{
				Content:  routeContent,
				Type:     "http_route",
				Package:  packageName,
				Filename: "routes.go",
				Language: "go",
				Imports:  append(imports, "net/http"),
			})
		}
	}

	// Extract Cobra commands
	if e.patterns["cobra_command"].MatchString(content) {
		cobraContent := e.extractCobraCommand(content)
		if cobraContent != "" {
			blocks = append(blocks, CodeBlock{
				Content:  cobraContent,
				Type:     "cobra_command",
				Package:  packageName,
				Filename: "cmd.go",
				Language: "go",
				Imports:  append(imports, "github.com/spf13/cobra"),
			})
		}
	}

	// Extract Wire providers
	if e.patterns["wire_inject"].MatchString(content) {
		wireContent := e.extractWireInjection(content)
		if wireContent != "" {
			blocks = append(blocks, CodeBlock{
				Content:  wireContent,
				Type:     "wire_injector",
				Package:  packageName,
				Filename: "wire.go",
				Language: "go",
				Imports:  append(imports, "github.com/google/wire"),
			})
		}
	}

	return blocks
}

// extractProtobufCode extracts gRPC/protobuf code
func (e *GoExtractor) extractProtobufCode(content string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	// Extract gRPC services
	for _, match := range e.patterns["grpc_service"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			serviceName := match[1]
			serviceContent := e.extractGrpcService(content, serviceName)

			blocks = append(blocks, CodeBlock{
				Content:  serviceContent,
				Type:     "grpc_service",
				Package:  "",
				Filename: serviceName + ".proto",
				Language: "protobuf",
				Imports:  append(imports, "google.golang.org/grpc"),
			})
		}
	}

	// Extract protobuf messages
	for _, match := range e.patterns["protobuf_message"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			messageName := match[1]
			messageContent := e.extractProtobufMessage(content, messageName)

			blocks = append(blocks, CodeBlock{
				Content:  messageContent,
				Type:     "protobuf_message",
				Package:  "",
				Filename: messageName + ".proto",
				Language: "protobuf",
				Imports:  imports,
			})
		}
	}

	return blocks
}

// Helper Methods

// structToFilename converts struct name to filename
func (e *GoExtractor) structToFilename(structName string) string {
	return strings.ToLower(structName) + ".go"
}

// interfaceToFilename converts interface name to filename
func (e *GoExtractor) interfaceToFilename(interfaceName string) string {
	return strings.ToLower(interfaceName) + ".go"
}

// functionToFilename converts function name to filename
func (e *GoExtractor) functionToFilename(funcName string) string {
	// Convert CamelCase to snake_case for Go conventions
	var result strings.Builder
	for i, r := range funcName {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String()) + ".go"
}

// methodToFilename converts method name to filename
func (e *GoExtractor) methodToFilename(methodName string) string {
	return e.functionToFilename(methodName)
}

// typeToFilename converts type name to filename
func (e *GoExtractor) typeToFilename(typeName string) string {
	return strings.ToLower(typeName) + ".go"
}

// variableToFilename converts variable name to filename
func (e *GoExtractor) variableToFilename(varName string) string {
	return strings.ToLower(varName) + ".go"
}

// constantToFilename converts constant name to filename
func (e *GoExtractor) constantToFilename(constName string) string {
	return strings.ToLower(constName) + ".go"
}

// testToFilename converts test name to filename
func (e *GoExtractor) testToFilename(testName string) string {
	return strings.ToLower(testName) + "_test.go"
}

// benchmarkToFilename converts benchmark name to filename
func (e *GoExtractor) benchmarkToFilename(benchmarkName string) string {
	return strings.ToLower(benchmarkName) + "_bench_test.go"
}

// exampleToFilename converts example name to filename
func (e *GoExtractor) exampleToFilename(exampleName string) string {
	return strings.ToLower(exampleName) + "_example.go"
}

// extractStructBody extracts complete struct body
func (e *GoExtractor) extractStructBody(content, structName string) string {
	// Find struct definition and everything until next type/function at same level
	pattern := regexp.MustCompile(
		`type\s+` + regexp.QuoteMeta(structName) +
			`\s+struct\s*\{[^}]*\}(?:\s*\n\s*\n|\s*$|\s*//)`)

	match := pattern.FindString(content)
	if match != "" {
		return strings.TrimSpace(match)
	}

	// Fallback: simple struct template
	return fmt.Sprintf("type %s struct {\n\t// fields\n}", structName)
}

// extractEmbeddedFields extracts embedded fields from struct
func (e *GoExtractor) extractEmbeddedFields(structContent string) []string {
	var fields []string
	lines := strings.Split(structContent, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//") || line == "" {
			continue
		}
		// Look for embedded fields (no field name before type)
		if !strings.Contains(line, " ") && line != "}" && !strings.Contains(line, "//") {
			fields = append(fields, line)
		}
	}

	return fields
}

// extractTags extracts struct tags from struct content
func (e *GoExtractor) extractTags(structContent string) []string {
	var tags []string

	// Extract JSON tags
	for _, match := range e.patterns["json_tag"].FindAllStringSubmatch(structContent, -1) {
		if len(match) > 1 {
			tags = append(tags, "json:"+match[1])
		}
	}

	// Extract YAML tags
	for _, match := range e.patterns["yaml_tag"].FindAllStringSubmatch(structContent, -1) {
		if len(match) > 1 {
			tags = append(tags, "yaml:"+match[1])
		}
	}

	// Extract XML tags
	for _, match := range e.patterns["xml_tag"].FindAllStringSubmatch(structContent, -1) {
		if len(match) > 1 {
			tags = append(tags, "xml:"+match[1])
		}
	}

	// Extract other tags
	for _, match := range e.patterns["tag"].FindAllStringSubmatch(structContent, -1) {
		if len(match) > 1 {
			tags = append(tags, match[1])
		}
	}

	return tags
}

// extractGormModel extracts GORM model struct
func (e *GoExtractor) extractGormModel(content, modelName string) string {
	pattern := regexp.MustCompile(
		`type\s+` + regexp.QuoteMeta(modelName) +
			`\s+struct\s*\{[^}]*gorm\.Model[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf(`type %s struct {
	gorm.Model
	// fields
}`, modelName)
}

// extractInterfaceBody extracts complete interface body
func (e *GoExtractor) extractInterfaceBody(content, interfaceName string) string {
	pattern := regexp.MustCompile(
		`type\s+` + regexp.QuoteMeta(interfaceName) +
			`\s+interface\s*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf("type %s interface {\n\t// methods\n}", interfaceName)
}

// extractInterfaceMethods extracts method signatures from interface
func (e *GoExtractor) extractInterfaceMethods(interfaceContent string) []string {
	var methods []string
	lines := strings.Split(interfaceContent, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "(") && strings.Contains(line, ")") &&
			!strings.HasPrefix(line, "//") && line != "{" && line != "}" {
			methods = append(methods, line)
		}
	}

	return methods
}

// extractFunctionBody extracts function body
func (e *GoExtractor) extractFunctionBody(content, funcName string) string {
	// Find function with its body
	pattern := regexp.MustCompile(
		`func\s+` + regexp.QuoteMeta(funcName) +
			`\s*\([^)]*\)(?:\s+\([^)]*\))?\s*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf("func %s() {\n\t// function body\n}", funcName)
}

// extractConstructorBody extracts constructor function
func (e *GoExtractor) extractConstructorBody(content, typeName string) string {
	pattern := regexp.MustCompile(
		`func\s+New` + regexp.QuoteMeta(typeName) +
			`\s*\([^)]*\)(?:\s+\*?` + regexp.QuoteMeta(typeName) + `)?\s*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf("func New%s() *%s {\n\treturn &%s{}\n}", typeName, typeName, typeName)
}

// extractHttpHandlerBody extracts HTTP handler function
func (e *GoExtractor) extractHttpHandlerBody(content, handlerName string) string {
	pattern := regexp.MustCompile(
		`func\s+` + regexp.QuoteMeta(handlerName) +
			`\s*\(\s*w\s+http\.ResponseWriter,\s*r\s+\*http\.Request\s*\)\s*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf(`func %s(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}`, handlerName)
}

// extractMiddlewareBody extracts middleware function
func (e *GoExtractor) extractMiddlewareBody(content, middlewareName string) string {
	pattern := regexp.MustCompile(
		`func\s+` + regexp.QuoteMeta(middlewareName) +
			`\s*\(\s*next\s+http\.Handler\s*\)\s*http\.Handler\s*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf(`func %s(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// middleware logic
		next.ServeHTTP(w, r)
	})
}`, middlewareName)
}

// extractMethodBody extracts method body
func (e *GoExtractor) extractMethodBody(content, methodName string) string {
	// Look for method with receiver
	pattern := regexp.MustCompile(
		`func\s+\([^)]+\)\s+` + regexp.QuoteMeta(methodName) +
			`\s*\([^)]*\)(?:\s+\([^)]*\))?\s*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf("func (t *Type) %s() {\n\t// method body\n}", methodName)
}

// extractFunctionTypeBody extracts function type definition
func (e *GoExtractor) extractFunctionTypeBody(content, typeName string) string {
	pattern := regexp.MustCompile(
		`type\s+` + regexp.QuoteMeta(typeName) +
			`\s+func\([^)]*\).*`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf("type %s func()", typeName)
}

// extractConstantBody extracts constant definition
func (e *GoExtractor) extractConstantBody(content, constName string) string {
	pattern := regexp.MustCompile(
		`const\s+` + regexp.QuoteMeta(constName) + `\s*=.*`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf("const %s = \"value\"", constName)
}

// extractTestBody extracts test function
func (e *GoExtractor) extractTestBody(content, testName string) string {
	pattern := regexp.MustCompile(
		`func\s+` + regexp.QuoteMeta(testName) +
			`\s*\(\s*t\s+\*testing\.T\s*\)\s*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf("func %s(t *testing.T) {\n\t// test code\n}", testName)
}

// extractBenchmarkBody extracts benchmark function
func (e *GoExtractor) extractBenchmarkBody(content, benchmarkName string) string {
	pattern := regexp.MustCompile(
		`func\s+` + regexp.QuoteMeta(benchmarkName) +
			`\s*\(\s*b\s+\*testing\.B\s*\)\s*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf("func %s(b *testing.B) {\n\tfor i := 0; i < b.N; i++ {\n\t\t// benchmark code\n\t}\n}", benchmarkName)
}

// extractExampleBody extracts example function
func (e *GoExtractor) extractExampleBody(content, exampleName string) string {
	pattern := regexp.MustCompile(
		`func\s+` + regexp.QuoteMeta(exampleName) +
			`\s*\(\s*\)\s*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf("func %s() {\n\t// example code\n}", exampleName)
}

// extractGinRoute extracts Gin framework route
func (e *GoExtractor) extractGinRoute(content, routePath string) string {
	pattern := regexp.MustCompile(
		`\.(?:GET|POST|PUT|DELETE|PATCH)\s*\(\s*"` +
			regexp.QuoteMeta(routePath) + `"[^)]*\)[^;]*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf(`router.GET("%s", func(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
})`, routePath)
}

// extractEchoRoute extracts Echo framework route
func (e *GoExtractor) extractEchoRoute(content, routePath string) string {
	pattern := regexp.MustCompile(
		`\.(?:GET|POST|PUT|DELETE|PATCH)\s*\(\s*"` +
			regexp.QuoteMeta(routePath) + `"[^)]*\)[^;]*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf(`e.GET("%s", func(c echo.Context) error {
	return c.String(200, "Hello, World!")
})`, routePath)
}

// extractHttpRoute extracts HTTP route
func (e *GoExtractor) extractHttpRoute(content, routePath string) string {
	pattern := regexp.MustCompile(
		`\.(?:Handle|HandleFunc|Get|Post|Put|Delete|Patch)\s*\(\s*"` +
			regexp.QuoteMeta(routePath) + `"[^)]*\)[^;]*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf(`http.HandleFunc("%s", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from %%s", r.URL.Path)
})`, routePath)
}

// extractCobraCommand extracts Cobra CLI command
func (e *GoExtractor) extractCobraCommand(content string) string {
	pattern := regexp.MustCompile(`cobra\.Command\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return `var cmd = &cobra.Command{
	Use:   "app",
	Short: "A brief description",
	Long:  "A longer description",
	Run: func(cmd *cobra.Command, args []string) {
		// command logic
	},
}`
}

// extractWireInjection extracts Wire dependency injection setup
func (e *GoExtractor) extractWireInjection(content string) string {
	if e.patterns["wire_inject"].MatchString(content) {
		pattern := regexp.MustCompile(`wire\.Build\([^)]*\)`)
		match := pattern.FindString(content)
		if match != "" {
			return match
		}
	}

	return `var Set = wire.NewSet(
	// providers
)`
}

// extractGrpcService extracts gRPC service definition
func (e *GoExtractor) extractGrpcService(content, serviceName string) string {
	pattern := regexp.MustCompile(
		`service\s+` + regexp.QuoteMeta(serviceName) + `\s*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf(`service %s {
	rpc Method(Request) returns (Response);
}`, serviceName)
}

// extractProtobufMessage extracts protobuf message definition
func (e *GoExtractor) extractProtobufMessage(content, messageName string) string {
	pattern := regexp.MustCompile(
		`message\s+` + regexp.QuoteMeta(messageName) + `\s*\{[^}]*\}`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf(`message %s {
	string field1 = 1;
	int32 field2 = 2;
}`, messageName)
}

// getLineWithPattern finds the line containing a pattern
func (e *GoExtractor) getLineWithPattern(content, pattern string) string {
	re := regexp.MustCompile(pattern)
	match := re.FindString(content)
	if match != "" {
		return match
	}
	return ""
}
