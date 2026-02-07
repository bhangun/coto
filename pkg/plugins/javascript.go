// javascript_extractor.go - JavaScript/TypeScript Code Extractor
package plugins

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bhangun/coto/pkg/extractor"
)

// JavaScriptExtractor extracts JavaScript and TypeScript code
type JavaScriptExtractor struct {
	patterns map[string]*regexp.Regexp
}

// NewJavaScriptExtractor creates a new JavaScript/TypeScript extractor
func NewJavaScriptExtractor() *JavaScriptExtractor {
	return &JavaScriptExtractor{}
}

// Name returns the extractor name
func (e *JavaScriptExtractor) Name() string {
	return "javascript"
}

// Extensions returns supported file extensions
func (e *JavaScriptExtractor) Extensions() []string {
	return []string{
		".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs",
		".vue", ".svelte", // Framework files
	}
}

// Initialize sets up regex patterns
func (e *JavaScriptExtractor) Initialize() error {
	e.patterns = make(map[string]*regexp.Regexp)
	
	patterns := map[string]string{
		// Import/Export patterns
		"import": `import\s+(?:(?:\*\s+as\s+\w+)|(?:\{[^}]+\})|(?:\w+(?:,\s*\w+)*)|(?:\w+\s+as\s+\w+)|)\s+from\s+['"]([^'"]+)['"]`,
		"import_default": `import\s+(\w+)\s+from\s+['"]([^'"]+)['"]`,
		"import_dynamic": `import\s*\(['"]([^'"]+)['"]\)`,
		"export": `export\s+(?:default\s+)?(?:const|let|var|function|class|async\s+function|interface|type|enum)\s+(\w+)`,
		"export_named": `export\s+\{[^}]+\}`,
		"export_default": `export\s+default\s+`,
		
		// CommonJS patterns
		"require": `require\s*\(['"]([^'"]+)['"]\)`,
		"module_export": `module\.exports\s*=\s*`,
		"exports": `exports\.(\w+)\s*=\s*`,
		
		// Class patterns
		"class": `class\s+(\w+)(?:<[^>]+>)?(?:\s+extends\s+\w+)?(?:\s+implements\s+[^{]+)?\s*\{`,
		"abstract_class": `abstract\s+class\s+(\w+)\s*\{`,
		
		// Function patterns
		"function": `function\s+(\w+)\s*\([^)]*\)\s*\{`,
		"async_function": `async\s+function\s+(\w+)\s*\([^)]*\)\s*\{`,
		"arrow_function": `(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s+)?\([^)]*\)\s*=>`,
		"arrow_function_const": `const\s+(\w+)\s*=\s*(?:async\s+)?\([^)]*\)\s*=>`,
		"arrow_function_let": `let\s+(\w+)\s*=\s*(?:async\s+)?\([^)]*\)\s*=>`,
		
		// Method patterns
		"method": `(\w+)\s*\([^)]*\)\s*\{`,
		"async_method": `async\s+(\w+)\s*\([^)]*\)\s*\{`,
		"getter": `get\s+(\w+)\s*\(\)\s*\{`,
		"setter": `set\s+(\w+)\s*\([^)]*\)\s*\{`,
		
		// React patterns
		"react_component": `(?:export\s+)?(?:default\s+)?(?:function\s+)?(\w+)\s*\([^)]*\)\s*\{[^}]*?\breturn\s*\(`,
		"react_class_component": `class\s+(\w+)\s+extends\s+(?:React\.)?(?:Component|PureComponent)`,
		"react_hook": `const\s+(\w+)\s*=\s*(?:use[A-Z][\w]*)\s*\(`,
		"jsx_element": `<\w+[^>]*>`,
		"jsx_self_closing": `<\w+[^>]*/>`,
		
		// Vue patterns
		"vue_component": `export\s+default\s*\{[^}]*\bname:\s*['"](\w+)['"]`,
		"vue_setup": `setup\s*\([^)]*\)\s*\{`,
		"vue_composition": `defineComponent\s*\(`,
		"vue_template": `<template>[^<]*</template>`,
		"vue_script": `<script[^>]*>[^<]*</script>`,
		"vue_style": `<style[^>]*>[^<]*</style>`,
		
		// TypeScript patterns
		"interface": `interface\s+(\w+)(?:<[^>]+>)?\s*\{`,
		"type_alias": `type\s+(\w+)(?:<[^>]+>)?\s*=`,
		"enum": `enum\s+(\w+)\s*\{`,
		"namespace": `namespace\s+(\w+)\s*\{`,
		"module": `module\s+['"]([^'"]+)['"]\s*\{`,
		"decorator": `@(\w+)(?:\([^)]*\))?`,
		
		// Variable declarations
		"const_declaration": `const\s+(\w+)\s*=`,
		"let_declaration": `let\s+(\w+)\s*=`,
		"var_declaration": `var\s+(\w+)\s*=`,
		
		// Configuration files
		"package_json": `\{\s*"name"\s*:\s*"[^"]+"`,
		"tsconfig": `\{\s*"compilerOptions"\s*:\s*\{`,
		"babel_config": `(?:module\.exports\s*=|export\s+default\s*)\s*\{`,
		"webpack_config": `module\.exports\s*=\s*\{[^}]*entry:`,
		"eslint_config": `module\.exports\s*=\s*\{[^}]*rules:`,
		
		// Test patterns
		"jest_test": `(?:test|it|describe)\(['"]([^'"]+)['"]`,
		"mocha_test": `it\(['"]([^'"]+)['"]`,
		"vitest_test": `(?:test|it)\(['"]([^'"]+)['"]`,
		
		// Node.js patterns
		"express_route": `(?:app|router)\.(?:get|post|put|delete|patch)\(['"]([^'"]+)['"]`,
		"middleware": `(?:app|router)\.use\(`,
		"controller": `class\s+(\w+)Controller\s*\{`,
		
		// CSS-in-JS patterns
		"styled_component": `const\s+(\w+)\s*=\s*styled\.[\w.]+`,
		"css_object": `const\s+styles\s*=\s*\{`,
		
		// Comment patterns
		"jsdoc": `/\*\*\s*\n(?:\s*\*.*\n)*\s*\*/`,
		"single_line_comment": `//[^\n]*`,
		"multi_line_comment": `/\*[^*]*\*+(?:[^/*][^*]*\*+)*/`,
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
func (e *JavaScriptExtractor) Cleanup() {
	e.patterns = nil
}

// ShouldProcess checks if this extractor should handle the file
func (e *JavaScriptExtractor) ShouldProcess(filename string) bool {
	lowerName := strings.ToLower(filename)
	
	// Check by extension
	ext := filepath.Ext(lowerName)
	switch ext {
	case ".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs", ".vue", ".svelte":
		return true
	}
	
	// Check by filename
	configFiles := []string{
		"package.json", "package-lock.json", "yarn.lock",
		"tsconfig.json", "tsconfig.*.json",
		"babel.config.js", "babel.config.json", ".babelrc",
		"webpack.config.js", "webpack.*.js",
		"eslint.config.js", ".eslintrc.js", ".eslintrc.json",
		"jest.config.js", "jest.config.json",
		"vite.config.js", "vite.config.ts",
		"next.config.js", "nuxt.config.js",
		"rollup.config.js", "vite.config.js",
		"postcss.config.js", "tailwind.config.js",
	}
	
	for _, configFile := range configFiles {
		if strings.Contains(lowerName, configFile) {
			return true
		}
	}
	
	// Check for React/Vue/Svelte framework indicators
	if strings.Contains(lowerName, "react") || strings.Contains(lowerName, "vue") || 
	   strings.Contains(lowerName, "svelte") || strings.Contains(lowerName, "angular") {
		return true
	}
	
	return false
}

// Extract extracts JavaScript/TypeScript code blocks from content
func (e *JavaScriptExtractor) Extract(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	// Extract imports first
	imports := e.extractImports(content)
	
	// Extract JavaScript/TypeScript specific blocks
	blocks = append(blocks, e.extractClasses(content, imports)...)
	blocks = append(blocks, e.extractFunctions(content, imports)...)
	blocks = append(blocks, e.extractVariables(content, imports)...)
	blocks = append(blocks, e.extractInterfaces(content, imports)...)
	blocks = append(blocks, e.extractTypes(content, imports)...)
	blocks = append(blocks, e.extractEnums(content, imports)...)
	
	// Extract framework-specific code
	blocks = append(blocks, e.extractReactComponents(content, imports)...)
	blocks = append(blocks, e.extractVueComponents(content, imports)...)
	blocks = append(blocks, e.extractNodeJSComponents(content, imports)...)
	
	// Extract configuration files
	blocks = append(blocks, e.extractPackageJson(content)...)
	blocks = append(blocks, e.extractConfigFiles(content)...)
	
	// Extract test files
	blocks = append(blocks, e.extractTests(content, imports)...)
	
	return blocks
}

// extractImports extracts import/require statements
func (e *JavaScriptExtractor) extractImports(content string) []string {
	var imports []string
	
	// Extract ES6 imports
	for _, match := range e.patterns["import"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 && match[1] != "" {
			imports = append(imports, match[1])
		}
	}
	
	// Extract default imports
	for _, match := range e.patterns["import_default"].FindAllStringSubmatch(content, -1) {
		if len(match) > 2 {
			imports = append(imports, match[2])
		}
	}
	
	// Extract dynamic imports
	for _, match := range e.patterns["import_dynamic"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			imports = append(imports, match[1])
		}
	}
	
	// Extract CommonJS requires
	for _, match := range e.patterns["require"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			imports = append(imports, match[1])
		}
	}
	
	return imports
}

