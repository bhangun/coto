package main

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sort"
	"strings"
	"sync"

	"github.com/bhangun/coto/pkg/extractor"
)

// Re-export types for compatibility
type ExtractorPlugin = extractor.ExtractorPlugin
type CodeBlock = extractor.CodeBlock

// PluginRegistry manages all registered extractor plugins
type PluginRegistry struct {
	plugins     map[string]ExtractorPlugin
	extToPlugin map[string]string
	initialized bool
	initMutex   sync.RWMutex
}

// NewPluginRegistry creates a new plugin registry
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins:     make(map[string]ExtractorPlugin),
		extToPlugin: make(map[string]string),
	}
}

// Register adds a plugin to the registry
func (r *PluginRegistry) Register(plugin ExtractorPlugin) error {
	r.initMutex.Lock()
	defer r.initMutex.Unlock()

	name := plugin.Name()
	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin '%s' already registered", name)
	}

	// Initialize plugin
	if err := plugin.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize plugin '%s': %w", name, err)
	}

	r.plugins[name] = plugin

	// Map extensions to plugin
	for _, ext := range plugin.Extensions() {
		r.extToPlugin[strings.ToLower(ext)] = name
	}

	r.initialized = true
	return nil
}

// GetPlugin returns a plugin by name
func (r *PluginRegistry) GetPlugin(name string) (ExtractorPlugin, bool) {
	r.initMutex.RLock()
	defer r.initMutex.RUnlock()

	plugin, exists := r.plugins[name]
	return plugin, exists
}

// GetPluginByExtension returns a plugin by file extension
func (r *PluginRegistry) GetPluginByExtension(ext string) (ExtractorPlugin, bool) {
	r.initMutex.RLock()
	defer r.initMutex.RUnlock()

	pluginName, exists := r.extToPlugin[strings.ToLower(ext)]
	if !exists {
		return nil, false
	}
	return r.GetPlugin(pluginName)
}

// GetPluginByFilename determines which plugin to use based on filename
func (r *PluginRegistry) GetPluginByFilename(filename string) (ExtractorPlugin, bool) {
	r.initMutex.RLock()
	defer r.initMutex.RUnlock()

	ext := strings.ToLower(filepath.Ext(filename))
	if plugin, exists := r.GetPluginByExtension(ext); exists {
		return plugin, true
	}

	// Fallback to checking all plugins
	for _, plugin := range r.plugins {
		if plugin.ShouldProcess(filename) {
			return plugin, true
		}
	}

	return nil, false
}

// GetPluginByLanguage returns a plugin by language name
func (r *PluginRegistry) GetPluginByLanguage(lang string) (ExtractorPlugin, bool) {
	r.initMutex.RLock()
	defer r.initMutex.RUnlock()

	plugin, exists := r.plugins[strings.ToLower(lang)]
	return plugin, exists
}

// ListPlugins returns all registered plugin names
func (r *PluginRegistry) ListPlugins() []string {
	r.initMutex.RLock()
	defer r.initMutex.RUnlock()

	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}

	// Sort the names for consistent output
	sort.Strings(names)
	return names
}

// Unregister removes a plugin from the registry
func (r *PluginRegistry) Unregister(name string) bool {
	r.initMutex.Lock()
	defer r.initMutex.Unlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return false
	}

	// Clean up extension mappings
	for _, ext := range plugin.Extensions() {
		delete(r.extToPlugin, strings.ToLower(ext))
	}

	// Cleanup the plugin
	plugin.Cleanup()

	// Remove from registry
	delete(r.plugins, name)

	return true
}

// Cleanup all plugins
func (r *PluginRegistry) Cleanup() {
	r.initMutex.Lock()
	defer r.initMutex.Unlock()

	for _, plugin := range r.plugins {
		plugin.Cleanup()
	}

	// Clear maps
	r.plugins = make(map[string]ExtractorPlugin)
	r.extToPlugin = make(map[string]string)
}

// PluginLoader handles dynamic loading of plugins
type PluginLoader struct {
	registry *PluginRegistry
	plugins  map[string]string // name -> path
	mutex    sync.RWMutex
}

// NewPluginLoader creates a new plugin loader
func NewPluginLoader(registry *PluginRegistry) *PluginLoader {
	return &PluginLoader{
		registry: registry,
		plugins:  make(map[string]string),
	}
}

