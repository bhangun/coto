package rename

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// RenameCommand handles the rename subcommand.
type RenameCommand struct {
	directory       string
	pattern         string
	prefix          string
	suffix          string
	regex           string
	search          string
	replacement     string
	verbose         bool
	quiet           bool
	dryRun          bool
	force           bool
	recursive       bool
	maxDepth        int
	maxDepthSet     bool
	targetType      string
	ignoreCase      bool
	includeHidden   bool
	regexMode       bool
	noRecursive     bool
	excludePatterns stringList

	cyan   func(...interface{}) string
	green  func(...interface{}) string
	yellow func(...interface{}) string
	red    func(...interface{}) string
}

type stringList []string

func (s *stringList) String() string { return strings.Join(*s, ",") }

func (s *stringList) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// NewRenameCommand creates a new rename command instance.
func NewRenameCommand() *RenameCommand {
	return &RenameCommand{
		maxDepth:   -1,
		targetType: "files",
		cyan:       color.New(color.FgCyan).SprintFunc(),
		green:      color.New(color.FgGreen).SprintFunc(),
		yellow:     color.New(color.FgYellow).SprintFunc(),
		red:        color.New(color.FgRed).SprintFunc(),
	}
}

// Run executes the rename command.
func (c *RenameCommand) Run(args []string) error {
	fs := flag.NewFlagSet("rename", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.StringVar(&c.directory, "dir", ".", "Directory to process (legacy alias)")
	fs.StringVar(&c.directory, "directory", ".", "Starting directory")
	fs.StringVar(&c.directory, "d", ".", "Starting directory (shorthand)")
	fs.StringVar(&c.pattern, "pattern", "", "Pattern to remove from filenames")
	fs.StringVar(&c.prefix, "prefix", "", "Prefix to remove from filenames")
	fs.StringVar(&c.suffix, "suffix", "", "Suffix to remove from filenames")
	fs.StringVar(&c.regex, "regex-pattern", "", "Regular expression pattern to match (legacy rename option)")
	fs.StringVar(&c.regex, "regex", "", "Regular expression pattern to match (legacy rename option)")
	fs.StringVar(&c.search, "s", "", "Search text or regular expression (use with -r)")
	fs.StringVar(&c.search, "search", "", "Search text or regular expression (use with -r)")
	fs.StringVar(&c.replacement, "r", "", "Replacement string (use with -s)")
	fs.StringVar(&c.replacement, "replace", "", "Replacement string (use with -s)")
	fs.StringVar(&c.replacement, "replacement", "", "Replacement string for -regex")
	fs.BoolVar(&c.regexMode, "regex-mode", false, "Treat -s/--search as a regular expression")
	fs.BoolVar(&c.verbose, "verbose", false, "Show detailed progress")
	fs.BoolVar(&c.verbose, "v", false, "Show detailed progress")
	fs.BoolVar(&c.quiet, "quiet", false, "Suppress non-essential output")
	fs.BoolVar(&c.dryRun, "dry-run", false, "Show what would be renamed without changing anything")
	fs.BoolVar(&c.force, "force", false, "Allow an existing target to be replaced")
	fs.BoolVar(&c.recursive, "recursive", false, "Process subdirectories recursively")
	fs.BoolVar(&c.noRecursive, "no-recursive", false, "Process only the starting directory")
	fs.IntVar(&c.maxDepth, "max-depth", -1, "Maximum depth below the starting directory (-1 is unlimited)")
	fs.IntVar(&c.maxDepth, "D", -1, "Maximum depth below the starting directory (-1 is unlimited)")
	fs.BoolVar(&c.ignoreCase, "ignore-case", false, "Match search text without regard to case")
	fs.BoolVar(&c.ignoreCase, "i", false, "Match search text without regard to case")
	fs.BoolVar(&c.includeHidden, "hidden", false, "Include hidden files and directories")
	fs.BoolVar(&c.includeHidden, "H", false, "Include hidden files and directories")
	fs.Var(&c.excludePatterns, "exclude", "Exclude paths matching a glob pattern (repeatable)")
	fs.Var(&c.excludePatterns, "e", "Exclude paths matching a glob pattern (repeatable)")

	var filesOnly, dirsOnly, both bool
	fs.BoolVar(&filesOnly, "file", false, "Rename files only (default for legacy Coto options)")
	fs.BoolVar(&filesOnly, "files", false, "Rename files only (default for legacy Coto options)")
	fs.BoolVar(&dirsOnly, "directories", false, "Rename directories only")
	fs.BoolVar(&dirsOnly, "dirs", false, "Rename directories only")
	fs.BoolVar(&both, "both", false, "Rename both files and directories")
	help := fs.Bool("help", false, "Show help")
	shortHelp := fs.Bool("h", false, "Show help (shorthand)")

	// Soto used --regex as a switch alongside -s. Retain Coto's existing
	// --regex PATTERN form while accepting Soto's switch-style spelling.
	args = normalizeRenameArgs(args)
	if err := fs.Parse(args); err != nil {
		return err
	}
	c.maxDepthSet = hasFlag(args, "-D", "--max-depth")
	if *help || *shortHelp {
		c.printHelp()
		return nil
	}

	sotoMode := hasFlag(args, "-s", "--search") || hasFlag(args, "-r", "--replace")
	if filesOnly && dirsOnly || (both && (filesOnly || dirsOnly)) {
		return fmt.Errorf("choose only one of --file, --directories, or --both")
	}
	switch {
	case dirsOnly:
		c.targetType = "dirs"
	case both:
		c.targetType = "both"
	default:
		if sotoMode {
			c.targetType = "both"
		} else {
			c.targetType = "files"
		}
	}
	if c.maxDepth < -1 {
		return fmt.Errorf("maximum depth must be -1 or greater")
	}
	if c.maxDepth >= 0 && c.maxDepth != -1 {
		c.recursive = true
	}
	if c.noRecursive {
		c.recursive = false
	} else if sotoMode {
		c.recursive = true
	}

	info, err := os.Stat(c.directory)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", c.directory)
		}
		return fmt.Errorf("cannot access directory %q: %w", c.directory, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", c.directory)
	}

	if sotoMode {
		if !hasFlag(args, "-s", "--search") || !hasFlag(args, "-r", "--replace") {
			return fmt.Errorf("both -s/--search and -r/--replace are required")
		}
		if c.search == "" {
			return fmt.Errorf("search pattern cannot be empty")
		}
	} else if c.pattern == "" && c.prefix == "" && c.suffix == "" && c.regex == "" {
		return fmt.Errorf("at least one renaming option must be specified (-pattern, -prefix, -suffix, or -regex)")
	}

	compiledRegex, err := c.compileSearchRegex(sotoMode)
	if err != nil {
		return err
	}
	for _, pattern := range c.excludePatterns {
		if _, err := filepath.Match(pattern, ""); err != nil {
			return fmt.Errorf("invalid exclude pattern %q: %w", pattern, err)
		}
	}

	if !c.quiet {
		fmt.Printf("%s Starting rename operation\n", c.cyan("→"))
		fmt.Printf("%s Directory: %s\n", c.cyan("→"), c.directory)
		if sotoMode {
			fmt.Printf("%s Search: %s\n", c.cyan("→"), c.search)
			fmt.Printf("%s Replacement: %s\n", c.cyan("→"), c.replacement)
		}
		fmt.Printf("%s Target: %s\n", c.cyan("→"), c.targetType)
		if c.dryRun {
			fmt.Printf("%s DRY RUN MODE - No paths will be renamed\n", c.yellow("⚠"))
		}
	}

	count, err := c.processDirectory(compiledRegex)
	if err != nil {
		return err
	}
	if !c.quiet {
		fmt.Printf("\n%s Renaming completed. %d path(s) renamed.\n", c.green("✓"), count)
	}
	return nil
}