// extractClasses extracts JavaScript/TypeScript classes
func (e *JavaScriptExtractor) extractClasses(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	// Extract regular classes
	for _, match := range e.patterns["class"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			className := match[1]
			classContent := e.extractClassBody(content, className)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  classContent,
				Type:     "class",
				Package:  e.extractPackageName(content),
				Filename: className + ".js",
				Language: "javascript",
				Imports:  imports,
			})
		}
	}
	
	// Extract abstract classes (TypeScript)
	for _, match := range e.patterns["abstract_class"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			className := match[1]
			classContent := e.extractAbstractClassBody(content, className)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:   classContent,
				Type:      "abstract_class",
				Package:   e.extractPackageName(content),
				Filename:  className + ".ts",
				Language:  "typescript",
				Imports:   imports,
				Modifiers: []string{"abstract"},
			})
		}
	}
	
	// Extract React class components
	for _, match := range e.patterns["react_class_component"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			componentName := match[1]
			componentContent := e.extractReactClassComponent(content, componentName)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  componentContent,
				Type:     "react_class_component",
				Package:  e.extractPackageName(content),
				Filename: componentName + ".jsx",
				Language: "javascript",
				Imports:  append(imports, "react"),
			})
		}
	}
	
	return blocks
}

// extractFunctions extracts functions
func (e *JavaScriptExtractor) extractFunctions(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	// Extract regular functions
	for _, match := range e.patterns["function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			funcContent := e.extractFunctionBody(content, funcName)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  funcContent,
				Type:     "function",
				Package:  e.extractPackageName(content),
				Filename: funcName + ".js",
				Language: "javascript",
				Imports:  imports,
			})
		}
	}
	
	// Extract async functions
	for _, match := range e.patterns["async_function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			funcContent := e.extractAsyncFunctionBody(content, funcName)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:   funcContent,
				Type:      "async_function",
				Package:   e.extractPackageName(content),
				Filename:  funcName + ".js",
				Language:  "javascript",
				Imports:   imports,
				Modifiers: []string{"async"},
			})
		}
	}
	
	// Extract arrow functions
	for _, match := range e.patterns["arrow_function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			funcContent := e.extractArrowFunctionBody(content, funcName)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  funcContent,
				Type:     "arrow_function",
				Package:  e.extractPackageName(content),
				Filename: funcName + ".js",
				Language: "javascript",
				Imports:  imports,
			})
		}
	}
	
	return blocks
}

