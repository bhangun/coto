package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// DartExtractor extracts Dart and Flutter code
type DartExtractor struct {
	patterns map[string]*regexp.Regexp
}

// NewDartExtractor creates a new Dart/Flutter extractor
func NewDartExtractor() *DartExtractor {
	return &DartExtractor{}
}

// Name returns the extractor name
func (e *DartExtractor) Name() string {
	return "dart"
}

// Extensions returns supported file extensions
func (e *DartExtractor) Extensions() []string {
	return []string{".dart", ".yaml", ".pub", ".lock"}
}

// Initialize sets up regex patterns
func (e *DartExtractor) Initialize() error {
	e.patterns = make(map[string]*regexp.Regexp)

	patterns := map[string]string{
		// Dart language patterns
		"import":  `^import\s+['"]([^'"]+)['"]\s*(?:as\s+(\w+))?\s*;`,
		"export":  `^export\s+['"]([^'"]+)['"]\s*;`,
		"library": `^library\s+([\w.]+)\s*;`,
		"part":    `^part\s+['"]([^'"]+)['"]\s*;`,
		"part_of": `^part\s+of\s+([\w.]+)\s*;`,

		// Class and struct patterns
		"class":          `(?s)(?:abstract\s+)?class\s+(\w+)(?:<[^>]+>)?(?:\s+extends\s+\w+)?(?:\s+implements\s+[^{]+)?(?:\s+with\s+[^{]+)?\s*\{`,
		"abstract_class": `abstract\s+class\s+(\w+)(?:<[^>]+>)?\s*\{`,
		"mixin":          `mixin\s+(\w+)(?:<[^>]+>)?(?:\s+on\s+[^{]+)?\s*\{`,
		"enum":           `enum\s+(\w+)\s*\{`,
		"typedef":        `typedef\s+(\w+)\s*=`,

		// Function patterns
		"function":    `(?s)(?:@[\w.]+\s*\n\s*)*(?:Future\s*<[^>]+>\s*)?(?:static\s+)?(?:const\s+)?(?:get\s+|set\s+)?(?:async\s*\*?\s*)?(?:void|\w+)\s+(\w+)\s*\([^)]*\)(?:\s*=>\s*[^;]+)?\s*\{`,
		"factory":     `factory\s+([\w.]+)\(`,
		"constructor": `(\w+)\s*\([^)]*\)\s*(?::\s*[^{]+)?\s*\{`,

		// Widget patterns (Flutter specific)
		"widget_class": `class\s+(\w+)\s+extends\s+(?:Stateless|Stateful)Widget\s*\{`,
		"state_class":  `class\s+(\w+)\s+extends\s+State\s*<\w+>\s*\{`,
		"build_method": `Widget\s+build\s*\([^)]*\)\s*\{`,

		// Flutter specific patterns
		"material_app":     `MaterialApp\s*\(`,
		"scaffold":         `Scaffold\s*\(`,
		"stateless_widget": `extends\s+StatelessWidget`,
		"stateful_widget":  `extends\s+StatefulWidget`,

		// Configuration files
		"pubspec_yaml":     `name:\s*([^\n]+)`,
		"dependencies":     `dependencies:\s*\n(?:[ \t]+[^\n]+\n)*`,
		"dev_dependencies": `dev_dependencies:\s*\n(?:[ \t]+[^\n]+\n)*`,
		"flutter_section":  `flutter:\s*\n(?:[ \t]+[^\n]+\n)*`,

		// JSON patterns for .idea files
		"json_object": `\{[^{}]*"name":\s*"[^"]*"[^{}]*\}`,
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
func (e *DartExtractor) Cleanup() {
	e.patterns = nil
}

// ShouldProcess checks if this extractor should handle the file
func (e *DartExtractor) ShouldProcess(filename string) bool {
	lowerName := strings.ToLower(filename)

	// Check by extension
	ext := filepath.Ext(lowerName)
	if ext == ".dart" || ext == ".yaml" || ext == ".pub" || ext == ".lock" {
		return true
	}

	// Check by filename
	if filename == "pubspec.yaml" || filename == "pubspec.yml" ||
		filename == "analysis_options.yaml" || filename == ".packages" ||
		strings.Contains(lowerName, "flutter") {
		return true
	}

	// Check content-based heuristics (if we had content)
	return false
}

// Extract extracts Dart/Flutter code blocks from content
func (e *DartExtractor) Extract(content string) []CodeBlock {
	var blocks []CodeBlock

	// Extract imports first (used by other blocks)
	imports := e.extractImports(content)

	// Extract Flutter/Dart specific blocks
	blocks = append(blocks, e.extractClasses(content, imports)...)
	blocks = append(blocks, e.extractWidgets(content, imports)...)
	blocks = append(blocks, e.extractFunctions(content, imports)...)
	blocks = append(blocks, e.extractMixins(content, imports)...)
	blocks = append(blocks, e.extractEnums(content, imports)...)

	// Extract configuration files
	blocks = append(blocks, e.extractPubspec(content)...)
	blocks = append(blocks, e.extractAnalysisOptions(content)...)
	blocks = append(blocks, e.extractBuildConfigs(content)...)

	return blocks
}

// extractImports extracts import statements
func (e *DartExtractor) extractImports(content string) []string {
	var imports []string

	for _, match := range e.patterns["import"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			importPath := match[1]
			imports = append(imports, importPath)
		}
	}

	return imports
}

// extractClasses extracts Dart classes
func (e *DartExtractor) extractClasses(content string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	// Extract regular classes
	for _, match := range e.patterns["class"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			className := match[1]
			classContent := e.extractClassBody(content, className)

			// Determine if it's abstract
			isAbstract := strings.Contains(match[0], "abstract")

			blocks = append(blocks, CodeBlock{
				Content:  classContent,
				Type:     "class",
				Package:  e.extractPackageName(content),
				Filename: className + ".dart",
				Language: "dart",
				Imports:  imports,
				Modifiers: func() []string {
					if isAbstract {
						return []string{"abstract"}
					}
					return nil
				}(),
			})
		}
	}

	// Extract abstract classes separately
	for _, match := range e.patterns["abstract_class"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			className := match[1]
			classContent := e.extractClassBody(content, className)

			blocks = append(blocks, CodeBlock{
				Content:   classContent,
				Type:      "abstract_class",
				Package:   e.extractPackageName(content),
				Filename:  className + ".dart",
				Language:  "dart",
				Imports:   imports,
				Modifiers: []string{"abstract"},
			})
		}
	}

	return blocks
}

