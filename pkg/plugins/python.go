// python_extractor.go - Python Code Extractor
package plugins

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	
	"github.com/bhangun/coto/pkg/extractor"
)

// PythonExtractor extracts Python code
type PythonExtractor struct {
	patterns map[string]*regexp.Regexp
}

// NewPythonExtractor creates a new Python extractor
func NewPythonExtractor() *PythonExtractor {
	return &PythonExtractor{}
}

// Name returns the extractor name
func (e *PythonExtractor) Name() string {
	return "python"
}

// Extensions returns supported file extensions
func (e *PythonExtractor) Extensions() []string {
	return []string{".py", ".pyw", ".pyi", ".pyx", ".pxd", ".pxi"}
}

// Initialize sets up regex patterns
func (e *PythonExtractor) Initialize() error {
	e.patterns = make(map[string]*regexp.Regexp)

	patterns := map[string]string{
		// Module-level patterns
		"import": `^(?:from\s+([\w.]+)\s+import\s+([\w*, ]+)|import\s+([\w., ]+))`,
		"import_as": `import\s+([\w.]+)\s+as\s+(\w+)`,
		"from_import": `from\s+([\w.]+)\s+import\s+\(([^)]+)\)`,

		// Class patterns
		"class": `^class\s+(\w+)(?:\(([^)]*)\))?\s*:`,
		"dataclass": `^@dataclass\s*\nclass\s+(\w+)`,
		"pydantic_model": `^class\s+(\w+)\((?:BaseModel|pydantic\.BaseModel)\)`,

		// Function patterns
		"function": `^def\s+(\w+)\(([^)]*)\)(?:\s*->\s*([^:]+))?\s*:`,
		"async_function": `^async\s+def\s+(\w+)\(([^)]*)\)(?:\s*->\s*([^:]+))?\s*:`,
		"lambda": `lambda\s+([^:]+):`,

		// Method patterns
		"method": `^\s+def\s+(\w+)\(([^)]*)\)(?:\s*->\s*([^:]+))?\s*:`,
		"classmethod": `^\s+@classmethod\s*\n\s+def\s+(\w+)\(cls[^)]*\)`,
		"staticmethod": `^\s+@staticmethod\s*\n\s+def\s+(\w+)\([^)]*\)`,
		"property_decorator": `^\s+@property\s*\n\s+def\s+(\w+)\(self[^)]*\)`,

		// Special methods
		"init_method": `^\s+def\s+__init__\(([^)]*)\)\s*:`,
		"str_method": `^\s+def\s+__str__\(([^)]*)\)\s*:`,
		"repr_method": `^\s+def\s+__repr__\(([^)]*)\)\s*:`,

		// Decorator patterns
		"decorator": `^@(\w+)(?:\(([^)]*)\))?`,
		"multiple_decorators": `^(?:@[\w.]+(?:\([^)]*\))?\s*\n)+`,

		// Type hints
		"type_hint": `:\s*([\w\[\], \.]+)(?:\s*=\s*[^,\n]+)?`,
		"return_hint": `->\s*([\w\[\], \.]+)`,

		// Configuration files
		"requirements": `^([\w\-\[\]]+)(?:[<>=!~]+[\d.,*]+)?(?:\s*#.*)?$`,
		"setup_py": `setup\s*\(`,
		"pyproject_toml": `\[tool\.(?:poetry|flit|setuptools)\]`,
		"setup_cfg": `^\[([\w:]+)\]`,

		// Django/Flask specific
		"django_model": `class\s+(\w+)\(models\.Model\)`,
		"django_view": `def\s+(\w+)\(request[^)]*\)`,
		"flask_route": `@app\.route\(['"]([^'"]+)['"]\)`,

		// FastAPI specific
		"fastapi_route": `@(?:app|router)\.(?:get|post|put|delete|patch|options|head)\(['"]([^'"]+)['"]\)`,

		// Test patterns
		"test_class": `class\s+(Test\w+)\([^)]*\)`,
		"test_function": `def\s+(test_\w+)\([^)]*\)`,
		"pytest_fixture": `@pytest\.fixture`,

		// Docstring patterns
		"docstring_triple_single": `'''([^']*?)'''`,
		"docstring_triple_double": `"""([^"]*?)"""`,

		// Comment patterns
		"shebang": `^#!.*python`,
		"encoding": `^#.*coding[:=]\s*([-\w.]+)`,

		// Async patterns
		"async_for": `async\s+for\s+(\w+)\s+in`,
		"async_with": `async\s+with\s+`,
		"await_expr": `await\s+(\w+)`,
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
func (e *PythonExtractor) Cleanup() {
	e.patterns = nil
}

// ShouldProcess checks if this extractor should handle the file
func (e *PythonExtractor) ShouldProcess(filename string) bool {
	lowerName := strings.ToLower(filename)

	// Check by extension
	ext := filepath.Ext(lowerName)
	if ext == ".py" || ext == ".pyw" || ext == ".pyi" ||
	   ext == ".pyx" || ext == ".pxd" || ext == ".pxi" {
		return true
	}

	// Check by filename
	if filename == "requirements.txt" || filename == "setup.py" ||
	   filename == "pyproject.toml" || filename == "setup.cfg" ||
	   filename == "MANIFEST.in" || filename == "Pipfile" ||
	   filename == "tox.ini" || filename == "pytest.ini" ||
	   filename == "mypy.ini" || filename == ".python-version" {
		return true
	}

	// Check for Python shebang in content (if we had content)
	return false
}

// Extract extracts Python code blocks from content
func (e *PythonExtractor) Extract(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract module-level imports
	imports := e.extractImports(content)

	// Extract Python-specific blocks
	blocks = append(blocks, e.extractClasses(content, imports)...)
	blocks = append(blocks, e.extractFunctions(content, imports)...)
	blocks = append(blocks, e.extractMethods(content, imports)...)
	blocks = append(blocks, e.extractDecoratedFunctions(content, imports)...)
	blocks = append(blocks, e.extractAsyncFunctions(content, imports)...)
	blocks = append(blocks, e.extractTestCases(content, imports)...)

	// Extract configuration files
	blocks = append(blocks, e.extractRequirements(content)...)
	blocks = append(blocks, e.extractSetupPy(content)...)
	blocks = append(blocks, e.extractPyprojectToml(content)...)
	blocks = append(blocks, e.extractConfigFiles(content)...)

	// Extract web framework specific code
	blocks = append(blocks, e.extractWebFrameworkCode(content, imports)...)

	// If no specific blocks found, extract the entire module
	if len(blocks) == 0 && strings.TrimSpace(content) != "" {
		blocks = append(blocks, e.extractModule(content, imports)...)
	}

	return blocks
}

// extractImports extracts import statements
func (e *PythonExtractor) extractImports(content string) []string {
	var imports []string

	// Extract simple imports
	for _, match := range e.patterns["import"].FindAllStringSubmatch(content, -1) {
		if len(match) > 3 {
			if match[1] != "" { // from X import Y
				module := match[1]
				items := strings.Split(match[2], ",")
				for _, item := range items {
					item = strings.TrimSpace(item)
					if item != "" {
						imports = append(imports, fmt.Sprintf("%s.%s", module, item))
					}
				}
			} else if match[3] != "" { // import X, Y, Z
				modules := strings.Split(match[3], ",")
				for _, module := range modules {
					module = strings.TrimSpace(module)
					if module != "" {
						imports = append(imports, module)
					}
				}
			}
		}
	}

	// Extract import with aliases
	for _, match := range e.patterns["import_as"].FindAllStringSubmatch(content, -1) {
		if len(match) > 2 {
			imports = append(imports, fmt.Sprintf("%s as %s", match[1], match[2]))
		}
	}

	// Extract multi-line from imports
	for _, match := range e.patterns["from_import"].FindAllStringSubmatch(content, -1) {
		if len(match) > 2 {
			module := match[1]
			items := strings.Split(match[2], "\n")
			for _, item := range items {
				item = strings.TrimSpace(item)
				if item != "" && !strings.HasPrefix(item, "#") {
					imports = append(imports, fmt.Sprintf("%s.%s", module, item))
				}
			}
		}
	}

	return imports
}

// extractClasses extracts Python classes
func (e *PythonExtractor) extractClasses(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract regular classes
	for _, match := range e.patterns["class"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			className := match[1]
			inheritance := ""
			if len(match) > 2 && match[2] != "" {
				inheritance = match[2]
			}

			classContent := e.extractClassBody(content, className, inheritance)

			// Check for decorators before class
			decorators := e.extractClassDecorators(content, className)

			blocks = append(blocks, extractor.CodeBlock{
				Content:     classContent,
				Type:        "class",
				Package:     e.extractModuleName(content),
				Filename:    e.classToFilename(className),
				Language:    "python",
				Imports:     imports,
				Annotations: decorators,
			})
		}
	}

	// Extract dataclasses
	for _, match := range e.patterns["dataclass"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			className := match[1]
			classContent := e.extractDataclassBody(content, className)

			blocks = append(blocks, extractor.CodeBlock{
				Content:     classContent,
				Type:        "dataclass",
				Package:     e.extractModuleName(content),
				Filename:    e.classToFilename(className),
				Language:    "python",
				Imports:     append(imports, "dataclasses"),
				Annotations: []string{"dataclass"},
			})
		}
	}

	// Extract Pydantic models
	for _, match := range e.patterns["pydantic_model"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			className := match[1]
			modelContent := e.extractPydanticModel(content, className)

			blocks = append(blocks, extractor.CodeBlock{
				Content:     modelContent,
				Type:        "pydantic_model",
				Package:     e.extractModuleName(content),
				Filename:    e.classToFilename(className),
				Language:    "python",
				Imports:     append(imports, "pydantic"),
				Annotations: []string{"BaseModel"},
			})
		}
	}

	// Extract Django models
	for _, match := range e.patterns["django_model"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			className := match[1]
			modelContent := e.extractDjangoModel(content, className)

			blocks = append(blocks, extractor.CodeBlock{
				Content:     modelContent,
				Type:        "django_model",
				Package:     e.extractModuleName(content),
				Filename:    e.classToFilename(className),
				Language:    "python",
				Imports:     append(imports, "django.db.models"),
				Annotations: []string{"Model"},
			})
		}
	}

	return blocks
}