// extractVariables extracts variable declarations
func (e *JavaScriptExtractor) extractVariables(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	// Extract const declarations
	for _, match := range e.patterns["const_declaration"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			varName := match[1]
			varContent := e.extractVariableBody(content, varName, "const")
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  varContent,
				Type:     "const_declaration",
				Package:  e.extractPackageName(content),
				Filename: varName + ".js",
				Language: "javascript",
				Imports:  imports,
			})
		}
	}
	
	// Extract let declarations
	for _, match := range e.patterns["let_declaration"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			varName := match[1]
			varContent := e.extractVariableBody(content, varName, "let")
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  varContent,
				Type:     "let_declaration",
				Package:  e.extractPackageName(content),
				Filename: varName + ".js",
				Language: "javascript",
				Imports:  imports,
			})
		}
	}
	
	// Extract styled components
	for _, match := range e.patterns["styled_component"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			componentName := match[1]
			styledContent := e.extractStyledComponent(content, componentName)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  styledContent,
				Type:     "styled_component",
				Package:  e.extractPackageName(content),
				Filename: componentName + ".js",
				Language: "javascript",
				Imports:  append(imports, "styled-components"),
			})
		}
	}
	
	return blocks
}

// extractInterfaces extracts TypeScript interfaces
func (e *JavaScriptExtractor) extractInterfaces(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	for _, match := range e.patterns["interface"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			interfaceName := match[1]
			interfaceContent := e.extractInterfaceBody(content, interfaceName)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  interfaceContent,
				Type:     "interface",
				Package:  e.extractPackageName(content),
				Filename: interfaceName + ".ts",
				Language: "typescript",
				Imports:  imports,
			})
		}
	}
	
	return blocks
}

