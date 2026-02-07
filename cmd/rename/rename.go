package rename

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// RenameCommand handles the rename subcommand
type RenameCommand struct {
	// Flags
	directory   string
	pattern     string
	prefix      string
	suffix      string
	regex       string
	replacement string
	verbose     bool
	quiet       bool
	dryRun      bool
	force       bool
	recursive   bool

	// Internal fields
	cyan   func(...interface{}) string
	green  func(...interface{}) string
	yellow func(...interface{}) string
	red    func(...interface{}) string
}

// NewRenameCommand creates a new rename command instance
func NewRenameCommand() *RenameCommand {
	return &RenameCommand{
		cyan:   color.New(color.FgCyan).SprintFunc(),
		green:  color.New(color.FgGreen).SprintFunc(),
		yellow: color.New(color.FgYellow).SprintFunc(),
		red:    color.New(color.FgRed).SprintFunc(),
	}
}

// Run executes the rename command
func (c *RenameCommand) Run(args []string) error {
	// Define flags
	fs := flag.NewFlagSet("rename", flag.ContinueOnError)
	fs.StringVar(&c.directory, "dir", ".", "Directory to rename files in")
	fs.StringVar(&c.pattern, "pattern", "", "Pattern to remove from filenames (prefix, suffix, or substring)")
	fs.StringVar(&c.prefix, "prefix", "", "Prefix to remove from filenames")
	fs.StringVar(&c.suffix, "suffix", "", "Suffix to remove from filenames")
	fs.StringVar(&c.regex, "regex", "", "Regular expression pattern to match")
	fs.StringVar(&c.replacement, "replacement", "", "Replacement string for regex (use with -regex)")
	fs.BoolVar(&c.verbose, "verbose", false, "Show detailed progress")
	fs.BoolVar(&c.quiet, "quiet", false, "Suppress non-essential output")
	fs.BoolVar(&c.dryRun, "dry-run", false, "Show what would be renamed without actually renaming")
	fs.BoolVar(&c.force, "force", false, "Force rename even if target file exists")
	fs.BoolVar(&c.recursive, "recursive", false, "Process subdirectories recursively")

	// Help flag
	help := fs.Bool("help", false, "Show help")
	h := fs.Bool("h", false, "Show help (shorthand)")

	// Parse flags
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Check for help
	if *help || *h {
		c.printHelp()
		return nil
	}

	// Validate directory exists
	if _, err := os.Stat(c.directory); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", c.directory)
	}

	// At least one renaming option must be specified
	hasOptions := c.pattern != "" || c.prefix != "" || c.suffix != "" || c.regex != ""
	if !hasOptions {
		return fmt.Errorf("at least one renaming option must be specified (-pattern, -prefix, -suffix, or -regex)")
	}

	// Validate regex if provided
	var compiledRegex *regexp.Regexp
	if c.regex != "" {
		var err error
		compiledRegex, err = regexp.Compile(c.regex)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %v", err)
		}
	}

	if !c.quiet {
		fmt.Printf("%s Starting rename operation\n", c.cyan("→"))
		fmt.Printf("%s Directory: %s\n", c.cyan("→"), c.directory)
		if c.pattern != "" {
			fmt.Printf("%s Pattern to remove: %s\n", c.cyan("→"), c.pattern)
		}
		if c.prefix != "" {
			fmt.Printf("%s Prefix to remove: %s\n", c.cyan("→"), c.prefix)
		}
		if c.suffix != "" {
			fmt.Printf("%s Suffix to remove: %s\n", c.cyan("→"), c.suffix)
		}
		if c.regex != "" {
			fmt.Printf("%s Regex pattern: %s\n", c.cyan("→"), c.regex)
			if c.replacement != "" {
				fmt.Printf("%s Replacement: %s\n", c.cyan("→"), c.replacement)
			}
		}
		if c.dryRun {
			fmt.Printf("%s DRY RUN MODE - No files will be renamed\n", c.yellow("⚠"))
		}
	}

	// Process files in the directory
	count, err := c.processDirectory(compiledRegex)
	if err != nil {
		return err
	}

	if !c.quiet {
		fmt.Printf("\n%s Renaming completed. %d files processed.\n", c.green("✓"), count)
	}

	return nil
}

// printHelp prints the help message
func (c *RenameCommand) printHelp() {
	fmt.Fprintf(os.Stderr, "%s Coto Rename v1.0.0 - Rename files based on patterns\n\n", c.cyan("📁"))
	fmt.Fprintf(os.Stderr, "Usage: coto rename [options]\n\n")

	fmt.Fprintf(os.Stderr, "%s Basic Options:\n", c.cyan("📋"))
	fmt.Fprintf(os.Stderr, "  -dir string            Directory to rename files in (default \".\")\n")
	fmt.Fprintf(os.Stderr, "  -pattern string        Pattern to remove from filenames\n")
	fmt.Fprintf(os.Stderr, "  -prefix string         Prefix to remove from filenames\n")
	fmt.Fprintf(os.Stderr, "  -suffix string         Suffix to remove from filenames\n")
	fmt.Fprintf(os.Stderr, "  -regex string          Regular expression pattern to match\n")
	fmt.Fprintf(os.Stderr, "  -replacement string    Replacement string for regex (use with -regex)\n")

	fmt.Fprintf(os.Stderr, "\n%s Mode Options:\n", c.cyan("🎯"))
	fmt.Fprintf(os.Stderr, "  -dry-run               Show what would be renamed without actually renaming\n")
	fmt.Fprintf(os.Stderr, "  -verbose               Show detailed progress\n")
	fmt.Fprintf(os.Stderr, "  -quiet                 Suppress non-essential output\n")
	fmt.Fprintf(os.Stderr, "  -force                 Force rename even if target file exists\n")
	fmt.Fprintf(os.Stderr, "  -recursive             Process subdirectories recursively\n")

	fmt.Fprintf(os.Stderr, "\n%s Information:\n", c.cyan("ℹ️"))
	fmt.Fprintf(os.Stderr, "  -h, -help              Show this help message\n")

	fmt.Fprintf(os.Stderr, "\n%s Examples:\n", c.cyan("🚀"))
	fmt.Fprintf(os.Stderr, "  coto rename -dir ./videos -pattern \"thisuffix\"\n")
	fmt.Fprintf(os.Stderr, "  coto rename -dir ./files -prefix \"old_\" -suffix \"_backup\"\n")
	fmt.Fprintf(os.Stderr, "  coto rename -dir ./data -regex \"^\\d+_(.*)\" -replacement \"$1\"\n")
	fmt.Fprintf(os.Stderr, "  coto rename -dir ./photos -suffix \".bak\" -dry-run\n")
}

