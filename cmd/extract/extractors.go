package extract

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bhangun/coto/pkg/extractor"
)

// JavaExtractor extracts Java code
type JavaExtractor struct {
	patterns map[string]*regexp.Regexp
}

func NewJavaExtractor() *JavaExtractor {
	return &JavaExtractor{}
}

func (e *JavaExtractor) Name() string { return "java" }

func (e *JavaExtractor) Extensions() []string {
	return []string{".java", ".jar"}
}

func (e *JavaExtractor) Initialize() error {
	e.patterns = make(map[string]*regexp.Regexp)

	// Compile regex patterns
	patterns := map[string]string{
		"class":      `(?s)(?:(public|private|protected|abstract|final)\s+)?class\s+(\w+)(?:<[^>]+>)?(?:\s+extends\s+\w+)?(?:\s+implements\s+[^{]+)?\s*\{(.*?)\}`,
		"interface":  `(?s)(?:(public|private|protected)\s+)?interface\s+(\w+)(?:<[^>]+>)?(?:\s+extends\s+[^{]+)?\s*\{(.*?)\}`,
		"enum":       `(?s)(?:(public|private|protected)\s+)?enum\s+(\w+)\s*\{(.*?)\}`,
		"package":    `package\s+([\w.]+)\s*;`,
		"import":     `import\s+(?:static\s+)?([\w.*]+)\s*;`,
		"annotation": `@(\w+)`,
		"maven_pom":  `(?s)<project[^>]*>(.*?)</project>`,
		"properties": `(?m)^([\w.-]+)\s*[=:]\s*(.*)$`,
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

func (e *JavaExtractor) Cleanup() {
	e.patterns = nil
}

func (e *JavaExtractor) ShouldProcess(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".java") ||
		strings.Contains(strings.ToLower(filename), "pom.xml")
}

func (e *JavaExtractor) Extract(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract package declaration
	packageName := ""
	if match := e.patterns["package"].FindStringSubmatch(content); match != nil {
		packageName = match[1]
	}

	// Extract imports
	var imports []string
	for _, match := range e.patterns["import"].FindAllStringSubmatch(content, -1) {
		imports = append(imports, match[1])
	}

	// Extract annotations
	var annotations []string
	for _, match := range e.patterns["annotation"].FindAllStringSubmatch(content, -1) {
		annotations = append(annotations, match[1])
	}

	// Extract classes
	blocks = append(blocks, e.extractType(content, "class", packageName, imports, annotations)...)

	// Extract interfaces
	blocks = append(blocks, e.extractType(content, "interface", packageName, imports, annotations)...)

	// Extract enums
	blocks = append(blocks, e.extractType(content, "enum", packageName, imports, annotations)...)

	// Extract Maven POM if present
	if strings.Contains(content, "<project") {
		if match := e.patterns["maven_pom"].FindStringSubmatch(content); match != nil {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  match[0],
				Type:     "maven-pom",
				Package:  "",
				Filename: "pom.xml",
				Language: "xml",
			})
		}
	}

	// Extract properties
	if props := e.extractProperties(content); len(props) > 0 {
		blocks = append(blocks, extractor.CodeBlock{
			Content:  props,
			Type:     "properties",
			Package:  "",
			Filename: "application.properties",
			Language: "properties",
		})
	}

	return blocks
}

func (e *JavaExtractor) extractType(content, typ, pkg string, imports, annotations []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	pattern := e.patterns[typ]

	for _, match := range pattern.FindAllStringSubmatch(content, -1) {
		if len(match) < 3 {
			continue
		}

		modifiers := []string{}
		if match[1] != "" {
			modifiers = strings.Fields(match[1])
		}

		typeName := match[2]
		body := match[3]

		blocks = append(blocks, extractor.CodeBlock{
			Content:     e.constructJavaCode(typ, typeName, body, pkg, imports),
			Type:        typ,
			Package:     pkg,
			Filename:    typeName + ".java",
			Language:    "java",
			Imports:     imports,
			Annotations: annotations,
			Modifiers:   modifiers,
		})
	}

	return blocks
}