// extractTypes extracts TypeScript type aliases
func (e *JavaScriptExtractor) extractTypes(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	for _, match := range e.patterns["type_alias"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			typeName := match[1]
			typeContent := e.extractTypeBody(content, typeName)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  typeContent,
				Type:     "type_alias",
				Package:  e.extractPackageName(content),
				Filename: typeName + ".ts",
				Language: "typescript",
				Imports:  imports,
			})
		}
	}
	
	return blocks
}

// extractEnums extracts TypeScript enums
func (e *JavaScriptExtractor) extractEnums(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	for _, match := range e.patterns["enum"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			enumName := match[1]
			enumContent := e.extractEnumBody(content, enumName)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  enumContent,
				Type:     "enum",
				Package:  e.extractPackageName(content),
				Filename: enumName + ".ts",
				Language: "typescript",
				Imports:  imports,
			})
		}
	}
	
	return blocks
}

// extractReactComponents extracts React components and hooks
func (e *JavaScriptExtractor) extractReactComponents(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	// Extract React function components
	for _, match := range e.patterns["react_component"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			componentName := match[1]
			componentContent := e.extractReactComponentBody(content, componentName)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  componentContent,
				Type:     "react_component",
				Package:  e.extractPackageName(content),
				Filename: componentName + ".jsx",
				Language: "javascript",
				Imports:  append(imports, "react"),
			})
		}
	}
	
	// Extract React hooks
	for _, match := range e.patterns["react_hook"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			hookName := match[1]
			hookContent := e.extractReactHookBody(content, hookName)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  hookContent,
				Type:     "react_hook",
				Package:  e.extractPackageName(content),
				Filename: "use" + strings.Title(hookName) + ".js",
				Language: "javascript",
				Imports:  append(imports, "react"),
			})
		}
	}
	
	// Extract JSX elements
	if e.patterns["jsx_element"].MatchString(content) || e.patterns["jsx_self_closing"].MatchString(content) {
		jsxContent := e.extractJSXContent(content)
		if jsxContent != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  jsxContent,
				Type:     "jsx",
				Package:  e.extractPackageName(content),
				Filename: "component.jsx",
				Language: "javascript",
				Imports:  append(imports, "react"),
			})
		}
	}
	
	return blocks
}

// extractVueComponents extracts Vue components
func (e *JavaScriptExtractor) extractVueComponents(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	// Extract Vue single-file components
	if e.patterns["vue_component"].MatchString(content) {
		for _, match := range e.patterns["vue_component"].FindAllStringSubmatch(content, -1) {
			if len(match) > 1 {
				componentName := match[1]
				vueContent := e.extractVueComponentBody(content, componentName)
				
				blocks = append(blocks, extractor.CodeBlock{
					Content:  vueContent,
					Type:     "vue_component",
					Package:  e.extractPackageName(content),
					Filename: componentName + ".vue",
					Language: "vue",
					Imports:  imports,
				})
			}
		}
	}
	
	// Extract Vue Composition API
	if e.patterns["vue_setup"].MatchString(content) || e.patterns["vue_composition"].MatchString(content) {
		setupContent := e.extractVueSetupContent(content)
		if setupContent != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  setupContent,
				Type:     "vue_setup",
				Package:  e.extractPackageName(content),
				Filename: "composition.js",
				Language: "javascript",
				Imports:  append(imports, "vue"),
			})
		}
	}
	
	// Extract Vue template
	if e.patterns["vue_template"].MatchString(content) {
		templateContent := e.extractVueTemplate(content)
		if templateContent != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  templateContent,
				Type:     "vue_template",
				Package:  e.extractPackageName(content),
				Filename: "template.html",
				Language: "html",
			})
		}
	}
	
	return blocks
}

