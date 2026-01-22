package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)


// ConfigFile represents a configuration file for the extractor
type ConfigFile struct {
	InputFiles   []string `json:"input_files"`
	OutputDir    string   `json:"output_dir"`
	Language     string   `json:"language"`
	GenerateRepo bool     `json:"generate_report"`
	DryRun       bool     `json:"dry_run"`
	Parallel     int      `json:"parallel"`
	Verbose      bool     `json:"verbose"`
	Quiet        bool     `json:"quiet"`
	PluginDir    string   `json:"plugin_dir"`
	Plugins      []string `json:"plugins"`
	MaxFileSize  int64    `json:"max_file_size"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*ConfigFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ConfigFile
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to a JSON file
func SaveConfig(filename string, config *ConfigFile) error {
	content, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, content, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// FindFilesRecursive finds files matching patterns recursively
func FindFilesRecursive(dir string, patterns []string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file matches any pattern
		for _, pattern := range patterns {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				return err
			}
			if matched {
				files = append(files, path)
				break
			}
		}

		return nil
	})

	return files, err
}

// CreateDirectoryStructure creates standard directory structure
func CreateDirectoryStructure(baseDir string) error {
	dirs := []string{
		filepath.Join(baseDir, "src"),
		filepath.Join(baseDir, "config"),
		filepath.Join(baseDir, "scripts"),
		filepath.Join(baseDir, "docs"),
		filepath.Join(baseDir, "tests"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// DetectLanguageFromContent tries to detect language from file content
func DetectLanguageFromContent(content string) string {
	content = strings.ToLower(content)

	// Check for language signatures
	switch {
	case strings.Contains(content, "package ") && strings.Contains(content, "import "):
		return "go"
	case strings.Contains(content, "public class") || strings.Contains(content, "private class"):
		return "java"
	case strings.Contains(content, "def ") && strings.Contains(content, ":"):
		return "python"
	case strings.Contains(content, "function ") || strings.Contains(content, "const ") || strings.Contains(content, "let "):
		return "javascript"
	case strings.Contains(content, "<?xml"):
		return "xml"
	case strings.Contains(content, "{") && strings.Contains(content, "}"):
		// Check if it's JSON
		if strings.HasPrefix(strings.TrimSpace(content), "{") {
			return "json"
		}
		return "text"
	default:
		return "text"
	}
}

// MergeBlocks merges similar code blocks
func MergeBlocks(blocks []CodeBlock) []CodeBlock {
	if len(blocks) == 0 {
		return blocks
	}

	// Group by filename and language
	groups := make(map[string][]CodeBlock)
	for _, block := range blocks {
		key := block.Filename + "|" + block.Language
		groups[key] = append(groups[key], block)
	}

	// Merge each group
	var merged []CodeBlock
	for _, group := range groups {
		if len(group) == 1 {
			merged = append(merged, group[0])
			continue
		}

		// Merge multiple blocks into one
		mergedBlock := group[0]
		for i := 1; i < len(group); i++ {
			// Append content with separator
			mergedBlock.Content += "\n\n" + strings.Repeat("-", 80) + "\n\n" + group[i].Content

			// Merge imports
			mergedBlock.Imports = mergeSlices(mergedBlock.Imports, group[i].Imports)

			// Merge annotations
			mergedBlock.Annotations = mergeSlices(mergedBlock.Annotations, group[i].Annotations)

			// Merge modifiers
			mergedBlock.Modifiers = mergeSlices(mergedBlock.Modifiers, group[i].Modifiers)
		}

		merged = append(merged, mergedBlock)
	}

	return merged
}

func mergeSlices(slice1, slice2 []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range slice1 {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	for _, item := range slice2 {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// ValidateBlock validates a code block
func ValidateBlock(block CodeBlock) error {
	if block.Filename == "" {
		return fmt.Errorf("block has empty filename")
	}

	if block.Content == "" {
		return fmt.Errorf("block has empty content")
	}

	// Check for common issues
	if len(block.Content) > 10*1024*1024 { // 10MB
		return fmt.Errorf("block content too large: %d bytes", len(block.Content))
	}

	// Check for invalid characters in filename
	if strings.ContainsAny(block.Filename, `<>:"/\|?*`) {
		return fmt.Errorf("invalid characters in filename: %s", block.Filename)
	}

	return nil
}