func normalizeRenameArgs(args []string) []string {
	normalized := make([]string, 0, len(args))
	for i, arg := range args {
		if arg == "--dir" && (i+1 == len(args) || strings.HasPrefix(args[i+1], "-")) {
			normalized = append(normalized, "--directories")
			continue
		}
		if arg == "--regex" && (i+1 == len(args) || strings.HasPrefix(args[i+1], "-")) {
			normalized = append(normalized, "--regex-mode")
			continue
		}
		normalized = append(normalized, arg)
	}
	return normalized
}

func hasFlag(args []string, names ...string) bool {
	for _, arg := range args {
		for _, name := range names {
			if arg == name || strings.HasPrefix(arg, name+"=") {
				return true
			}
		}
	}
	return false
}

func (c *RenameCommand) compileSearchRegex(sotoMode bool) (*regexp.Regexp, error) {
	expression := c.regex
	if sotoMode && c.regexMode {
		expression = c.search
	}
	if expression == "" {
		return nil, nil
	}
	if c.ignoreCase {
		expression = "(?i)" + expression
	}
	compiled, err := regexp.Compile(expression)
	if err != nil {
		if sotoMode && c.regexMode {
			return nil, fmt.Errorf("invalid search regular expression: %w", err)
		}
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}
	return compiled, nil
}