// extractNodeJSComponents extracts Node.js specific code
func (e *JavaScriptExtractor) extractNodeJSComponents(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	// Extract Express routes
	for _, match := range e.patterns["express_route"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			routePath := match[1]
			routeContent := e.extractExpressRoute(content, routePath)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  routeContent,
				Type:     "express_route",
				Package:  e.extractPackageName(content),
				Filename: "routes.js",
				Language: "javascript",
				Imports:  append(imports, "express"),
			})
		}
	}
	
	// Extract middleware
	if e.patterns["middleware"].MatchString(content) {
		middlewareContent := e.extractMiddleware(content)
		if middlewareContent != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  middlewareContent,
				Type:     "middleware",
				Package:  e.extractPackageName(content),
				Filename: "middleware.js",
				Language: "javascript",
				Imports:  imports,
			})
		}
	}
	
	// Extract controllers
	for _, match := range e.patterns["controller"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			controllerName := match[1]
			controllerContent := e.extractControllerBody(content, controllerName)
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  controllerContent,
				Type:     "controller",
				Package:  e.extractPackageName(content),
				Filename: controllerName + "Controller.js",
				Language: "javascript",
				Imports:  imports,
			})
		}
	}
	
	return blocks
}

// extractPackageJson extracts package.json
func (e *JavaScriptExtractor) extractPackageJson(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	if e.patterns["package_json"].MatchString(content) {
		packageJson := e.extractJSONObject(content, `^\s*\{[^}]*\}`)
		if packageJson != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  packageJson,
				Type:     "package_json",
				Package:  "",
				Filename: "package.json",
				Language: "json",
			})
		}
	}
	
	return blocks
}

// extractConfigFiles extracts configuration files
func (e *JavaScriptExtractor) extractConfigFiles(content string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	// Extract TypeScript config
	if e.patterns["tsconfig"].MatchString(content) {
		tsconfig := e.extractJSONObject(content, `\{\s*"compilerOptions"\s*:\s*\{[^}]*\}`)
		if tsconfig != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  tsconfig,
				Type:     "tsconfig",
				Package:  "",
				Filename: "tsconfig.json",
				Language: "json",
			})
		}
	}
	
	// Extract Babel config
	if e.patterns["babel_config"].MatchString(content) {
		babelConfig := e.extractModuleExports(content)
		if babelConfig != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  babelConfig,
				Type:     "babel_config",
				Package:  "",
				Filename: "babel.config.js",
				Language: "javascript",
			})
		}
	}
	
	// Extract Webpack config
	if e.patterns["webpack_config"].MatchString(content) {
		webpackConfig := e.extractModuleExports(content)
		if webpackConfig != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  webpackConfig,
				Type:     "webpack_config",
				Package:  "",
				Filename: "webpack.config.js",
				Language: "javascript",
			})
		}
	}
	
	// Extract ESLint config
	if e.patterns["eslint_config"].MatchString(content) {
		eslintConfig := e.extractModuleExports(content)
		if eslintConfig != "" {
			blocks = append(blocks, extractor.CodeBlock{
				Content:  eslintConfig,
				Type:     "eslint_config",
				Package:  "",
				Filename: "eslint.config.js",
				Language: "javascript",
			})
		}
	}
	
	return blocks
}

// extractTests extracts test files
func (e *JavaScriptExtractor) extractTests(content string, imports []string) []extractor.CodeBlock {
	var blocks []extractor.CodeBlock
	
	// Extract Jest tests
	for _, match := range e.patterns["jest_test"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			testName := match[1]
			testContent := e.extractTestBody(content, testName, "jest")
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  testContent,
				Type:     "jest_test",
				Package:  e.extractPackageName(content),
				Filename: testName + ".test.js",
				Language: "javascript",
				Imports:  append(imports, "@testing-library/react", "jest"),
			})
		}
	}
	
	// Extract Mocha tests
	for _, match := range e.patterns["mocha_test"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			testName := match[1]
			testContent := e.extractTestBody(content, testName, "mocha")
			
			blocks = append(blocks, extractor.CodeBlock{
				Content:  testContent,
				Type:     "mocha_test",
				Package:  e.extractPackageName(content),
				Filename: testName + ".spec.js",
				Language: "javascript",
				Imports:  append(imports, "chai", "mocha"),
			})
		}
	}
	
	return blocks
}