// LoadPlugin loads a plugin from a .so file
func (l *PluginLoader) LoadPlugin(pluginPath string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Check if plugin already loaded
	for _, path := range l.plugins {
		if path == pluginPath {
			return fmt.Errorf("plugin already loaded: %s", pluginPath)
		}
	}

	// Open the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %w", pluginPath, err)
	}

	// Look for the plugin symbol
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin %s doesn't export 'Plugin' symbol: %w", pluginPath, err)
	}

	// Cast to ExtractorPlugin
	extractor, ok := sym.(ExtractorPlugin)
	if !ok {
		return fmt.Errorf("plugin %s doesn't implement ExtractorPlugin interface", pluginPath)
	}

	// Register the plugin
	if err := l.registry.Register(extractor); err != nil {
		return fmt.Errorf("failed to register plugin %s: %w", pluginPath, err)
	}

	// Store plugin info
	pluginName := extractor.Name()
	l.plugins[pluginName] = pluginPath

	return nil
}

// LoadPluginsFromDir loads all plugins from a directory
func (l *PluginLoader) LoadPluginsFromDir(dirPath string) ([]string, error) {
	var loaded []string

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin directory %s: %w", dirPath, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) == ".so" {
			pluginPath := filepath.Join(dirPath, file.Name())
			if err := l.LoadPlugin(pluginPath); err != nil {
				fmt.Printf("Warning: Failed to load plugin %s: %v\n", pluginPath, err)
				continue
			}
			loaded = append(loaded, file.Name())
		}
	}

	return loaded, nil
}

// LoadPluginByLanguage loads a specific plugin by language name from a directory
func (l *PluginLoader) LoadPluginByLanguage(dirPath, language string) error {
	pluginFileName := fmt.Sprintf("%s_plugin.so", strings.ToLower(language))
	pluginPath := filepath.Join(dirPath, pluginFileName)

	return l.LoadPlugin(pluginPath)
}

// UnloadPlugin unloads a plugin
func (l *PluginLoader) UnloadPlugin(pluginName string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	_, exists := l.plugins[pluginName]
	if !exists {
		return fmt.Errorf("plugin not loaded: %s", pluginName)
	}

	// Note: Go plugins cannot be truly unloaded.
	// We just remove it from our registry.
	delete(l.plugins, pluginName)

	// Cleanup the plugin if possible
	if plugin, exists := l.registry.GetPlugin(pluginName); exists {
		plugin.Cleanup()
	}

	// Remove from registry
	l.registry.Unregister(pluginName)

	return nil
}

// ListLoadedPlugins returns list of loaded plugins
func (l *PluginLoader) ListLoadedPlugins() []string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	plugins := make([]string, 0, len(l.plugins))
	for name := range l.plugins {
		plugins = append(plugins, name)
	}

	return plugins
}

// ValidatePlugin checks if a plugin file is valid
func (l *PluginLoader) ValidatePlugin(pluginPath string) error {
	// Open the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %w", pluginPath, err)
	}

	// Look for the plugin symbol
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin %s doesn't export 'Plugin' symbol: %w", pluginPath, err)
	}

	// Cast to ExtractorPlugin
	_, ok := sym.(ExtractorPlugin)
	if !ok {
		return fmt.Errorf("plugin %s doesn't implement ExtractorPlugin interface", pluginPath)
	}

	return nil
}

// ReloadPlugin reloads a plugin from file
func (l *PluginLoader) ReloadPlugin(pluginName string, pluginPath string) error {
	// Unload the existing plugin
	if err := l.UnloadPlugin(pluginName); err != nil {
		return fmt.Errorf("failed to unload plugin: %w", err)
	}

	// Load the plugin again
	return l.LoadPlugin(pluginPath)
}

// PluginBuilder helps create and compile plugins
type PluginBuilder struct {
	pluginDir string
}

// NewPluginBuilder creates a new plugin builder
func NewPluginBuilder(pluginDir string) *PluginBuilder {
	return &PluginBuilder{
		pluginDir: pluginDir,
	}
}

// BuildPluginTemplate generates a template for a new plugin
func (b *PluginBuilder) BuildPluginTemplate(language, author string) (string, error) {
	// This function generates plugin templates but is currently disabled due to template complexity
	return "", fmt.Errorf("plugin template generation is temporarily disabled")
}

// BuildAdvancedPluginTemplate generates a more comprehensive template for a new plugin
func (b *PluginBuilder) BuildAdvancedPluginTemplate(language, author string) (string, error) {
	// This function generates advanced plugin templates but is currently disabled due to template complexity
	return "", fmt.Errorf("advanced plugin template generation is temporarily disabled")
}