func (e *JavaExtractor) constructJavaCode(typ, name, body, pkg string, imports []string) string {
	var code strings.Builder

	if pkg != "" {
		code.WriteString("package " + pkg + ";\n\n")
	}

	for _, imp := range imports {
		code.WriteString("import " + imp + ";\n")
	}

	if len(imports) > 0 {
		code.WriteString("\n")
	}

	code.WriteString("public " + typ + " " + name + " {\n")
	code.WriteString(body + "\n")
	code.WriteString("}")

	return code.String()
}

func (e *JavaExtractor) extractProperties(content string) string {
	var props strings.Builder

	for _, match := range e.patterns["properties"].FindAllStringSubmatch(content, -1) {
		if len(match) == 3 {
			props.WriteString(match[1] + "=" + match[2] + "\n")
		}
	}

	return props.String()
}

// GoExtractor extracts Go code
type GoExtractor struct {
	patterns map[string]*regexp.Regexp
}

func NewGoExtractor() *GoExtractor {
	return &GoExtractor{}
}

func (e *GoExtractor) Name() string { return "go" }

func (e *GoExtractor) Extensions() []string {
	return []string{".go"}
}

func (e *GoExtractor) Initialize() error {
	e.patterns = make(map[string]*regexp.Regexp)

	patterns := map[string]string{
		"package":   `package\s+(\w+)`,
		"import":    `import\s+(?:\(([^)]+)\)|"([^"]+)")`,
		"func":      `func\s+(?:\([^)]+\)\s+)?(\w+)\s*\([^)]*\)(?:\s+\([^)]*\))?\s*(?:\{[^}]*\})?`,
		"struct":    `type\s+(\w+)\s+struct\s*\{[^}]*\}`,
		"interface": `type\s+(\w+)\s+interface\s*\{[^}]*\}`,
		"const":     `const\s*\(([^)]+)\)`,
		"var":       `var\s*\(([^)]+)\)`,
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

func (e *GoExtractor) Cleanup() {
	e.patterns = nil
}

func (e *GoExtractor) ShouldProcess(filename string) bool {
	return strings.HasSuffix(filename, ".go")
}

func (e *GoExtractor) Extract(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract package name
	packageName := ""
	if match := e.patterns["package"].FindStringSubmatch(content); match != nil {
		packageName = match[1]
	}

	// Extract imports
	var imports []string
	if match := e.patterns["import"].FindStringSubmatch(content); match != nil {
		if match[1] != "" {
			// Multi-line imports
			importLines := strings.Split(match[1], "\n")
			for _, line := range importLines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, `"`) && strings.HasSuffix(line, `"`) {
					imports = append(imports, strings.Trim(line, `"`))
				}
			}
		} else if match[2] != "" {
			// Single import
			imports = append(imports, strings.Trim(match[2], `"`))
		}
	}

	// Extract functions
	for _, match := range e.patterns["func"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractFunction(content, funcName),
				Type:     "function",
				Package:  packageName,
				Filename: funcName + ".go",
				Language: "go",
				Imports:  imports,
			})
		}
	}

	// Extract structs
	for _, match := range e.patterns["struct"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			structName := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractStruct(content, structName),
				Type:     "struct",
				Package:  packageName,
				Filename: structName + ".go",
				Language: "go",
				Imports:  imports,
			})
		}
	}

	// Extract interfaces
	for _, match := range e.patterns["interface"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			interfaceName := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractInterface(content, interfaceName),
				Type:     "interface",
				Package:  packageName,
				Filename: interfaceName + ".go",
				Language: "go",
				Imports:  imports,
			})
		}
	}

	// Extract constants
	if match := e.patterns["const"].FindStringSubmatch(content); match != nil {
		blocks = append(blocks, extractor.CodeBlock{
			Content:  e.extractConstant(content, match[1]),
			Type:     "constant",
			Package:  packageName,
			Filename: match[1] + ".go",
			Language: "go",
			Imports:  imports,
		})
	}

	return blocks
}