// extractFunctions extracts Python functions
func (e *PythonExtractor) extractFunctions(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract module-level functions
	for _, match := range e.patterns["function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			params := ""
			returnType := ""

			if len(match) > 2 {
				params = match[2]
			}
			if len(match) > 3 {
				returnType = match[3]
			}

			// Skip if it's actually a method (indented)
			funcLine := e.getLineWithPattern(content, `def\s+`+regexp.QuoteMeta(funcName))
			if strings.HasPrefix(funcLine, " ") || strings.HasPrefix(funcLine, "\t") {
				continue
			}

			funcContent := e.extractFunctionBody(content, funcName, params, returnType, false)

			blocks = append(blocks, extractor.CodeBlock{
				Content:  funcContent,
				Type:     "function",
				Package:  e.extractModuleName(content),
				Filename: e.functionToFilename(funcName),
				Language: "python",
				Imports:  imports,
			})
		}
	}

	return blocks
}

// extractMethods extracts class methods
func (e *PythonExtractor) extractMethods(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract __init__ methods
	for _, match := range e.patterns["init_method"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			params := match[1]
			funcContent := e.extractInitMethodBody(content, params)

			blocks = append(blocks, extractor.CodeBlock{
				Content:  funcContent,
				Type:     "__init__",
				Package:  e.extractModuleName(content),
				Filename: "__init__.py",
				Language: "python",
				Imports:  imports,
			})
		}
	}

	// Extract class methods
	for _, match := range e.patterns["classmethod"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			methodName := match[1]
			methodContent := e.extractClassMethodBody(content, methodName)

			blocks = append(blocks, extractor.CodeBlock{
				Content:     methodContent,
				Type:        "classmethod",
				Package:     e.extractModuleName(content),
				Filename:    e.methodToFilename(methodName),
				Language:    "python",
				Imports:     imports,
				Annotations: []string{"classmethod"},
			})
		}
	}

	// Extract static methods
	for _, match := range e.patterns["staticmethod"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			methodName := match[1]
			methodContent := e.extractStaticMethodBody(content, methodName)

			blocks = append(blocks, extractor.CodeBlock{
				Content:     methodContent,
				Type:        "staticmethod",
				Package:     e.extractModuleName(content),
				Filename:    e.methodToFilename(methodName),
				Language:    "python",
				Imports:     imports,
				Annotations: []string{"staticmethod"},
			})
		}
	}

	return blocks
}

