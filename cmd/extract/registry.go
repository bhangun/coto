package extract

import (
	"strings"

	"github.com/bhangun/coto/pkg/extractor"
)

// PluginRegistry manages all extractor plugins
type PluginRegistry struct {
	plugins        map[string]extractor.ExtractorPlugin // name -> plugin
	extensionMap   map[string]extractor.ExtractorPlugin // extension -> plugin
	languageMap    map[string]extractor.ExtractorPlugin // language -> plugin
}

// NewPluginRegistry creates a new plugin registry
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins:      make(map[string]extractor.ExtractorPlugin),
		extensionMap: make(map[string]extractor.ExtractorPlugin),
		languageMap:  make(map[string]extractor.ExtractorPlugin),
	}
}

// Register adds a new extractor plugin to the registry
func (r *PluginRegistry) Register(plugin extractor.ExtractorPlugin) error {
	name := plugin.Name()
	r.plugins[name] = plugin

	// Register extensions
	for _, ext := range plugin.Extensions() {
		r.extensionMap[strings.ToLower(ext)] = plugin
	}

	// Register by language name (using plugin name as language identifier)
	r.languageMap[strings.ToLower(name)] = plugin

	return nil
}

// GetExtractorByName returns an extractor by its name
func (r *PluginRegistry) GetExtractorByName(name string) extractor.ExtractorPlugin {
	return r.plugins[strings.ToLower(name)]
}

// GetExtractorByExtension returns an extractor by file extension
func (r *PluginRegistry) GetExtractorByExtension(ext string) extractor.ExtractorPlugin {
	return r.extensionMap[strings.ToLower(ext)]
}

// GetExtractorByLanguage returns an extractor by language name
func (r *PluginRegistry) GetExtractorByLanguage(lang string) extractor.ExtractorPlugin {
	return r.languageMap[strings.ToLower(lang)]
}

// GetAllPlugins returns all registered plugins
func (r *PluginRegistry) GetAllPlugins() []extractor.ExtractorPlugin {
	var plugins []extractor.ExtractorPlugin
	for _, plugin := range r.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// registerBuiltInPlugins registers all built-in extractor plugins
func registerBuiltInPlugins(registry *PluginRegistry) {
	// Register Java extractor
	javaExtractor := NewJavaExtractor()
	registry.Register(javaExtractor)

	// Register Go extractor
	goExtractor := NewGoExtractor()
	registry.Register(goExtractor)

	// Register Python extractor
	pythonExtractor := NewPythonExtractor()
	registry.Register(pythonExtractor)

	// Register JavaScript extractor
	jsExtractor := NewJavaScriptExtractor()
	registry.Register(jsExtractor)

	// Register Rust extractor
	rustExtractor := NewRustExtractor()
	registry.Register(rustExtractor)

	// Register Dart extractor
	dartExtractor := NewDartExtractor()
	registry.Register(dartExtractor)

	// Register generic extractor as fallback
	genericExtractor := NewGenericExtractor()
	registry.Register(genericExtractor)
}