// extractWidgets extracts Flutter widgets
func (e *DartExtractor) extractWidgets(content string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	// Extract widget classes
	for _, match := range e.patterns["widget_class"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			widgetName := match[1]
			widgetContent := e.extractWidgetBody(content, widgetName)

			blocks = append(blocks, CodeBlock{
				Content:  widgetContent,
				Type:     "widget",
				Package:  e.extractPackageName(content),
				Filename: widgetName + ".dart",
				Language: "dart",
				Imports:  imports,
			})
		}
	}

	// Extract state classes
	for _, match := range e.patterns["state_class"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			stateName := match[1]
			stateContent := e.extractStateBody(content, stateName)

			blocks = append(blocks, CodeBlock{
				Content:  stateContent,
				Type:     "state",
				Package:  e.extractPackageName(content),
				Filename: stateName + ".dart",
				Language: "dart",
				Imports:  imports,
			})
		}
	}

	// Extract MaterialApp widgets
	if e.patterns["material_app"].MatchString(content) {
		materialAppContent := e.extractMaterialApp(content)
		if materialAppContent != "" {
			blocks = append(blocks, CodeBlock{
				Content:  materialAppContent,
				Type:     "material_app",
				Package:  e.extractPackageName(content),
				Filename: "main.dart",
				Language: "dart",
				Imports:  imports,
			})
		}
	}

	// Extract Scaffold widgets
	if e.patterns["scaffold"].MatchString(content) {
		scaffoldContent := e.extractScaffold(content)
		if scaffoldContent != "" {
			blocks = append(blocks, CodeBlock{
				Content:  scaffoldContent,
				Type:     "scaffold",
				Package:  e.extractPackageName(content),
				Filename: "home_screen.dart",
				Language: "dart",
				Imports:  imports,
			})
		}
	}

	return blocks
}

// extractFunctions extracts Dart functions
func (e *DartExtractor) extractFunctions(content string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	for _, match := range e.patterns["function"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			funcName := match[1]
			funcContent := e.extractFunctionBody(content, funcName)

			// Determine function type
			funcType := "function"
			if strings.Contains(match[0], "async") {
				funcType = "async_function"
			}
			if strings.Contains(match[0], "get ") {
				funcType = "getter"
			}
			if strings.Contains(match[0], "set ") {
				funcType = "setter"
			}

			blocks = append(blocks, CodeBlock{
				Content:  funcContent,
				Type:     funcType,
				Package:  e.extractPackageName(content),
				Filename: funcName + ".dart",
				Language: "dart",
				Imports:  imports,
			})
		}
	}

	return blocks
}

// extractMixins extracts Dart mixins
func (e *DartExtractor) extractMixins(content string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	for _, match := range e.patterns["mixin"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			mixinName := match[1]
			mixinContent := e.extractMixinBody(content, mixinName)

			blocks = append(blocks, CodeBlock{
				Content:  mixinContent,
				Type:     "mixin",
				Package:  e.extractPackageName(content),
				Filename: mixinName + ".dart",
				Language: "dart",
				Imports:  imports,
			})
		}
	}

	return blocks
}