func (c *RenameCommand) printHelp() {
	fmt.Fprintf(os.Stderr, "%s Coto Rename - Rename files and directories\n\n", c.cyan("📁"))
	fmt.Fprintln(os.Stderr, "Usage: coto rename [options]")
	fmt.Fprintln(os.Stderr, "\nSearch and replace (Soto features):")
	fmt.Fprintln(os.Stderr, "  -s, --search string     Search text or regex (requires -r)")
	fmt.Fprintln(os.Stderr, "  -r, --replace string    Replacement text; empty text is allowed")
	fmt.Fprintln(os.Stderr, "      --regex             Use regex for -s (or use --regex PATTERN for legacy Coto syntax)")
	fmt.Fprintln(os.Stderr, "  -i, --ignore-case       Match without regard to case")
	fmt.Fprintln(os.Stderr, "  -d, --directory string Starting directory (default \".\"); -dir remains supported")
	fmt.Fprintln(os.Stderr, "  -D, --max-depth int     Maximum depth below the starting directory")
	fmt.Fprintln(os.Stderr, "      --recursive         Process subdirectories recursively (default for -s/-r)")
	fmt.Fprintln(os.Stderr, "      --no-recursive      Process only the starting directory")
	fmt.Fprintln(os.Stderr, "      --file              Rename files only (default for legacy Coto options)")
	fmt.Fprintln(os.Stderr, "      --directories       Rename directories only")
	fmt.Fprintln(os.Stderr, "      --both              Rename files and directories (default for -s/-r)")
	fmt.Fprintln(os.Stderr, "  -e, --exclude pattern   Exclude matching paths (repeatable)")
	fmt.Fprintln(os.Stderr, "  -H, --hidden            Include hidden files and directories")
	fmt.Fprintln(os.Stderr, "\nExisting Coto options:")
	fmt.Fprintln(os.Stderr, "      --pattern string    Remove a substring")
	fmt.Fprintln(os.Stderr, "      --prefix string     Remove a filename prefix")
	fmt.Fprintln(os.Stderr, "      --suffix string     Remove a filename suffix")
	fmt.Fprintln(os.Stderr, "      --regex PATTERN     Regex pattern (use --replacement for replacement)")
	fmt.Fprintln(os.Stderr, "      --dry-run           Preview without changing paths")
	fmt.Fprintln(os.Stderr, "      --force             Allow an existing target to be replaced")
	fmt.Fprintln(os.Stderr, "      --verbose           Show detailed progress")
	fmt.Fprintln(os.Stderr, "      --quiet             Suppress non-essential output")
	fmt.Fprintln(os.Stderr, "  -h, --help              Show this help")
	fmt.Fprintln(os.Stderr, "\nExamples:")
	fmt.Fprintln(os.Stderr, "  coto rename -d ./files -s old -r new --recursive --dry-run")
	fmt.Fprintln(os.Stderr, "  coto rename -d ./media --regex -s '^test_' -r prod_ --both")
	fmt.Fprintln(os.Stderr, "  coto rename -dir ./videos -prefix old_ --recursive")
}

type renameCandidate struct {
	path  string
	depth int
	dir   bool
}

func (c *RenameCommand) processDirectory(regex *regexp.Regexp) (int, error) {
	root, err := filepath.Abs(c.directory)
	if err != nil {
		return 0, fmt.Errorf("resolve starting directory: %w", err)
	}
	maxDepth := c.maxDepth
	if !c.maxDepthSet && maxDepth == 0 {
		maxDepth = -1
	}

	var candidates []renameCandidate
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}

		relative, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("resolve path %q: %w", path, err)
		}
		depth := strings.Count(filepath.Clean(relative), string(filepath.Separator)) + 1
		if maxDepth >= 0 && depth > maxDepth {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !c.includeHidden && hasHiddenComponent(relative) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if c.isExcluded(relative) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			if c.targetType == "dirs" || c.targetType == "both" {
				candidates = append(candidates, renameCandidate{path: path, depth: depth, dir: true})
			}
			if !c.recursive {
				return filepath.SkipDir
			}
			return nil
		}
		if c.targetType == "" || c.targetType == "files" || c.targetType == "both" {
			if entry.Type().IsRegular() {
				candidates = append(candidates, renameCandidate{path: path, depth: depth})
			}
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("walk directory tree: %w", err)
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].dir != candidates[j].dir {
			return !candidates[i].dir
		}
		if candidates[i].dir && candidates[i].depth != candidates[j].depth {
			return candidates[i].depth > candidates[j].depth
		}
		return candidates[i].path < candidates[j].path
	})

	renamed := 0
	for _, candidate := range candidates {
		changed, err := c.renamePath(candidate, regex)
		if err != nil {
			return renamed, err
		}
		if changed {
			renamed++
		}
	}
	return renamed, nil
}

func hasHiddenComponent(path string) bool {
	for _, component := range strings.Split(filepath.Clean(path), string(filepath.Separator)) {
		if strings.HasPrefix(component, ".") {
			return true
		}
	}
	return false
}

func (c *RenameCommand) isExcluded(relative string) bool {
	return c.matchesAnyPath(relative, c.excludePatterns)
}