func (e *GoExtractor) extractFunction(content, funcName string) string {
	// Find the function definition and its body
	pattern := regexp.MustCompile(`func\s+(?:\([^)]+\)\s+)?` + regexp.QuoteMeta(funcName) + `\s*\([^)]*\)(?:\s+\([^)]*\))?\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	// Fallback: return just the signature
	return "func " + funcName + "() {\n    // Function body\n}"
}

func (e *GoExtractor) extractStruct(content, structName string) string {
	pattern := regexp.MustCompile(`type\s+` + regexp.QuoteMeta(structName) + `\s+struct\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "type " + structName + " struct {\n    // Fields\n}"
}

func (e *GoExtractor) extractInterface(content, interfaceName string) string {
	pattern := regexp.MustCompile(`type\s+` + regexp.QuoteMeta(interfaceName) + `\s+interface\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "type " + interfaceName + " interface {\n    // Methods\n}"
}

func (e *GoExtractor) extractConstant(content, constName string) string {
	// Find the constant definition in the content
	pattern := regexp.MustCompile(`const\s+` + regexp.QuoteMeta(constName) + `\s+=\s+[^;\n]+`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	// Fallback: return a basic constant declaration
	return "const " + constName + " = /* value */"
}

// PythonExtractor extracts Python code
type PythonExtractor struct {
	patterns map[string]*regexp.Regexp
}

func NewPythonExtractor() *PythonExtractor {
	return &PythonExtractor{}
}

func (e *PythonExtractor) Name() string { return "python" }

func (e *PythonExtractor) Extensions() []string {
	return []string{".py", ".pyw"}
}

func (e *PythonExtractor) Initialize() error {
	e.patterns = make(map[string]*regexp.Regexp)

	patterns := map[string]string{
		"import":         `^(?:from\s+(\w+)\s+import|import\s+(\w+))`,
		"class":          `class\s+(\w+)(?:\([^)]*\))?:`,
		"function":       `def\s+(\w+)\([^)]*\):`,
		"async_function": `async\s+def\s+(\w+)\([^)]*\):`,
		"decorator":      `^@(\w+)`,
		"requirements":   `^([\w-]+)(?:[<>=!~]+[\d.,*]+)?`,
		"setup_py":       `from setuptools import setup.*?setup\(`,
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

func (e *PythonExtractor) Cleanup() {
	e.patterns = nil
}

func (e *PythonExtractor) ShouldProcess(filename string) bool {
	return strings.HasSuffix(filename, ".py") ||
		strings.HasSuffix(filename, ".pyw") ||
		strings.HasSuffix(filename, "requirements.txt") ||
		filename == "setup.py"
}

func (e *PythonExtractor) Extract(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract imports
	var imports []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if match := e.patterns["import"].FindStringSubmatch(line); match != nil {
			if match[1] != "" {
				imports = append(imports, match[1])
			} else if match[2] != "" {
				imports = append(imports, match[2])
			}
		}
	}

	// Extract classes
	for _, match := range e.patterns["class"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			className := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractClass(content, className),
				Type:     "class",
				Package:  "",
				Filename: className + ".py",
				Language: "python",
				Imports:  imports,
			})
		}
	}

	// Extract functions
	for _, match := range e.patterns["function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractFunction(content, funcName),
				Type:     "function",
				Package:  "",
				Filename: funcName + ".py",
				Language: "python",
				Imports:  imports,
			})
		}
	}

	// Extract async functions
	for _, match := range e.patterns["async_function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractAsyncFunction(content, funcName),
				Type:     "async_function",
				Package:  "",
				Filename: funcName + ".py",
				Language: "python",
				Imports:  imports,
			})
		}
	}

	// Extract requirements.txt if present
	if strings.Contains(content, "requirements.txt") ||
		(len(lines) > 0 && strings.Contains(strings.ToLower(content), "requirement")) {
		if reqs := e.extractRequirements(content); reqs != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  reqs,
				Type:     "requirements",
				Package:  "",
				Filename: "requirements.txt",
				Language: "text",
			})
		}
	}

	// Extract setup.py if present
	if match := e.patterns["setup_py"].FindString(content); match != "" {
		blocks = append(blocks, extractor.CodeBlock{
			Content:  e.extractSetupPy(content),
			Type:     "setup",
			Package:  "",
			Filename: "setup.py",
			Language: "python",
		})
	}

	return blocks
}

func (e *PythonExtractor) extractClass(content, className string) string {
	// Find class definition and its body
	pattern := regexp.MustCompile(`class\s+` + regexp.QuoteMeta(className) + `(?:\([^)]*\))?:\s*\n(?:[ \t]+[^\n]*\n)*`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "class " + className + ":\n    pass\n"
}

func (e *PythonExtractor) extractFunction(content, funcName string) string {
	pattern := regexp.MustCompile(`def\s+` + regexp.QuoteMeta(funcName) + `\([^)]*\):\s*\n(?:[ \t]+[^\n]*\n)*`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "def " + funcName + "():\n    pass\n"
}

func (e *PythonExtractor) extractAsyncFunction(content, funcName string) string {
	pattern := regexp.MustCompile(`async\s+def\s+` + regexp.QuoteMeta(funcName) + `\([^)]*\):\s*\n(?:[ \t]+[^\n]*\n)*`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "async def " + funcName + "():\n    pass\n"
}

func (e *PythonExtractor) extractRequirements(content string) string {
	var reqs strings.Builder
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if match := e.patterns["requirements"].FindStringSubmatch(line); match != nil {
			reqs.WriteString(line + "\n")
		}
	}

	return reqs.String()
}

func (e *PythonExtractor) extractSetupPy(content string) string {
	// Find setup() call and surrounding context
	pattern := regexp.MustCompile(`(?s)from setuptools import setup.*?setup\(.*?\)`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	// Fallback minimal setup.py
	return `from setuptools import setup, find_packages