// extractDecoratedFunctions extracts functions with decorators
func (e *PythonExtractor) extractDecoratedFunctions(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Find lines with decorators followed by function definitions
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "@") {
			// Look ahead for function definition
			for j := i + 1; j < len(lines) && j < i+5; j++ {
				if strings.Contains(lines[j], "def ") {
					// Extract function name from next line
					if match := e.patterns["function"].FindStringSubmatch(lines[j]); match != nil {
						if len(match) > 1 {
							funcName := match[1]
							funcContent := e.extractDecoratedFunction(content, funcName, line)

							// Extract decorator name
							decorator := strings.TrimSpace(strings.TrimPrefix(line, "@"))
							if idx := strings.Index(decorator, "("); idx != -1 {
								decorator = decorator[:idx]
							}

							blocks = append(blocks, extractor.CodeBlock{
								Content:     funcContent,
								Type:        "decorated_function",
								Package:     e.extractModuleName(content),
								Filename:    e.functionToFilename(funcName),
								Language:    "python",
								Imports:     imports,
								Annotations: []string{decorator},
							})
						}
					}
					break
				}
			}
		}
	}

	// Extract property methods
	for _, match := range e.patterns["property_decorator"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			methodName := match[1]
			methodContent := e.extractPropertyMethod(content, methodName)

			blocks = append(blocks, extractor.CodeBlock{
				Content:     methodContent,
				Type:        "property",
				Package:     e.extractModuleName(content),
				Filename:    e.methodToFilename(methodName),
				Language:    "python",
				Imports:     imports,
				Annotations: []string{"property"},
			})
		}
	}

	return blocks
}