// processDirectory processes all files in the specified directory and subdirectories if recursive is enabled
func (c *RenameCommand) processDirectory(regex *regexp.Regexp) (int, error) {
	count := 0
	
	if c.recursive {
		// Walk the directory tree recursively
		err := filepath.Walk(c.directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			// Skip directories
			if info.IsDir() {
				return nil
			}
			
			// Get the directory of the current file
			dir := filepath.Dir(path)
			filename := info.Name()
			
			newName := c.renameFile(filename, regex)
			
			// If the name didn't change, skip
			if newName == filename {
				if c.verbose && !c.quiet {
					fmt.Printf("%s Skipping %s (no change)\n", c.cyan("→"), filename)
				}
				return nil
			}
			
			newPath := filepath.Join(dir, newName)
			
			// Check if target file already exists
			if !c.force {
				if _, err := os.Stat(newPath); err == nil {
					if !c.quiet {
						fmt.Printf("%s Skipped %s -> %s (target already exists)\n", c.yellow("⚠"), filename, newName)
					}
					return nil
				}
			}
			
			if !c.quiet {
				fmt.Printf("%s %s -> %s\n", c.cyan("→"), filename, newName)
			}
			
			if !c.dryRun {
				if err := os.Rename(path, newPath); err != nil {
					if !c.quiet {
						fmt.Printf("%s Error renaming %s: %v\n", c.red("✗"), filename, err)
					}
					return nil // Continue processing other files
				}
			} else {
				if c.verbose && !c.quiet {
					fmt.Printf("%s Would rename %s -> %s\n", c.yellow("→"), filename, newName)
				}
			}
			
			count++
			
			return nil
		})
		
		if err != nil {
			return count, fmt.Errorf("error walking directory tree: %v", err)
		}
	} else {
		// Original behavior - only process the specified directory
		entries, err := os.ReadDir(c.directory)
		if err != nil {
			return 0, fmt.Errorf("failed to read directory: %v", err)
		}

		for _, entry := range entries {
			// Skip directories
			if entry.IsDir() {
				continue
			}

			oldPath := filepath.Join(c.directory, entry.Name())
			newName := c.renameFile(entry.Name(), regex)

			// If the name didn't change, skip
			if newName == entry.Name() {
				if c.verbose && !c.quiet {
					fmt.Printf("%s Skipping %s (no change)\n", c.cyan("→"), entry.Name())
				}
				continue
			}

			newPath := filepath.Join(c.directory, newName)

			// Check if target file already exists
			if !c.force {
				if _, err := os.Stat(newPath); err == nil {
					if !c.quiet {
						fmt.Printf("%s Skipped %s -> %s (target already exists)\n", c.yellow("⚠"), entry.Name(), newName)
					}
					continue
				}
			}

			if !c.quiet {
				fmt.Printf("%s %s -> %s\n", c.cyan("→"), entry.Name(), newName)
			}

			if !c.dryRun {
				if err := os.Rename(oldPath, newPath); err != nil {
					if !c.quiet {
						fmt.Printf("%s Error renaming %s: %v\n", c.red("✗"), entry.Name(), err)
					}
					continue
				}
			} else {
				if c.verbose && !c.quiet {
					fmt.Printf("%s Would rename %s -> %s\n", c.yellow("→"), entry.Name(), newName)
				}
			}

			count++
		}
	}

	return count, nil
}

// renameFile applies the renaming rules to a single filename
func (c *RenameCommand) renameFile(filename string, regex *regexp.Regexp) string {
	result := filename

	// Apply regex replacement first if specified
	if c.regex != "" && regex != nil {
		if c.replacement != "" {
			result = regex.ReplaceAllString(result, c.replacement)
		} else {
			// If no replacement is specified, just remove the matched pattern
			result = regex.ReplaceAllString(result, "")
		}
	}

	// Apply pattern removal (removes the pattern anywhere in the filename)
	if c.pattern != "" {
		result = strings.ReplaceAll(result, c.pattern, "")
	}

	// Apply prefix removal
	if c.prefix != "" && strings.HasPrefix(result, c.prefix) {
		result = strings.TrimPrefix(result, c.prefix)
	}

	// Apply suffix removal
	if c.suffix != "" && strings.HasSuffix(result, c.suffix) {
		result = strings.TrimSuffix(result, c.suffix)
	}

	return result
}