setup(
    name="package",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[]
)`
}

// JavaScriptExtractor extracts JavaScript/TypeScript code
type JavaScriptExtractor struct {
	patterns map[string]*regexp.Regexp
}

func NewJavaScriptExtractor() *JavaScriptExtractor {
	return &JavaScriptExtractor{}
}

func (e *JavaScriptExtractor) Name() string { return "javascript" }

func (e *JavaScriptExtractor) Extensions() []string {
	return []string{".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs"}
}

func (e *JavaScriptExtractor) Initialize() error {
	e.patterns = make(map[string]*regexp.Regexp)

	patterns := map[string]string{
		"import":          `(?:import|require)\(?['"]([^'"]+)['"]\)?`,
		"export":          `export\s+(?:default\s+)?(?:class|function|const|let|var|async\s+function)\s+(\w+)`,
		"class":           `(?:export\s+)?(?:default\s+)?class\s+(\w+)(?:.*?)\{`,
		"function":        `(?:export\s+)?(?:default\s+)?(?:async\s+)?function\s+(\w+)\s*\(`,
		"arrow_function":  `(?:export\s+)?(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s+)?\([^)]*\)\s*=>`,
		"react_component": `(?:export\s+)?(?:default\s+)?(?:function\s+)?(\w+)\s*\([^)]*\)\s*\{.*?\n\s*return\s*\(`,
		"package_json":    `(?s)\{\s*"name"\s*:\s*"[^"]+"`,
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

func (e *JavaScriptExtractor) Cleanup() {
	e.patterns = nil
}

func (e *JavaScriptExtractor) ShouldProcess(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".js" || ext == ".jsx" || ext == ".ts" || ext == ".tsx" ||
		ext == ".mjs" || ext == ".cjs" || filename == "package.json"
}

func (e *JavaScriptExtractor) Extract(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract imports
	var imports []string
	for _, match := range e.patterns["import"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			imports = append(imports, match[1])
		}
	}

	// Extract classes
	for _, match := range e.patterns["class"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			className := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractClass(content, className),
				Type:     "class",
				Package:  "",
				Filename: className + ".js",
				Language: "javascript",
				Imports:  imports,
			})
		}
	}

	// Extract functions
	for _, match := range e.patterns["function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractFunction(content, funcName),
				Type:     "function",
				Package:  "",
				Filename: funcName + ".js",
				Language: "javascript",
				Imports:  imports,
			})
		}
	}

	// Extract arrow functions
	for _, match := range e.patterns["arrow_function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractArrowFunction(content, funcName),
				Type:     "arrow_function",
				Package:  "",
				Filename: funcName + ".js",
				Language: "javascript",
				Imports:  imports,
			})
		}
	}

	// Extract React components
	for _, match := range e.patterns["react_component"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			componentName := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractReactComponent(content, componentName),
				Type:     "react_component",
				Package:  "",
				Filename: componentName + ".jsx",
				Language: "javascript",
				Imports:  imports,
			})
		}
	}

	// Extract package.json if present
	if strings.Contains(content, `"name"`) && strings.Contains(content, `"version"`) {
		if pkg := e.extractPackageJson(content); pkg != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  pkg,
				Type:     "package_json",
				Package:  "",
				Filename: "package.json",
				Language: "json",
			})
		}
	}

	return blocks
}

