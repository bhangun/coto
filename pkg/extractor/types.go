package extractor

// Version information
const Version = "1.0.0"

// CodeBlock represents an extracted code block
type CodeBlock struct {
	Content     string   `json:"content"`
	Type        string   `json:"type"`
	Package     string   `json:"package,omitempty"`
	Filename    string   `json:"filename"`
	Language    string   `json:"language"`
	Imports     []string `json:"imports,omitempty"`
	Annotations []string `json:"annotations,omitempty"`
	Modifiers   []string `json:"modifiers,omitempty"`
}

// ExtractorPlugin defines the interface for language-specific extractors
type ExtractorPlugin interface {
	Name() string
	Extensions() []string
	Extract(content string) []CodeBlock
	ShouldProcess(filename string) bool
	Initialize() error
	Cleanup()
}