// extractEnums extracts Dart enums
func (e *DartExtractor) extractEnums(content string, imports []string) []CodeBlock {
	var blocks []CodeBlock

	for _, match := range e.patterns["enum"].FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			enumName := match[1]
			enumContent := e.extractEnumBody(content, enumName)

			blocks = append(blocks, CodeBlock{
				Content:  enumContent,
				Type:     "enum",
				Package:  e.extractPackageName(content),
				Filename: enumName + ".dart",
				Language: "dart",
				Imports:  imports,
			})
		}
	}

	return blocks
}

// extractPubspec extracts pubspec.yaml configuration
func (e *DartExtractor) extractPubspec(content string) []CodeBlock {
	var blocks []CodeBlock

	// Check if this looks like a pubspec.yaml
	if strings.Contains(content, "name:") &&
		(strings.Contains(content, "dependencies:") ||
			strings.Contains(content, "flutter:")) {

		// Extract project name
		projectName := "app"
		if match := e.patterns["pubspec_yaml"].FindStringSubmatch(content); match != nil {
			if len(match) > 1 {
				projectName = strings.TrimSpace(match[1])
			}
		}

		// Extract dependencies section
		depsContent := e.extractYamlSection(content, "dependencies")
		devDepsContent := e.extractYamlSection(content, "dev_dependencies")
		flutterContent := e.extractYamlSection(content, "flutter")

		// Create comprehensive pubspec
		pubspec := fmt.Sprintf(`name: %s
description: A Flutter application
version: 1.0.0+1

environment:
  sdk: ">=2.18.0 <3.0.0"

dependencies:%s

dev_dependencies:%s

flutter:%s`,
			projectName,
			depsContent,
			devDepsContent,
			flutterContent)

		blocks = append(blocks, CodeBlock{
			Content:  pubspec,
			Type:     "pubspec",
			Package:  "",
			Filename: "pubspec.yaml",
			Language: "yaml",
		})
	}

	return blocks
}

// extractAnalysisOptions extracts analysis_options.yaml
func (e *DartExtractor) extractAnalysisOptions(content string) []CodeBlock {
	var blocks []CodeBlock

	if strings.Contains(content, "analyzer:") || strings.Contains(content, "linter:") {
		// Try to extract the analysis options section
		analyzerContent := e.extractYamlSection(content, "analyzer")
		linterContent := e.extractYamlSection(content, "linter")

		if analyzerContent != "" || linterContent != "" {
			analysisOpts := fmt.Sprintf(`include: package:flutter_lints/flutter.yaml

analyzer:%s

linter:%s`,
				analyzerContent,
				linterContent)

			blocks = append(blocks, CodeBlock{
				Content:  analysisOpts,
				Type:     "analysis_options",
				Package:  "",
				Filename: "analysis_options.yaml",
				Language: "yaml",
			})
		}
	}

	return blocks
}

// extractBuildConfigs extracts build configuration files
func (e *DartExtractor) extractBuildConfigs(content string) []CodeBlock {
	var blocks []CodeBlock

	// Extract build.yaml if present
	if strings.Contains(content, "builders:") || strings.Contains(content, "targets:") {
		buildersContent := e.extractYamlSection(content, "builders")
		targetsContent := e.extractYamlSection(content, "targets")

		if buildersContent != "" || targetsContent != "" {
			buildConfig := fmt.Sprintf(`builders:%s

targets:%s`,
				buildersContent,
				targetsContent)

			blocks = append(blocks, CodeBlock{
				Content:  buildConfig,
				Type:     "build_config",
				Package:  "",
				Filename: "build.yaml",
				Language: "yaml",
			})
		}
	}

	return blocks
}

// Helper Methods

// extractPackageName extracts package name from content
func (e *DartExtractor) extractPackageName(content string) string {
	// Check for library declaration
	if match := e.patterns["library"].FindStringSubmatch(content); match != nil {
		if len(match) > 1 {
			return match[1]
		}
	}

	// Check for part-of declaration
	if match := e.patterns["part_of"].FindStringSubmatch(content); match != nil {
		if len(match) > 1 {
			return match[1]
		}
	}

	// Default package name
	return "main"
}