func (e *JavaScriptExtractor) extractClass(content, className string) string {
	// Find class definition
	pattern := regexp.MustCompile(`(?:export\s+)?(?:default\s+)?class\s+` +
		regexp.QuoteMeta(className) + `\s*(?:extends\s+\w+)?\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "class " + className + " {\n    constructor() {\n        // constructor\n    }\n}"
}

func (e *JavaScriptExtractor) extractFunction(content, funcName string) string {
	pattern := regexp.MustCompile(`(?:export\s+)?(?:default\s+)?(?:async\s+)?function\s+` +
		regexp.QuoteMeta(funcName) + `\s*\([^)]*\)\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "function " + funcName + "() {\n    // function body\n}"
}

func (e *JavaScriptExtractor) extractArrowFunction(content, funcName string) string {
	pattern := regexp.MustCompile(`(?:export\s+)?(?:const|let|var)\s+` +
		regexp.QuoteMeta(funcName) + `\s*=\s*(?:async\s+)?\([^)]*\)\s*=>\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "const " + funcName + " = () => {\n    // arrow function body\n}"
}

func (e *JavaScriptExtractor) extractReactComponent(content, componentName string) string {
	// Find React component with JSX
	pattern := regexp.MustCompile(`(?:export\s+)?(?:default\s+)?(?:function\s+)?` +
		regexp.QuoteMeta(componentName) + `\s*\([^)]*\)\s*\{.*?\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "function " + componentName + "() {\n    return (\n        <div>\n            {/* JSX content */}\n        </div>\n    );\n}"
}

func (e *JavaScriptExtractor) extractPackageJson(content string) string {
	// Try to extract JSON object
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start != -1 && end != -1 && end > start {
		jsonStr := content[start : end+1]
		// Validate it looks like package.json
		if strings.Contains(jsonStr, `"name"`) && strings.Contains(jsonStr, `"version"`) {
			return jsonStr
		}
	}

	// Fallback minimal package.json
	return `{
  "name": "package",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
  "scripts": {
    "start": "node index.js"
  },
  "dependencies": {}
}`
}

// RustExtractor extracts Rust code
type RustExtractor struct {
	patterns map[string]*regexp.Regexp
}

func NewRustExtractor() *RustExtractor {
	return &RustExtractor{}
}

func (e *RustExtractor) Name() string { return "rust" }

func (e *RustExtractor) Extensions() []string {
	return []string{".rs"}
}

func (e *RustExtractor) Initialize() error {
	e.patterns = make(map[string]*regexp.Regexp)

	patterns := map[string]string{
		"module":      `mod\s+(\w+)\s*\{`,
		"function":    `fn\s+(\w+)\s*\([^)]*\)(?:\s*->\s*[^{]+)?\s*\{`,
		"struct":      `struct\s+(\w+)\s*\{`,
		"enum":        `enum\s+(\w+)\s*\{`,
		"trait":       `trait\s+(\w+)\s*\{`,
		"impl":        `impl\s+(?:<[^>]+>\s+)?(\w+)\s*\{`,
		"use":         `use\s+([\w::]+);`,
		"macro_rules": `macro_rules!\s+(\w+)\s*\{`,
		"crate_toml":  `\[package\][\s\S]*?name\s*=\s*"[^"]+"`,
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

func (e *RustExtractor) Cleanup() {
	e.patterns = nil
}

func (e *RustExtractor) ShouldProcess(filename string) bool {
	return strings.HasSuffix(filename, ".rs") || strings.Contains(filename, "Cargo.toml")
}

func (e *RustExtractor) Extract(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract modules
	for _, match := range e.patterns["module"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			moduleName := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractModule(content, moduleName),
				Type:     "module",
				Package:  "",
				Filename: moduleName + ".rs",
				Language: "rust",
			})
		}
	}

	// Extract functions
	for _, match := range e.patterns["function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractFunction(content, funcName),
				Type:     "function",
				Package:  "",
				Filename: funcName + ".rs",
				Language: "rust",
			})
		}
	}

	// Extract structs
	for _, match := range e.patterns["struct"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			structName := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractStruct(content, structName),
				Type:     "struct",
				Package:  "",
				Filename: structName + ".rs",
				Language: "rust",
			})
		}
	}

	// Extract enums
	for _, match := range e.patterns["enum"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			enumName := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractEnum(content, enumName),
				Type:     "enum",
				Package:  "",
				Filename: enumName + ".rs",
				Language: "rust",
			})
		}
	}

	// Extract Cargo.toml if present
	if strings.Contains(content, "[package]") {
		if cargo := e.extractCargoToml(content); cargo != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  cargo,
				Type:     "cargo_toml",
				Package:  "",
				Filename: "Cargo.toml",
				Language: "toml",
			})
		}
	}

	return blocks
}