// extractAsyncFunctions extracts async functions
func (e *PythonExtractor) extractAsyncFunctions(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract async functions
	for _, match := range e.patterns["async_function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			params := ""
			returnType := ""

			if len(match) > 2 {
				params = match[2]
			}
			if len(match) > 3 {
				returnType = match[3]
			}

			funcContent := e.extractFunctionBody(content, funcName, params, returnType, true)

			blocks = append(blocks, extractor.CodeBlock{
				Content:  funcContent,
				Type:     "async_function",
				Package:  e.extractModuleName(content),
				Filename: e.functionToFilename(funcName),
				Language: "python",
				Imports:  imports,
				Modifiers: []string{"async"},
			})
		}
	}

	return blocks
}

// extractTestCases extracts test functions and classes
func (e *PythonExtractor) extractTestCases(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract test classes
	for _, match := range e.patterns["test_class"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			className := match[1]
			classContent := e.extractClassBody(content, className, "unittest.TestCase")

			blocks = append(blocks, extractor.CodeBlock{
				Content:  classContent,
				Type:     "test_class",
				Package:  e.extractModuleName(content),
				Filename: e.classToFilename(className),
				Language: "python",
				Imports:  append(imports, "unittest", "pytest"),
			})
		}
	}

	// Extract test functions
	for _, match := range e.patterns["test_function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			funcContent := e.extractTestFunctionBody(content, funcName)

			blocks = append(blocks, extractor.CodeBlock{
				Content:  funcContent,
				Type:     "test_function",
				Package:  e.extractModuleName(content),
				Filename: e.functionToFilename(funcName),
				Language: "python",
				Imports:  append(imports, "pytest"),
			})
		}
	}

	// Extract pytest fixtures
	if e.patterns["pytest_fixture"].MatchString(content) {
		fixtureContent := e.extractPytestFixtures(content)
		if fixtureContent != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  fixtureContent,
				Type:     "pytest_fixture",
				Package:  e.extractModuleName(content),
				Filename: "conftest.py",
				Language: "python",
				Imports:  append(imports, "pytest"),
			})
		}
	}

	return blocks
}

// extractRequirements extracts requirements.txt
func (e *PythonExtractor) extractRequirements(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Check if this looks like requirements.txt
	if strings.Contains(content, "requirements") ||
	   (len(strings.Split(strings.TrimSpace(content), "\n")) > 3 &&
	    e.patterns["requirements"].MatchString(content)) {

		var requirements []string
		lines := strings.Split(content, "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if match := e.patterns["requirements"].FindStringSubmatch(line); match != nil {
				if len(match) > 1 {
					requirements = append(requirements, line)
				}
			}
		}

		if len(requirements) > 0 {
			reqContent := strings.Join(requirements, "\n")
			blocks = append(blocks, extractor.CodeBlock{
				Content:  reqContent,
				Type:     "requirements",
				Package:  "",
				Filename: "requirements.txt",
				Language: "text",
			})
		}
	}

	return blocks
}

