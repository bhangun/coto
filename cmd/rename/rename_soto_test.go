package rename

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunSearchReplaceRecursiveBothAndExclude(t *testing.T) {
	root := t.TempDir()
	writeRenameTestFile(t, filepath.Join(root, "old-root.txt"))
	writeRenameTestFile(t, filepath.Join(root, "old-dir", "old-file.txt"))
	writeRenameTestFile(t, filepath.Join(root, "old-dir", "excluded", "old-secret.txt"))
	writeRenameTestFile(t, filepath.Join(root, ".hidden", "old-hidden.txt"))

	cmd := NewRenameCommand()
	if err := cmd.Run([]string{
		"-d", root, "-s", "old", "-r", "new",
		"-e", "excluded", "-quiet",
	}); err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	for _, path := range []string{
		filepath.Join(root, "new-root.txt"),
		filepath.Join(root, "new-dir", "new-file.txt"),
		filepath.Join(root, "new-dir", "excluded", "old-secret.txt"),
		filepath.Join(root, ".hidden", "old-hidden.txt"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected path %q: %v", path, err)
		}
	}
	for _, path := range []string{
		filepath.Join(root, "old-root.txt"),
		filepath.Join(root, "old-dir"),
	} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("expected old path %q to be renamed, stat error: %v", path, err)
		}
	}
}

func TestRunRegexIgnoreCaseAndDepth(t *testing.T) {
	root := t.TempDir()
	writeRenameTestFile(t, filepath.Join(root, "OLD-file.txt"))
	writeRenameTestFile(t, filepath.Join(root, "nested", "OLD-nested.txt"))

	cmd := NewRenameCommand()
	if err := cmd.Run([]string{
		"-d", root, "--regex", "-s", "^old", "-r", "new",
		"-i", "-D", "1", "--quiet",
	}); err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(root, "new-file.txt")); err != nil {
		t.Errorf("expected depth-one file to be renamed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "nested", "OLD-nested.txt")); err != nil {
		t.Errorf("expected file deeper than max depth to remain: %v", err)
	}
}

func TestRunSotoDirAliasSelectsDirectories(t *testing.T) {
	root := t.TempDir()
	writeRenameTestFile(t, filepath.Join(root, "old-dir", "keep.txt"))

	cmd := NewRenameCommand()
	if err := cmd.Run([]string{
		"-d", root, "-s", "old", "-r", "new", "--dir", "--quiet",
	}); err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "new-dir", "keep.txt")); err != nil {
		t.Errorf("expected directory to be renamed: %v", err)
	}
}

func TestRunRegexUsesSotoReplacementSyntax(t *testing.T) {
	root := t.TempDir()
	writeRenameTestFile(t, filepath.Join(root, "old_name.txt"))

	cmd := NewRenameCommand()
	if err := cmd.Run([]string{
		"-d", root, "--regex", "-s", "^(old)_(.*)$",
		"-r", `\2-\1-&`, "--file", "--no-recursive", "--quiet",
	}); err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "name.txt-old-old_name.txt")); err != nil {
		t.Errorf("expected regex replacement using Soto syntax: %v", err)
	}
}

func writeRenameTestFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("create test directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("create test file: %v", err)
	}
}