func (e *RustExtractor) extractModule(content, moduleName string) string {
	pattern := regexp.MustCompile(`mod\s+` + regexp.QuoteMeta(moduleName) + `\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "mod " + moduleName + " {\n    // Module content\n}"
}

func (e *RustExtractor) extractFunction(content, funcName string) string {
	pattern := regexp.MustCompile(`fn\s+` + regexp.QuoteMeta(funcName) + `\s*\([^)]*\)(?:\s*->\s*[^{]+)?\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "fn " + funcName + "() {\n    // Function body\n}"
}

func (e *RustExtractor) extractStruct(content, structName string) string {
	pattern := regexp.MustCompile(`struct\s+` + regexp.QuoteMeta(structName) + `\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "struct " + structName + " {\n    // Fields\n}"
}

func (e *RustExtractor) extractEnum(content, enumName string) string {
	pattern := regexp.MustCompile(`enum\s+` + regexp.QuoteMeta(enumName) + `\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "enum " + enumName + " {\n    // Variants\n}"
}

func (e *RustExtractor) extractCargoToml(content string) string {
	// Look for Cargo.toml content
	pattern := regexp.MustCompile(`(?s)\[package\].*?name\s*=\s*"[^"]+".*?version\s*=\s*"[^"]+".*?`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	// Fallback minimal Cargo.toml
	return `[package]
name = "my_package"
version = "0.1.0"
edition = "2021"

[dependencies]
`
}

// DartExtractor extracts Dart code
type DartExtractor struct {
	patterns map[string]*regexp.Regexp
}

func NewDartExtractor() *DartExtractor {
	return &DartExtractor{}
}

func (e *DartExtractor) Name() string { return "dart" }

func (e *DartExtractor) Extensions() []string {
	return []string{".dart"}
}

func (e *DartExtractor) Initialize() error {
	e.patterns = make(map[string]*regexp.Regexp)

	patterns := map[string]string{
		"library":    `library\s+([\w.]+);`,
		"import":     `import\s+['"]([^'"]+)['"];`,
		"class":      `class\s+(\w+)(?:\s+extends\s+\w+)?(?:\s+with\s+\w+)?(?:\s+implements\s+[^{]+)?\s*\{`,
		"function":   `(\w+)?\s*(\w+)\s*\([^)]*\)\s*\{`,
		"typedef":    `typedef\s+(\w+)\s*=\s*`,
		"mixin":      `mixin\s+(\w+)\s*\{`,
		"enum":       `enum\s+(\w+)\s*\{`,
		"pubspec":    `name:\s+[\w-]+`,
		"main":       `void\s+main\s*\([^)]*\)\s*\{`,
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

func (e *DartExtractor) Cleanup() {
	e.patterns = nil
}

func (e *DartExtractor) ShouldProcess(filename string) bool {
	return strings.HasSuffix(filename, ".dart") || strings.Contains(filename, "pubspec.yaml")
}

func (e *DartExtractor) Extract(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract library declarations
	if match := e.patterns["library"].FindStringSubmatch(content); match != nil {
		blocks = append(blocks, extractor.CodeBlock{
			Content:  match[0],
			Type:     "library",
			Package:  match[1],
			Filename: "library.dart",
			Language: "dart",
		})
	}

	// Extract imports
	var imports []string
	for _, match := range e.patterns["import"].FindAllStringSubmatch(content, -1) {
		imports = append(imports, match[1])
	}

	// Extract classes
	for _, match := range e.patterns["class"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			className := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  e.extractClass(content, className),
				Type:     "class",
				Package:  "",
				Filename: className + ".dart",
				Language: "dart",
				Imports:  imports,
			})
		}
	}

	// Extract functions
	for _, match := range e.patterns["function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 2 {
			// Skip if first group is a type (like 'int' or 'void')
			if match[1] != "void" && match[1] != "bool" && match[1] != "String" && match[1] != "int" && match[1] != "double" {
				funcName := match[2]
				blocks = append(blocks, extractor.CodeBlock{
					Content:  e.extractFunction(content, funcName),
					Type:     "function",
					Package:  "",
					Filename: funcName + ".dart",
					Language: "dart",
					Imports:  imports,
				})
			}
		}
	}

	// Extract pubspec.yaml if present
	if strings.Contains(content, "name:") && strings.Contains(content, "dependencies:") {
		if pubspec := e.extractPubspec(content); pubspec != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  pubspec,
				Type:     "pubspec",
				Package:  "",
				Filename: "pubspec.yaml",
				Language: "yaml",
			})
		}
	}

	return blocks
}

