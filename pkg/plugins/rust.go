package main

import (
	"fmt"
	"regexp"
	"strings"
)

// RustExtractor extracts Rust code
type RustExtractor struct {
	patterns map[string]*regexp.Regexp
}

func NewRustExtractor() *RustExtractor {
	return &RustExtractor{}
}

func (e *RustExtractor) Name() string {
	return "rust"
}

func (e *RustExtractor) Extensions() []string {
	return []string{".rs", ".toml"}
}

func (e *RustExtractor) Initialize() error {
	e.patterns = make(map[string]*regexp.Regexp)

	patterns := map[string]string{
		"module":     `mod\s+(\w+)\s*\{`,
		"struct":     `struct\s+(\w+)\s*\{`,
		"enum":       `enum\s+(\w+)\s*\{`,
		"impl":       `impl\s+(\w+)\s*\{`,
		"function":   `fn\s+(\w+)\s*\([^)]*\)`,
		"use":        `use\s+([\w:*]+)`,
		"cargo_toml": `\[package\]`,
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
	return strings.HasSuffix(filename, ".rs") ||
		strings.HasSuffix(filename, "Cargo.toml")
}

func (e *RustExtractor) Extract(content string) []CodeBlock {
	var blocks []CodeBlock

	// Extract use statements
	var imports []string
	for _, match := range e.patterns["use"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			imports = append(imports, match[1])
		}
	}

	// Extract structs
	for _, match := range e.patterns["struct"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			structName := match[1]
			blocks = append(blocks, CodeBlock{
				Content:  e.extractStruct(content, structName),
				Type:     "struct",
				Package:  "",
				Filename: structName + ".rs",
				Language: "rust",
				Imports:  imports,
			})
		}
	}

	// Extract enums
	for _, match := range e.patterns["enum"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			enumName := match[1]
			blocks = append(blocks, CodeBlock{
				Content:  e.extractEnum(content, enumName),
				Type:     "enum",
				Package:  "",
				Filename: enumName + ".rs",
				Language: "rust",
				Imports:  imports,
			})
		}
	}

	// Extract functions
	for _, match := range e.patterns["function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			blocks = append(blocks, CodeBlock{
				Content:  e.extractFunction(content, funcName),
				Type:     "function",
				Package:  "",
				Filename: funcName + ".rs",
				Language: "rust",
				Imports:  imports,
			})
		}
	}

	// Extract Cargo.toml if present
	if strings.Contains(content, "[package]") {
		if cargo := e.extractCargoToml(content); cargo != "" {
			blocks = append(blocks, CodeBlock{
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

func (e *RustExtractor) extractStruct(content, structName string) string {
	pattern := regexp.MustCompile(`struct\s+` + regexp.QuoteMeta(structName) + `\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "struct " + structName + " {\n    // fields\n}"
}

func (e *RustExtractor) extractEnum(content, enumName string) string {
	pattern := regexp.MustCompile(`enum\s+` + regexp.QuoteMeta(enumName) + `\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "enum " + enumName + " {\n    // variants\n}"
}

func (e *RustExtractor) extractFunction(content, funcName string) string {
	pattern := regexp.MustCompile(`fn\s+` + regexp.QuoteMeta(funcName) + `\s*\([^)]*\)\s*\{[^}]*\}`)
	match := pattern.FindString(content)
	if match != "" {
		return match
	}
	return "fn " + funcName + "() {\n    // function body\n}"
}

func (e *RustExtractor) extractCargoToml(content string) string {
	// Extract Cargo.toml section
	start := strings.Index(content, "[package]")
	if start == -1 {
		return "[package]\nname = \"project\"\nversion = \"0.1.0\"\nedition = \"2021\"\n"
	}

	end := strings.Index(content[start:], "\n\n")
	if end == -1 {
		end = len(content)
	} else {
		end = start + end
	}

	return content[start:end]
}

// Export the plugin
var Plugin ExtractorPlugin = NewRustExtractor()