// extractClassBody extracts the complete class body
func (e *DartExtractor) extractClassBody(content, className string) string {
	// Find class definition and capture until matching closing brace
	classPattern := regexp.MustCompile(
		`(?s)((?:abstract\s+)?class\s+` + regexp.QuoteMeta(className) +
			`(?:<[^>]+>)?(?:\s+extends\s+\w+)?(?:\s+implements\s+[^{]+)?(?:\s+with\s+[^{]+)?\s*\{[^}]*?(?:\{[^}]*\}[^}]*?)*\})`)

	match := classPattern.FindString(content)
	if match != "" {
		return match
	}

	// Fallback: simple class template
	return fmt.Sprintf("class %s {\n  // Class body\n}", className)
}

// extractWidgetBody extracts Flutter widget body
func (e *DartExtractor) extractWidgetBody(content, widgetName string) string {
	// Find widget class with build method
	widgetPattern := regexp.MustCompile(
		`(?s)(class\s+` + regexp.QuoteMeta(widgetName) +
			`\s+extends\s+(?:Stateless|Stateful)Widget\s*\{[^}]*?(?:@override[^}]*?)?(?:Widget\s+build[^}]*\{[^}]*\}[^}]*?)*\})`)

	match := widgetPattern.FindString(content)
	if match != "" {
		return match
	}

	// Fallback: simple widget template
	return fmt.Sprintf(`import 'package:flutter/material.dart';

class %s extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      child: Text('$widgetName'),
    );
  }
}`, widgetName)
}

// extractStateBody extracts State class body
func (e *DartExtractor) extractStateBody(content, stateName string) string {
	statePattern := regexp.MustCompile(
		`(?s)(class\s+` + regexp.QuoteMeta(stateName) +
			`\s+extends\s+State\s*<[^>]+>\s*\{[^}]*?(?:@override[^}]*?)?(?:Widget\s+build[^}]*\{[^}]*\}[^}]*?)*\})`)

	match := statePattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf(`class %s extends State<MyWidget> {
  @override
  Widget build(BuildContext context) {
    return Container();
  }
}`, stateName)
}

// extractMaterialApp extracts MaterialApp widget
func (e *DartExtractor) extractMaterialApp(content string) string {
	// Find MaterialApp widget
	pattern := regexp.MustCompile(`MaterialApp\s*\([^)]*\)`)
	match := pattern.FindString(content)
	if match != "" {
		return fmt.Sprintf(`import 'package:flutter/material.dart';

void main() {
  runApp(
    %s
  );
}`, match)
	}

	return ""
}

// extractScaffold extracts Scaffold widget
func (e *DartExtractor) extractScaffold(content string) string {
	pattern := regexp.MustCompile(`Scaffold\s*\([^)]*\)`)
	match := pattern.FindString(content)
	if match != "" {
		return fmt.Sprintf(`import 'package:flutter/material.dart';

class HomeScreen extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return %s;
  }
}`, match)
	}

	return ""
}

// extractFunctionBody extracts function body
func (e *DartExtractor) extractFunctionBody(content, funcName string) string {
	// Find function with body
	funcPattern := regexp.MustCompile(
		`(?s)((?:@[\w.]+\s*\n\s*)*(?:Future\s*<[^>]+>\s*)?(?:static\s+)?(?:const\s+)?(?:get\s+|set\s+)?(?:async\s*\*?\s*)?(?:void|\w+)\s+` +
			regexp.QuoteMeta(funcName) + `\s*\([^)]*\)(?:\s*=>\s*[^;]+)?\s*\{[^}]*\})`)

	match := funcPattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf("void %s() {\n  // Function body\n}", funcName)
}

// extractMixinBody extracts mixin body
func (e *DartExtractor) extractMixinBody(content, mixinName string) string {
	pattern := regexp.MustCompile(
		`(?s)(mixin\s+` + regexp.QuoteMeta(mixinName) +
			`(?:<[^>]+>)?(?:\s+on\s+[^{]+)?\s*\{[^}]*\})`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf("mixin %s {\n  // Mixin body\n}", mixinName)
}

// extractEnumBody extracts enum body
func (e *DartExtractor) extractEnumBody(content, enumName string) string {
	pattern := regexp.MustCompile(
		`(?s)(enum\s+` + regexp.QuoteMeta(enumName) + `\s*\{[^}]*\})`)

	match := pattern.FindString(content)
	if match != "" {
		return match
	}

	return fmt.Sprintf("enum %s {\n  value1,\n  value2,\n}", enumName)
}

// extractYamlSection extracts a section from YAML content
func (e *DartExtractor) extractYamlSection(content, section string) string {
	sectionPattern := regexp.MustCompile(
		section + `:\s*\n(?:[ \t]+[^\n]+\n)*`)

	match := sectionPattern.FindString(content)
	if match != "" {
		return "\n" + match[len(section)+1:] // Remove "section:"
	}

	return ""
}

// Register the Dart extractor
func init() {
	// This would be called in main.go during initialization
}