func (e *DartExtractor) extractClass(content, className string) string {
	pattern := regexp.MustCompile(`class\s+` + regexp.QuoteMeta(className) + `(?:\s+extends\s+\w+)?(?:\s+with\s+\w+)?(?:\s+implements\s+[^{]+)?\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "class " + className + " {\n  // Class body\n}"
}

func (e *DartExtractor) extractFunction(content, funcName string) string {
	pattern := regexp.MustCompile(`(?:\w+\s+)?` + regexp.QuoteMeta(funcName) + `\s*\([^)]*\)\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return funcName + "() {\n  // Function body\n}"
}

func (e *DartExtractor) extractPubspec(content string) string {
	// Look for pubspec.yaml content
	pattern := regexp.MustCompile(`(?s)name:\s+[\w-]+.*?dependencies:\s*\n(?:\s+[\w-]+\s*:.*)*`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	// Fallback minimal pubspec.yaml
	return `name: my_app
description: A new Flutter project.
version: 1.0.0

environment:
  sdk: '>=2.19.0 <4.0.0'

dependencies:
  flutter:
    sdk: flutter

dev_dependencies:
  flutter_test:
    sdk: flutter
`
}

// GenericExtractor extracts generic code blocks from any text
type GenericExtractor struct {
	patterns map[string]*regexp.Regexp
}

func NewGenericExtractor() *GenericExtractor {
	return &GenericExtractor{}
}

func (e *GenericExtractor) Name() string { return "generic" }

func (e *GenericExtractor) Extensions() []string {
	return []string{".txt", ".md", ".rst", ".yml", ".yaml", ".xml", ".json", ".ini", ".cfg", ".conf"}
}

func (e *GenericExtractor) Initialize() error {
	e.patterns = make(map[string]*regexp.Regexp)

	patterns := map[string]string{
		"code_block":      "```(\\w+)\\n([\\s\\S]*?)```",
		"markdown_header": "^#{1,6}\\s+(.+)$",
		"yaml_block":      "(?:^|\\n)([\\w-]+:\\s*(?:[^\\n]+|(?:\\n(?:  |\\t)+[^\\n]+)+))",
		"json_block":      "\\{[\\s\\S]*?\\}",
		"xml_block":       "<[\\w]+[^>]*>[\\s\\S]*?</[\\w]+>",
		"ini_block":       "(?:^|\\n)(\\[[^\\]]+\\]\\n(?:[^\\[\\n].*\\n)*)",
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

func (e *GenericExtractor) Cleanup() {
	e.patterns = nil
}

func (e *GenericExtractor) ShouldProcess(filename string) bool {
	// Generic extractor should process any file as fallback
	return true
}

func (e *GenericExtractor) Extract(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract code blocks from markdown
	for _, match := range e.patterns["code_block"].FindAllStringSubmatch(content, -1) {
		if len(match) > 2 {
			language := match[1]
			code := match[2]

			blocks = append(blocks, extractor.CodeBlock{
				Content:  code,
				Type:     "code_block",
				Package:  "",
				Filename: "code." + strings.ToLower(language),
				Language: language,
			})
		}
	}

	// Extract YAML blocks
	for _, match := range e.patterns["yaml_block"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			yamlContent := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  yamlContent,
				Type:     "yaml",
				Package:  "",
				Filename: "config.yaml",
				Language: "yaml",
			})
		}
	}

	// Extract JSON blocks
	for _, match := range e.patterns["json_block"].FindAllStringSubmatch(content, -1) {
		jsonContent := match[0]
		// Validate JSON
		if strings.Contains(jsonContent, "{") && strings.Contains(jsonContent, "}") {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  jsonContent,
				Type:     "json",
				Package:  "",
				Filename: "config.json",
				Language: "json",
			})
		}
	}

	// Extract XML blocks
	for _, match := range e.patterns["xml_block"].FindAllStringSubmatch(content, -1) {
		xmlContent := match[0]
		blocks = append(blocks, extractor.CodeBlock{
			Content:  xmlContent,
			Type:     "xml",
			Package:  "",
			Filename: "config.xml",
			Language: "xml",
		})
	}

	// Extract INI blocks
	for _, match := range e.patterns["ini_block"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			iniContent := match[1]
			blocks = append(blocks, extractor.CodeBlock{
				Content:  iniContent,
				Type:     "ini",
				Package:  "",
				Filename: "config.ini",
				Language: "ini",
			})
		}
	}

	// If no specific blocks found, treat the entire content as a block
	if len(blocks) == 0 && strings.TrimSpace(content) != "" {
		// Try to determine file type from content
		filename, language := e.determineFileType(content)

		blocks = append(blocks, extractor.CodeBlock{
			Content:  content,
			Type:     "text",
			Package:  "",
			Filename: filename,
			Language: language,
		})
	}

	return blocks
}

func (e *GenericExtractor) determineFileType(content string) (string, string) {
	// Check for common patterns
	if strings.HasPrefix(strings.TrimSpace(content), "{") &&
		strings.HasSuffix(strings.TrimSpace(content), "}") {
		return "data.json", "json"
	}

	if strings.Contains(content, "<?xml") {
		return "data.xml", "xml"
	}

	if strings.Contains(content, "---") &&
		(strings.Contains(content, ":") || strings.Contains(content, "- ")) {
		return "data.yaml", "yaml"
	}

	if strings.Contains(content, "[") && strings.Contains(content, "]") &&
		strings.Contains(content, "=") {
		return "config.ini", "ini"
	}

	// Default to text
	return "content.txt", "text"
}