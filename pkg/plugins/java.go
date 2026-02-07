package plugins

import (
	"fmt"
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