// Helper Methods

// extractPackageName extracts package/module name from content
func (e *JavaScriptExtractor) extractPackageName(content string) string {
	// Look for package.json name
	if match := e.patterns["package_json"].FindStringSubmatch(content); match != nil {
		// Try to extract name from JSON
		namePattern := regexp.MustCompile(`"name"\s*:\s*"([^"]+)"`)
		if nameMatch := namePattern.FindStringSubmatch(content); nameMatch != nil && len(nameMatch) > 1 {
			return nameMatch[1]
		}
	}
	
	// Look for export default name
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "export default") {
			// Extract name after export default
			parts := strings.Split(line, "export default")
			if len(parts) > 1 {
				name := strings.TrimSpace(parts[1])
				// Remove semicolon and whitespace
				name = strings.TrimSuffix(name, ";")
				name = strings.TrimSpace(name)
				return name
			}
		}
	}
	
	// Default package name
	return "app"
}

// extractClassBody extracts complete class body
func (e *JavaScriptExtractor) extractClassBody(content, className string) string {
	// Find class definition and everything until next class/function at same level
	classPattern := regexp.MustCompile(
		`(?s)class\s+` + regexp.QuoteMeta(className) + 
		`(?:<[^>]+>)?(?:\s+extends\s+\w+)?(?:\s+implements\s+[^{]+)?\s*\{[^}]*?(?:\{[^}]*\}[^}]*?)*\}`)
	
	match := classPattern.FindString(content)
	if match != "" {
		return match
	}
	
	// Fallback: simple class template
	return fmt.Sprintf("class %s {\n  constructor() {\n    // constructor\n  }\n}", className)
}