func (c *RenameCommand) matchesAnyPath(relative string, patterns []string) bool {
	relative = filepath.Clean(relative)
	parts := strings.Split(relative, string(filepath.Separator))
	for _, pattern := range patterns {
		pattern = filepath.Clean(pattern)
		for start := range parts {
			candidate := filepath.Join(parts[start:]...)
			if matched, _ := filepath.Match(pattern, candidate); matched {
				return true
			}
		}
	}
	return false
}

func (c *RenameCommand) renamePath(candidate renameCandidate, regex *regexp.Regexp) (bool, error) {
	oldName := filepath.Base(candidate.path)
	newName := c.renameFile(oldName, regex)
	if c.search != "" && regex == nil {
		if c.ignoreCase {
			var err error
			newName, err = replaceAllFold(newName, c.search, c.replacement)
			if err != nil {
				return false, fmt.Errorf("build case-insensitive search: %w", err)
			}
		} else {
			newName = strings.ReplaceAll(newName, c.search, c.replacement)
		}
	} else if c.search != "" && c.regexMode && regex != nil {
		newName = regex.ReplaceAllString(newName, sedReplacementToGo(c.replacement))
	}
	if newName == oldName {
		if c.verbose && !c.quiet {
			fmt.Printf("%s Skipping %s (no change)\n", c.cyan("→"), candidate.path)
		}
		return false, nil
	}
	if newName == "" {
		if c.verbose && !c.quiet {
			fmt.Printf("%s Skipping %s (result would be an empty name)\n", c.yellow("⚠"), candidate.path)
		}
		return false, nil
	}
	if newName == "." || newName == ".." || strings.ContainsAny(newName, `/\`) || strings.ContainsRune(newName, '\x00') {
		return false, fmt.Errorf("refusing invalid new name %q for %q", newName, candidate.path)
	}

	target := filepath.Join(filepath.Dir(candidate.path), newName)
	if _, err := os.Lstat(target); err == nil && !c.force {
		if !c.quiet {
			fmt.Printf("%s Skipped %s -> %s (target already exists)\n", c.yellow("⚠"), candidate.path, target)
		}
		return false, nil
	} else if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("check target %q: %w", target, err)
	}
	if !c.quiet {
		if c.dryRun {
			fmt.Printf("%s Would rename %s -> %s\n", c.yellow("→"), candidate.path, target)
		} else {
			fmt.Printf("%s %s -> %s\n", c.cyan("→"), candidate.path, target)
		}
	}
	if c.dryRun {
		return true, nil
	}
	if err := os.Rename(candidate.path, target); err != nil {
		return false, fmt.Errorf("rename %q to %q: %w", candidate.path, target, err)
	}
	return true, nil
}

func sedReplacementToGo(replacement string) string {
	var translated strings.Builder
	for i := 0; i < len(replacement); i++ {
		switch replacement[i] {
		case '\\':
			if i+1 == len(replacement) {
				translated.WriteByte('\\')
				continue
			}
			i++
			next := replacement[i]
			if next >= '0' && next <= '9' {
				start := i
				for i+1 < len(replacement) && replacement[i+1] >= '0' && replacement[i+1] <= '9' {
					i++
				}
				translated.WriteByte('$')
				translated.WriteString(replacement[start : i+1])
			} else {
				translated.WriteByte(next)
			}
		case '&':
			translated.WriteString("$0")
		case '$':
			translated.WriteString("$$")
		default:
			translated.WriteByte(replacement[i])
		}
	}
	return translated.String()
}

func replaceAllFold(input, search, replacement string) (string, error) {
	if search == "" {
		return input, nil
	}
	re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(search))
	if err != nil {
		return "", err
	}
	indices := re.FindAllStringIndex(input, -1)
	if len(indices) == 0 {
		return input, nil
	}
	var result strings.Builder
	last := 0
	for _, index := range indices {
		result.WriteString(input[last:index[0]])
		result.WriteString(replacement)
		last = index[1]
	}
	result.WriteString(input[last:])
	return result.String(), nil
}

// renameFile applies Coto's prefix, suffix, substring, and regex rules.
func (c *RenameCommand) renameFile(filename string, regex *regexp.Regexp) string {
	result := filename
	if c.regex != "" && regex != nil {
		if c.replacement != "" {
			result = regex.ReplaceAllString(result, c.replacement)
		} else {
			result = regex.ReplaceAllString(result, "")
		}
	}
	if c.pattern != "" {
		result = strings.ReplaceAll(result, c.pattern, "")
	}
	if c.prefix != "" && strings.HasPrefix(result, c.prefix) {
		result = strings.TrimPrefix(result, c.prefix)
	}
	if c.suffix != "" && strings.HasSuffix(result, c.suffix) {
		result = strings.TrimSuffix(result, c.suffix)
	}
	return result
}
