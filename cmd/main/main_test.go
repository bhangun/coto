package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShouldProcessFileIncludeFolderAndExcludeChildFolder(t *testing.T) {
	root := t.TempDir()
	config := Config{
		InputDir:       root,
		IncludePattern: "lib",
		ExcludePattern: "child2",
	}

	tests := []struct {
		path string
		want bool
	}{
		{path: filepath.Join(root, "lib", "root.go"), want: true},
		{path: filepath.Join(root, "lib", "child1", "file.go"), want: true},
		{path: filepath.Join(root, "lib", "child1", "child2", "nested.go"), want: false},
		{path: filepath.Join(root, "src", "child1", "file.go"), want: false},
		{path: filepath.Join(root, "library", "file.go"), want: false},
		{path: filepath.Join(root, "lib.go"), want: false},
	}
	for _, test := range tests {
		t.Run(filepath.Base(filepath.Dir(test.path))+"/"+filepath.Base(test.path), func(t *testing.T) {
			info := testFileInfo(t, test.path)
			got := shouldProcessFile(test.path, info, config, nil, nil)
			if got != test.want {
				t.Errorf("shouldProcessFile(%q) = %v, want %v", test.path, got, test.want)
			}
		})
	}
}

func TestShouldProcessFileFiltersWorkIndependently(t *testing.T) {
	root := t.TempDir()
	includeOnly := Config{InputDir: root, IncludePattern: "lib"}
	includePath := filepath.Join(root, "lib", "child1", "file.go")
	if !shouldProcessFile(includePath, testFileInfo(t, includePath), includeOnly, nil, nil) {
		t.Fatal("expected include-only filter to accept a descendant")
	}
	if shouldProcessFile(filepath.Join(root, "src", "file.go"), testFileInfo(t, filepath.Join(root, "src", "file.go")), includeOnly, nil, nil) {
		t.Fatal("expected include-only filter to reject paths outside the selected folder")
	}

	excludeOnly := Config{InputDir: root, ExcludePattern: "child2"}
	excludedPath := filepath.Join(root, "lib", "child1", "child2", "file.go")
	if shouldProcessFile(excludedPath, testFileInfo(t, excludedPath), excludeOnly, nil, nil) {
		t.Fatal("expected exclude-only filter to reject the matching folder subtree")
	}
	allowedPath := filepath.Join(root, "lib", "child1", "file.go")
	if !shouldProcessFile(allowedPath, testFileInfo(t, allowedPath), excludeOnly, nil, nil) {
		t.Fatal("expected exclude-only filter to preserve other paths")
	}
}

func TestShouldProcessFilePathGlobAndRegexFilters(t *testing.T) {
	root := t.TempDir()
	globConfig := Config{InputDir: root, IncludePattern: "src/*.go"}
	globPath := filepath.Join(root, "src", "main.go")
	globFilter, err := compilePathFilter(globConfig.IncludePattern)
	if err != nil {
		t.Fatalf("compile glob filter: %v", err)
	}
	if !shouldProcessFile(globPath, testFileInfo(t, globPath), globConfig, nil, globFilter) {
		t.Fatal("expected glob path to be included")
	}

	regexConfig := Config{InputDir: root, IncludePattern: `\.go$`}
	regexFilter, err := compilePathFilter(regexConfig.IncludePattern)
	if err != nil {
		t.Fatalf("compile regex filter: %v", err)
	}
	regexPath := filepath.Join(root, "pkg", "main.go")
	if !shouldProcessFile(regexPath, testFileInfo(t, regexPath), regexConfig, nil, regexFilter) {
		t.Fatal("expected regex-matched file to be included")
	}
	textPath := filepath.Join(root, "pkg", "README.md")
	if shouldProcessFile(textPath, testFileInfo(t, textPath), regexConfig, nil, regexFilter) {
		t.Fatal("expected non-matching file to be excluded")
	}

	legacyConfig := Config{InputDir: root, ExcludePattern: `\.git|vendor`}
	excludeRegex, err := compilePathFilter(legacyConfig.ExcludePattern)
	if err != nil {
		t.Fatalf("compile legacy regex filter: %v", err)
	}
	vendorPath := filepath.Join(root, "pkg", "vendor", "file.go")
	if shouldProcessFile(vendorPath, testFileInfo(t, vendorPath), legacyConfig, excludeRegex, nil) {
		t.Fatal("expected legacy regex alternation to exclude vendor files")
	}

	if _, err := compilePathFilter("["); err == nil {
		t.Fatal("expected invalid path glob to be rejected")
	}
}

func testFileInfo(t *testing.T, path string) os.FileInfo {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("create test directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("create test file: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat test file: %v", err)
	}
	return info
}