// extractSetupPy extracts setup.py
func (e *PythonExtractor) extractSetupPy(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	if e.patterns["setup_py"].MatchString(content) {
		setupContent := e.extractSetupContent(content)

		if setupContent != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  setupContent,
				Type:     "setup_py",
				Package:  "",
				Filename: "setup.py",
				Language: "python",
				Imports:  []string{"setuptools"},
			})
		}
	}

	return blocks
}

// extractPyprojectToml extracts pyproject.toml
func (e *PythonExtractor) extractPyprojectToml(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	if e.patterns["pyproject_toml"].MatchString(content) {
		// Extract TOML sections
		sections := e.extractTomlSections(content)

		if len(sections) > 0 {
			tomlContent := strings.Join(sections, "\n\n")
			blocks = append(blocks, extractor.CodeBlock{
				Content:  tomlContent,
				Type:     "pyproject_toml",
				Package:  "",
				Filename: "pyproject.toml",
				Language: "toml",
			})
		}
	}

	return blocks
}

// extractConfigFiles extracts other config files
func (e *PythonExtractor) extractConfigFiles(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract setup.cfg
	if e.patterns["setup_cfg"].MatchString(content) {
		cfgSections := e.extractSetupCfgSections(content)

		if len(cfgSections) > 0 {
			cfgContent := strings.Join(cfgSections, "\n\n")
			blocks = append(blocks, extractor.CodeBlock{
				Content:  cfgContent,
				Type:     "setup_cfg",
				Package:  "",
				Filename: "setup.cfg",
				Language: "ini",
			})
		}
	}

	// Extract tox.ini
	if strings.Contains(content, "[tox]") {
		toxContent := e.extractToxIni(content)
		if toxContent != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  toxContent,
				Type:     "tox_ini",
				Package:  "",
				Filename: "tox.ini",
				Language: "ini",
			})
		}
	}

	return blocks
}

// extractWebFrameworkCode extracts web framework specific code
func (e *PythonExtractor) extractWebFrameworkCode(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	// Extract Flask routes
	for _, match := range e.patterns["flask_route"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			routePath := match[1]
			routeContent := e.extractFlaskRoute(content, routePath)

			blocks = append(blocks, extractor.CodeBlock{
				Content:  routeContent,
				Type:     "flask_route",
				Package:  e.extractModuleName(content),
				Filename: "routes.py",
				Language: "python",
				Imports:  append(imports, "flask"),
			})
		}
	}

	// Extract FastAPI routes
	for _, match := range e.patterns["fastapi_route"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			routePath := match[1]
			routeContent := e.extractFastAPIRoute(content, routePath)

			blocks = append(blocks, extractor.CodeBlock{
				Content:  routeContent,
				Type:     "fastapi_route",
				Package:  e.extractModuleName(content),
				Filename: "api.py",
				Language: "python",
				Imports:  append(imports, "fastapi"),
			})
		}
	}

	// Extract Django views
	for _, match := range e.patterns["django_view"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			viewName := match[1]
			viewContent := e.extractDjangoView(content, viewName)

			blocks = append(blocks, extractor.CodeBlock{
				Content:  viewContent,
				Type:     "django_view",
				Package:  e.extractModuleName(content),
				Filename: "views.py",
				Language: "python",
				Imports:  append(imports, "django.http"),
			})
		}
	}

	return blocks
}

// extractModule extracts entire module if no specific blocks found
func (e *PythonExtractor) extractModule(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock

	moduleName := e.extractModuleName(content)

	// Add shebang if present
	shebang := ""
	if match := e.patterns["shebang"].FindString(content); match != "" {
		shebang = match + "\n\n"
	}

	// Add encoding if present
	encoding := ""
	if match := e.patterns["encoding"].FindString(content); match != "" {
		encoding = match + "\n\n"
	}

	// Create module with imports
	moduleContent := shebang + encoding + e.reconstructImports(imports) + "\n\n" + content

	blocks = append(blocks, extractor.CodeBlock{
		Content:  moduleContent,
		Type:     "module",
		Package:  moduleName,
		Filename: moduleName + ".py",
		Language: "python",
		Imports:  imports,
	})

	return blocks
}