// extractAbstractClassBody extracts abstract class body
func (e *JavaScriptExtractor) extractAbstractClassBody(content, className string) string {
	pattern := regexp.MustCompile(
		`(?s)abstract\s+class\s+` + regexp.QuoteMeta(className) + 
		`\s*\{[^}]*?(?:\{[^}]*\}[^}]*?)*\}`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf("abstract class %s {\n  // abstract methods\n}", className)
}

// extractFunctionBody extracts function body
func (e *JavaScriptExtractor) extractFunctionBody(content, funcName string) string {
	pattern := regexp.MustCompile(
		`(?s)function\s+` + regexp.QuoteMeta(funcName) + 
		`\s*\([^)]*\)\s*\{[^}]*\}`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf("function %s() {\n  // function body\n}", funcName)
}

// extractAsyncFunctionBody extracts async function body
func (e *JavaScriptExtractor) extractAsyncFunctionBody(content, funcName string) string {
	pattern := regexp.MustCompile(
		`(?s)async\s+function\s+` + regexp.QuoteMeta(funcName) + 
		`\s*\([^)]*\)\s*\{[^}]*\}`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf("async function %s() {\n  // async function body\n}", funcName)
}

// extractArrowFunctionBody extracts arrow function body
func (e *JavaScriptExtractor) extractArrowFunctionBody(content, funcName string) string {
	// Look for const/let/var declaration with arrow function
	pattern := regexp.MustCompile(
		`(?:const|let|var)\s+` + regexp.QuoteMeta(funcName) + 
		`\s*=\s*(?:async\s+)?\([^)]*\)\s*=>\s*\{[^}]*\}`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf("const %s = () => {\n  // arrow function body\n}", funcName)
}

// extractVariableBody extracts variable declaration body
func (e *JavaScriptExtractor) extractVariableBody(content, varName, varType string) string {
	pattern := regexp.MustCompile(
		regexp.QuoteMeta(varType) + `\s+` + regexp.QuoteMeta(varName) + 
		`\s*=.*?;(?:\n|$)`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf("%s %s = null;", varType, varName)
}

// extractInterfaceBody extracts TypeScript interface body
func (e *JavaScriptExtractor) extractInterfaceBody(content, interfaceName string) string {
	pattern := regexp.MustCompile(
		`(?s)interface\s+` + regexp.QuoteMeta(interfaceName) + 
		`(?:<[^>]+>)?\s*\{[^}]*\}`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf("interface %s {\n  // properties\n}", interfaceName)
}

// extractTypeBody extracts TypeScript type alias body
func (e *JavaScriptExtractor) extractTypeBody(content, typeName string) string {
	pattern := regexp.MustCompile(
		`type\s+` + regexp.QuoteMeta(typeName) + 
		`(?:<[^>]+>)?\s*=.*?;(?:\n|$)`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf("type %s = any;", typeName)
}

// extractEnumBody extracts TypeScript enum body
func (e *JavaScriptExtractor) extractEnumBody(content, enumName string) string {
	pattern := regexp.MustCompile(
		`(?s)enum\s+` + regexp.QuoteMeta(enumName) + `\s*\{[^}]*\}`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf("enum %s {\n  VALUE1,\n  VALUE2,\n}", enumName)
}

// extractReactComponentBody extracts React component body
func (e *JavaScriptExtractor) extractReactComponentBody(content, componentName string) string {
	// Find component function with JSX return
	pattern := regexp.MustCompile(
		`(?s)(?:export\s+)?(?:default\s+)?(?:function\s+)?` + 
		regexp.QuoteMeta(componentName) + 
		`\s*\([^)]*\)\s*\{[^}]*?\breturn\s*\([^)]*\)[^}]*\}`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf(`function %s() {
  return (
    <div>
      <h1>%s</h1>
    </div>
  );
}`, componentName, componentName)
}

// extractReactClassComponent extracts React class component
func (e *JavaScriptExtractor) extractReactClassComponent(content, componentName string) string {
	pattern := regexp.MustCompile(
		`(?s)class\s+` + regexp.QuoteMeta(componentName) + 
		`\s+extends\s+(?:React\.)?(?:Component|PureComponent)[^}]*\{[^}]*\}`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf(`class %s extends React.Component {
  render() {
    return (
      <div>
        <h1>%s</h1>
      </div>
    );
  }
}`, componentName, componentName)
}

// extractReactHookBody extracts React hook body
func (e *JavaScriptExtractor) extractReactHookBody(content, hookName string) string {
	pattern := regexp.MustCompile(
		`const\s+` + regexp.QuoteMeta(hookName) + 
		`\s*=\s*(?:use[A-Z][\w]*)\s*\([^)]*\)\s*\{[^}]*\}`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf("const %s = () => {\n  // hook logic\n  return null;\n};", hookName)
}

// extractJSXContent extracts JSX content
func (e *JavaScriptExtractor) extractJSXContent(content string) string {
	// Find JSX block
	pattern := regexp.MustCompile(`(?s)<[^>]*>[^<]*</[^>]*>`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return "<div>\n  {/* JSX content */}\n</div>"
}

// extractVueComponentBody extracts Vue component body
func (e *JavaScriptExtractor) extractVueComponentBody(content, componentName string) string {
	// Try to extract Vue SFC structure
	template := e.extractVueTemplate(content)
	script := e.extractVueScript(content)
	style := e.extractVueStyle(content)
	
	vueContent := ""
	if template != "" {
		vueContent += template + "\n\n"
	}
	if script != "" {
		vueContent += script + "\n\n"
	}
	if style != "" {
		vueContent += style + "\n"
	}
	
	if vueContent != "" {
		return vueContent
	}
	
	// Fallback Vue component
	return fmt.Sprintf(`<template>
  <div>
    <h1>%s</h1>
  </div>
</template>

<script>
export default {
  name: '%s',
  data() {
    return {
      message: 'Hello Vue!'
    };
  }
};
</script>

<style scoped>
h1 {
  color: #42b983;
}
</style>`, componentName, componentName)
}

// extractVueSetupContent extracts Vue 3 setup function
func (e *JavaScriptExtractor) extractVueSetupContent(content string) string {
	if e.patterns["vue_setup"].MatchString(content) {
		pattern := regexp.MustCompile(`(?s)setup\s*\([^)]*\)\s*\{[^}]*\}`)
		match := pattern.FindString(content)
		if match != "" {
			return match
		}
	}
	
	return `setup() {
  const count = ref(0);
  return { count };
}`
}

// extractVueTemplate extracts Vue template
func (e *JavaScriptExtractor) extractVueTemplate(content string) string {
	pattern := regexp.MustCompile(`(?s)<template[^>]*>([^<]*)</template>`)
	match := pattern.FindStringSubmatch(content)
	if match != nil && len(match) > 1 {
		return "<template>" + match[1] + "</template>"
	}
	return ""
}

// extractVueScript extracts Vue script
func (e *JavaScriptExtractor) extractVueScript(content string) string {
	pattern := regexp.MustCompile(`(?s)<script[^>]*>([^<]*)</script>`)
	match := pattern.FindStringSubmatch(content)
	if match != nil && len(match) > 1 {
		return "<script>" + match[1] + "</script>"
	}
	return ""
}

// extractVueStyle extracts Vue style
func (e *JavaScriptExtractor) extractVueStyle(content string) string {
	pattern := regexp.MustCompile(`(?s)<style[^>]*>([^<]*)</style>`)
	match := pattern.FindStringSubmatch(content)
	if match != nil && len(match) > 1 {
		return "<style>" + match[1] + "</style>"
	}
	return ""
}

// extractExpressRoute extracts Express route handler
func (e *JavaScriptExtractor) extractExpressRoute(content, routePath string) string {
	pattern := regexp.MustCompile(
		`(?:app|router)\.(?:get|post|put|delete|patch)\(['"]` + 
		regexp.QuoteMeta(routePath) + `['"][^)]*\)[^;]*\{[^}]*\}`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf(`app.get('%s', (req, res) => {
  res.send('Response from %s');
});`, routePath, routePath)
}

// extractMiddleware extracts middleware function
func (e *JavaScriptExtractor) extractMiddleware(content string) string {
	pattern := regexp.MustCompile(`(?s)(?:app|router)\.use\([^)]*\)[^;]*\{[^}]*\}`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return `app.use((req, res, next) => {
  console.log('Middleware executed');
  next();
});`
}

// extractControllerBody extracts controller class
func (e *JavaScriptExtractor) extractControllerBody(content, controllerName string) string {
	pattern := regexp.MustCompile(
		`(?s)class\s+` + regexp.QuoteMeta(controllerName) + 
		`Controller\s*\{[^}]*\}`)
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf(`class %sController {
  constructor() {
    // controller initialization
  }
  
  index(req, res) {
    res.json({ message: 'Index method' });
  }
}`, controllerName)
}

// extractStyledComponent extracts styled component
func (e *JavaScriptExtractor) extractStyledComponent(content, componentName string) string {
	pattern := regexp.MustCompile(
		`const\s+` + regexp.QuoteMeta(componentName) + 
		`\s*=\s*styled\.[\w.]+\s*` + "`[^`]*`")
	
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf("const %s = styled.div`\n  color: #333;\n  padding: 1rem;\n`;", componentName)
}

// extractJSONObject extracts JSON object from content
func (e *JavaScriptExtractor) extractJSONObject(content, patternStr string) string {
	pattern := regexp.MustCompile(patternStr)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return ""
}

// extractModuleExports extracts module.exports or export default
func (e *JavaScriptExtractor) extractModuleExports(content string) string {
	// Look for module.exports assignment
	pattern := regexp.MustCompile(`(?s)module\.exports\s*=\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	// Look for export default
	pattern = regexp.MustCompile(`(?s)export\s+default\s*\{[^}]*\}`)
	match = pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return ""
}

// extractTestBody extracts test function body
func (e *JavaScriptExtractor) extractTestBody(content, testName, testFramework string) string {
	patternStr := ""
	switch testFramework {
	case "jest":
		patternStr = `(?:test|it|describe)\(['"]` + regexp.QuoteMeta(testName) + `['"][^)]*\)[^;]*\{[^}]*\}`
	case "mocha":
		patternStr = `it\(['"]` + regexp.QuoteMeta(testName) + `['"][^)]*\)[^;]*\{[^}]*\}`
	default:
		patternStr = `(?:test|it)\(['"]` + regexp.QuoteMeta(testName) + `['"][^)]*\)[^;]*\{[^}]*\}`
	}
	
	pattern := regexp.MustCompile(patternStr)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	
	return fmt.Sprintf(`test('%s', () => {
  expect(true).toBe(true);
});`, testName)
}