// Helper Methods

// extractModuleName extracts module name from content
func (e *PythonExtractor) extractModuleName(content string) string {
	// Look for module docstring or first meaningful line
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") ||
		   strings.HasPrefix(line, "\"\"\"") || strings.HasPrefix(line, "'''") {
			continue
		}

		// Check for module-level assignment
		if strings.Contains(line, "__name__") || strings.Contains(line, "__package__") {
			continue
		}

		// Use first non-import, non-empty line to guess module name
		if !strings.HasPrefix(line, "import ") && !strings.HasPrefix(line, "from ") {
			// Try to extract a name from assignment
			if strings.Contains(line, "=") {
				parts := strings.Split(line, "=")
				if len(parts) > 0 {
					name := strings.TrimSpace(parts[0])
					if !strings.Contains(name, " ") {
						return name
					}
				}
			}
		}
	}

	// Default module name
	return "module"
}

// classToFilename converts class name to filename
func (e *PythonExtractor) classToFilename(className string) string {
	// Convert CamelCase to snake_case
	var result strings.Builder
	for i, r := range className {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String()) + ".py"
}

// functionToFilename converts function name to filename
func (e *PythonExtractor) functionToFilename(funcName string) string {
	// Convert snake_case if needed, otherwise use as-is
	if strings.Contains(funcName, "_") {
		return funcName + ".py"
	}

	// Convert camelCase to snake_case
	var result strings.Builder
	for i, r := range funcName {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String()) + ".py"
}

// methodToFilename converts method name to filename
func (e *PythonExtractor) methodToFilename(methodName string) string {
	return e.functionToFilename(methodName)
}

// extractClassBody extracts complete class body
func (e *PythonExtractor) extractClassBody(content, className, inheritance string) string {
	// Find class definition and everything until next class/function at same indent
	classPattern := regexp.MustCompile(
		`(?s)^class\s+` + regexp.QuoteMeta(className) +
		`(?:\([^)]*\))?\s*:\s*\n(.*?)(?=^\S|\z)`)

	match := classPattern.FindStringSubmatch(content)
	if match != nil && len(match) > 1 {
		return "class " + className + "(" + inheritance + "):\n" + match[1]
	}
	
	return "class " + className + "(" + inheritance + "):\n    pass"
}

// extractClassDecorators extracts decorators applied to a class
func (e *PythonExtractor) extractClassDecorators(content, className string) []string {
	// Look for decorators before the class definition
	lines := strings.Split(content, "\n")
	
	var decorators []string
	classFound := false
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		if strings.HasPrefix(trimmed, "class "+className) {
			classFound = true
			break
		}
		
		if strings.HasPrefix(trimmed, "@") {
			decorator := strings.TrimPrefix(trimmed, "@")
			if idx := strings.Index(decorator, "("); idx != -1 {
				decorator = decorator[:idx]
			}
			decorators = append(decorators, decorator)
		}
	}
	
	if !classFound {
		return []string{}
	}
	
	return decorators
}

// extractDataclassBody extracts dataclass body
func (e *PythonExtractor) extractDataclassBody(content, className string) string {
	return e.extractClassBody(content, className, "")
}

// extractPydanticModel extracts Pydantic model
func (e *PythonExtractor) extractPydanticModel(content, className string) string {
	return e.extractClassBody(content, className, "BaseModel")
}

// extractDjangoModel extracts Django model
func (e *PythonExtractor) extractDjangoModel(content, className string) string {
	return e.extractClassBody(content, className, "models.Model")
}

// extractFunctionBody extracts function body
func (e *PythonExtractor) extractFunctionBody(content, funcName, params, returnType string, isAsync bool) string {
	prefix := ""
	if isAsync {
		prefix = "async "
	}
	
	funcPattern := regexp.MustCompile(
		`(?s)^` + prefix + `def\s+` + regexp.QuoteMeta(funcName) + 
		`\([^)]*\)(?:\s*->\s*[^\s:]+)?\s*:\s*\n(.*?)(?=^\S|\z)`)
	
	match := funcPattern.FindStringSubmatch(content)
	if match != nil && len(match) > 1 {
		return prefix + "def " + funcName + "(" + params + ")" + 
			func() string {
				if returnType != "" {
					return " -> " + returnType
				}
				return ""
			}() + ":\n" + match[1]
	}
	
	return prefix + "def " + funcName + "(" + params + ")" + 
		func() string {
			if returnType != "" {
				return " -> " + returnType
			}
			return ""
		}() + ":\n    pass"
}

// extractInitMethodBody extracts __init__ method body
func (e *PythonExtractor) extractInitMethodBody(content, params string) string {
	return e.extractFunctionBody(content, "__init__", params, "", false)
}

// extractClassMethodBody extracts class method body
func (e *PythonExtractor) extractClassMethodBody(content, methodName string) string {
	return e.extractMethodBody(content, methodName)
}

// extractStaticMethodBody extracts static method body
func (e *PythonExtractor) extractStaticMethodBody(content, methodName string) string {
	return e.extractMethodBody(content, methodName)
}

// extractMethodBody extracts method body
func (e *PythonExtractor) extractMethodBody(content, methodName string) string {
	// Find method definition and everything until next method/function at same indent
	methodPattern := regexp.MustCompile(
		`(?s)^\s+def\s+` + regexp.QuoteMeta(methodName) + 
		`\([^)]*\)(?:\s*->\s*[^\s:]+)?\s*:\s*\n(.*?)(?=\n\s+\w|\n\s*$|\n\w|\z)`)

	match := methodPattern.FindStringSubmatch(content)
	if match != nil && len(match) > 1 {
		return "def " + methodName + "():\n" + match[1]
	}
	
	return "def " + methodName + "():\n    pass"
}

// extractPropertyMethod extracts property method
func (e *PythonExtractor) extractPropertyMethod(content, methodName string) string {
	return e.extractMethodBody(content, methodName)
}

// extractDecoratedFunction extracts decorated function
func (e *PythonExtractor) extractDecoratedFunction(content, funcName, decorator string) string {
	return decorator + "\n" + e.extractFunctionBody(content, funcName, "", "", false)
}

// extractTestBody extracts test function body
func (e *PythonExtractor) extractTestBody(content, testName string) string {
	return e.extractFunctionBody(content, testName, "assert", "", false)
}

// extractTestFunctionBody extracts test function body
func (e *PythonExtractor) extractTestFunctionBody(content, funcName string) string {
	return e.extractFunctionBody(content, funcName, "assert", "", false)
}

// extractBenchmarkBody extracts benchmark function body
func (e *PythonExtractor) extractBenchmarkBody(content, benchmarkName string) string {
	return e.extractFunctionBody(content, benchmarkName, "assert", "", false)
}

// extractExampleBody extracts example function body
func (e *PythonExtractor) extractExampleBody(content, exampleName string) string {
	return e.extractFunctionBody(content, exampleName, "assert", "", false)
}

// extractPytestFixtures extracts pytest fixtures
func (e *PythonExtractor) extractPytestFixtures(content string) string {
	// Find @pytest.fixture decorated functions
	fixturePattern := regexp.MustCompile(`(?s)@pytest\.fixture.*?def\s+(\w+)\([^)]*\)\s*:\s*\n(.*?)(?=\n\w|\n\s*$|\z)`)
	
	matches := fixturePattern.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return ""
	}
	
	var result strings.Builder
	for _, match := range matches {
		if len(match) > 2 {
			result.WriteString(fmt.Sprintf("@pytest.fixture\ndef %s():\n%s\n\n", match[1], match[2]))
		}
	}
	
	return strings.TrimSpace(result.String())
}

// extractGinRoute extracts Gin route
func (e *PythonExtractor) extractGinRoute(content, routePath string) string {
	return fmt.Sprintf("Route: %s\nContent: %s", routePath, content)
}

// extractEchoRoute extracts Echo route
func (e *PythonExtractor) extractEchoRoute(content, routePath string) string {
	return fmt.Sprintf("Route: %s\nContent: %s", routePath, content)
}

// extractHttpRoute extracts HTTP route
func (e *PythonExtractor) extractHttpRoute(content, routePath string) string {
	return fmt.Sprintf("Route: %s\nContent: %s", routePath, content)
}

// extractCobraCommand extracts Cobra command
func (e *PythonExtractor) extractCobraCommand(content string) string {
	return content
}

// extractWireInjection extracts Wire injection
func (e *PythonExtractor) extractWireInjection(content string) string {
	return content
}

// extractGrpcService extracts gRPC service
func (e *PythonExtractor) extractGrpcService(content, serviceName string) string {
	return fmt.Sprintf("Service: %s\nContent: %s", serviceName, content)
}

// extractProtobufMessage extracts protobuf message
func (e *PythonExtractor) extractProtobufMessage(content, messageName string) string {
	return fmt.Sprintf("Message: %s\nContent: %s", messageName, content)
}

// extractSetupContent extracts setup.py content
func (e *PythonExtractor) extractSetupContent(content string) string {
	return content
}

// extractTomlSections extracts TOML sections
func (e *PythonExtractor) extractTomlSections(content string) []string {
	var sections []string
	lines := strings.Split(content, "\n")
	
	var currentSection strings.Builder
	inSection := false
	
	for _, line := range lines {
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]") {
			if inSection && currentSection.Len() > 0 {
				sections = append(sections, currentSection.String())
				currentSection.Reset()
			}
			inSection = true
		}
		
		if inSection {
			currentSection.WriteString(line + "\n")
		}
	}
	
	if inSection && currentSection.Len() > 0 {
		sections = append(sections, currentSection.String())
	}
	
	return sections
}

// extractSetupCfgSections extracts setup.cfg sections
func (e *PythonExtractor) extractSetupCfgSections(content string) []string {
	var sections []string
	lines := strings.Split(content, "\n")
	
	var currentSection strings.Builder
	inSection := false
	
	for _, line := range lines {
		if strings.HasPrefix(line, "[") && strings.Contains(line, "]") {
			if inSection && currentSection.Len() > 0 {
				sections = append(sections, currentSection.String())
				currentSection.Reset()
			}
			inSection = true
		}
		
		if inSection {
			currentSection.WriteString(line + "\n")
		}
	}
	
	if inSection && currentSection.Len() > 0 {
		sections = append(sections, currentSection.String())
	}
	
	return sections
}

// extractToxIni extracts tox.ini content
func (e *PythonExtractor) extractToxIni(content string) string {
	return content
}

// extractDjangoView extracts Django view
func (e *PythonExtractor) extractDjangoView(content, viewName string) string {
	return e.extractFunctionBody(content, viewName, "request", "", false)
}

// extractFlaskRoute extracts Flask route
func (e *PythonExtractor) extractFlaskRoute(content, routePath string) string {
	return fmt.Sprintf("Route: %s\nContent: %s", routePath, content)
}

// extractFastAPIRoute extracts FastAPI route
func (e *PythonExtractor) extractFastAPIRoute(content, routePath string) string {
	return fmt.Sprintf("Route: %s\nContent: %s", routePath, content)
}

// reconstructImports reconstructs import statements
func (e *PythonExtractor) reconstructImports(imports []string) string {
	if len(imports) == 0 {
		return ""
	}
	
	var result strings.Builder
	for _, imp := range imports {
		if strings.Contains(imp, " as ") {
			parts := strings.Split(imp, " as ")
			result.WriteString(fmt.Sprintf("import %s as %s\n", parts[0], parts[1]))
		} else {
			result.WriteString(fmt.Sprintf("import %s\n", imp))
		}
	}
	
	return result.String()
}

// getLineWithPattern finds a line containing a specific pattern
func (e *PythonExtractor) getLineWithPattern(content, pattern string) string {
	re := regexp.MustCompile(pattern)
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		if re.MatchString(line) {
			return line
		}
	}
	
	return ""
}

// testToFilename converts test name to filename
func (e *PythonExtractor) testToFilename(testName string) string {
	return e.functionToFilename(testName)
}

// benchmarkToFilename converts benchmark name to filename
func (e *PythonExtractor) benchmarkToFilename(benchmarkName string) string {
	return e.functionToFilename(benchmarkName)
}

// exampleToFilename converts example name to filename
func (e *PythonExtractor) exampleToFilename(exampleName string) string {
	return e.functionToFilename(exampleName)
}

// variableToFilename converts variable name to filename
func (e *PythonExtractor) variableToFilename(varName string) string {
	return varName + ".py"
}

// constantToFilename converts constant name to filename
func (e *PythonExtractor) constantToFilename(constName string) string {
	return constName + ".py"
}