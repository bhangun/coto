package extract

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/bhangun/coto/pkg/extractor"
)

// ExtractCommand handles the extract subcommand
type ExtractCommand struct {
	// Flags
	inputFiles    string
	outputDir     string
	language      string
	parallel      int
	verbose       bool
	quiet         bool
	dryRun        bool
	report        bool
	listPlugins   bool
	pluginDir     string
	configFile    string

	// Internal fields
	cyan   func(...interface{}) string
	green  func(...interface{}) string
	yellow func(...interface{}) string
	red    func(...interface{}) string
}

// NewExtractCommand creates a new extract command instance
func NewExtractCommand() *ExtractCommand {
	return &ExtractCommand{
		cyan:   color.New(color.FgCyan).SprintFunc(),
		green:  color.New(color.FgGreen).SprintFunc(),
		yellow: color.New(color.FgYellow).SprintFunc(),
		red:    color.New(color.FgRed).SprintFunc(),
	}
}

// Run executes the extract command
func (c *ExtractCommand) Run(args []string) error {
	// Define flags
	fs := flag.NewFlagSet("extract", flag.ContinueOnError)
	fs.StringVar(&c.inputFiles, "input", "", "Comma-separated input files")
	fs.StringVar(&c.outputDir, "output", "extracted", "Output directory")
	fs.StringVar(&c.language, "language", "", "Target language (auto-detected)")
	fs.IntVar(&c.parallel, "parallel", 1, "Parallel processing")
	fs.BoolVar(&c.verbose, "verbose", false, "Show detailed progress")
	fs.BoolVar(&c.quiet, "quiet", false, "Suppress non-essential output")
	fs.BoolVar(&c.dryRun, "dry-run", false, "Show what would be extracted")
	fs.BoolVar(&c.report, "report", false, "Generate detailed report")
	fs.BoolVar(&c.listPlugins, "list-plugins", false, "List available plugins")
	fs.StringVar(&c.pluginDir, "plugin-dir", "", "Directory containing plugins")
	fs.StringVar(&c.configFile, "config", "", "Configuration file path")

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

	// Check for list-plugins
	if c.listPlugins {
		c.listAvailablePlugins()
		return nil
	}

	// Validate required arguments
	if c.inputFiles == "" && len(fs.Args()) == 0 {
		return fmt.Errorf("input files are required")
	}

	// Prepare input files
	var inputPaths []string
	if c.inputFiles != "" {
		// Split comma-separated input files
		inputPaths = strings.Split(c.inputFiles, ",")
		for i, path := range inputPaths {
			inputPaths[i] = strings.TrimSpace(path)
		}
	} else {
		// Use remaining arguments as input files
		inputPaths = fs.Args()
	}

	// Expand glob patterns
	var expandedPaths []string
	for _, path := range inputPaths {
		matches, err := filepath.Glob(path)
		if err != nil {
			return fmt.Errorf("invalid glob pattern: %s", path)
		}
		if matches != nil {
			expandedPaths = append(expandedPaths, matches...)
		} else {
			// If no matches, treat as literal path
			expandedPaths = append(expandedPaths, path)
		}
	}

	// Validate input files exist
	for _, path := range expandedPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("input file does not exist: %s", path)
		}
	}

	// Create output directory if not in dry-run mode
	if !c.dryRun {
		if err := os.MkdirAll(c.outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	if !c.quiet {
		fmt.Printf("%s Starting extraction\n", c.cyan("→"))
		fmt.Printf("%s Input files: %d\n", c.cyan("→"), len(expandedPaths))
		fmt.Printf("%s Output directory: %s\n", c.cyan("→"), c.outputDir)
		if c.dryRun {
			fmt.Printf("%s DRY RUN MODE - No files will be written\n", c.yellow("⚠"))
		}
	}

	// Initialize plugin registry
	registry := NewPluginRegistry()

	// Register built-in plugins
	registerBuiltInPlugins(registry)

	// Load custom plugins if plugin directory is specified
	if c.pluginDir != "" {
		if err := c.loadCustomPlugins(registry); err != nil {
			return fmt.Errorf("failed to load custom plugins: %v", err)
		}
	}

	// Process files
	results, err := c.processFiles(expandedPaths, registry)
	if err != nil {
		return err
	}

	// Generate report if requested
	if c.report {
		c.generateReport(results)
	}

	if !c.quiet {
		// Calculate statistics
		totalBlocks := 0
		totalFilesWritten := 0
		for _, result := range results {
			totalBlocks += len(result.CodeBlocks)
			totalFilesWritten += len(result.WrittenFiles)
		}

		fmt.Printf("\n%s %s\n", c.cyan("┌"), strings.Repeat("─", 50))
		fmt.Printf("%s Extraction Summary\n", c.cyan("│"))
		fmt.Printf("%s %s\n", c.cyan("├"), strings.Repeat("─", 50))
		fmt.Printf("%s Files processed:     %s\n", c.cyan("│"), c.green(strconv.Itoa(len(results))))
		fmt.Printf("%s Code blocks found:   %s\n", c.cyan("│"), c.green(strconv.Itoa(totalBlocks)))
		fmt.Printf("%s Files written:       %s\n", c.cyan("│"), c.green(strconv.Itoa(totalFilesWritten)))
		if c.language != "" {
			fmt.Printf("%s Target language:     %s\n", c.cyan("│"), c.green(c.language))
		}
		fmt.Printf("%s %s\n", c.cyan("└"), strings.Repeat("─", 50))

		fmt.Printf("\n%s Extraction completed successfully!\n", c.green("✓"))
	}

	return nil
}

// printHelp prints the help message
func (c *ExtractCommand) printHelp() {
	fmt.Fprintf(os.Stderr, "%s Coto Extract v1.0.0 - Extract code blocks\n\n", c.cyan("📁"))
	fmt.Fprintf(os.Stderr, "Usage: coto extract [options]\n\n")

	fmt.Fprintf(os.Stderr, "%s Basic Options:\n", c.cyan("📋"))
	fmt.Fprintf(os.Stderr, "  -input string        Comma-separated input files\n")
	fmt.Fprintf(os.Stderr, "  -output string       Output directory (default \"extracted\")\n")
	fmt.Fprintf(os.Stderr, "  -language string     Target language (auto-detected)\n")
	fmt.Fprintf(os.Stderr, "  -parallel int        Parallel processing (default 1)\n")

	fmt.Fprintf(os.Stderr, "\n%s Mode Options:\n", c.cyan("🎯"))
	fmt.Fprintf(os.Stderr, "  -dry-run             Show what would be extracted\n")
	fmt.Fprintf(os.Stderr, "  -verbose             Show detailed progress\n")
	fmt.Fprintf(os.Stderr, "  -quiet               Suppress non-essential output\n")
	fmt.Fprintf(os.Stderr, "  -report              Generate detailed report\n")

	fmt.Fprintf(os.Stderr, "\n%s Plugin Options:\n", c.cyan("🔌"))
	fmt.Fprintf(os.Stderr, "  -list-plugins        List available plugins\n")
	fmt.Fprintf(os.Stderr, "  -plugin-dir string   Directory containing plugins\n")

	fmt.Fprintf(os.Stderr, "\n%s Information:\n", c.cyan("ℹ️"))
	fmt.Fprintf(os.Stderr, "  -h, -help            Show this help message\n")

	fmt.Fprintf(os.Stderr, "\n%s Examples:\n", c.cyan("🚀"))
	fmt.Fprintf(os.Stderr, "  coto extract -input file1.txt,file2.txt -language java -output extracted/\n")
	fmt.Fprintf(os.Stderr, "  coto extract -input *.txt -language python -parallel 4 -verbose\n")
	fmt.Fprintf(os.Stderr, "  coto extract -input code.txt -language javascript -report -dry-run\n")
	fmt.Fprintf(os.Stderr, "  coto extract -list-plugins\n")
	fmt.Fprintf(os.Stderr, "  coto extract -input files/ -plugin-dir ./plugins/\n")
}

// listAvailablePlugins lists all registered plugins
func (c *ExtractCommand) listAvailablePlugins() {
	registry := NewPluginRegistry()
	registerBuiltInPlugins(registry)

	fmt.Printf("%s Available Extractor Plugins:\n", c.cyan("🔌"))
	
	plugins := registry.GetAllPlugins()
	for _, plugin := range plugins {
		fmt.Printf("  %s: %s\n", c.green(plugin.Name()), strings.Join(plugin.Extensions(), ", "))
	}
}

// loadCustomPlugins loads plugins from the specified directory
func (c *ExtractCommand) loadCustomPlugins(registry *PluginRegistry) error {
	// For now, we'll just scan the directory for plugin files
	// In a real implementation, this would dynamically load Go plugins
	files, err := os.ReadDir(c.pluginDir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			// In a real implementation, we would load the plugin here
			if !c.quiet {
				fmt.Printf("%s Loading plugin: %s\n", c.cyan("→"), file.Name())
			}
		}
	}

	return nil
}

// processFiles processes the input files with the given registry
func (c *ExtractCommand) processFiles(inputPaths []string, registry *PluginRegistry) ([]ExtractionResult, error) {
	var results []ExtractionResult

	if c.parallel > 1 {
		// Use parallel processing
		results = c.processFilesParallel(inputPaths, registry)
	} else {
		// Use sequential processing
		results = c.processFilesSequential(inputPaths, registry)
	}

	return results, nil
}

// processFilesSequential processes files sequentially
func (c *ExtractCommand) processFilesSequential(inputPaths []string, registry *PluginRegistry) []ExtractionResult {
	var results []ExtractionResult

	for i, path := range inputPaths {
		if c.verbose && !c.quiet {
			fmt.Printf("%s Processing file %d/%d: %s\n",
				c.cyan("↳"), i+1, len(inputPaths), path)
		}

		result, err := c.processSingleFile(path, registry)
		if err != nil {
			if !c.quiet {
				fmt.Printf("%s Error processing %s: %v\n", c.red("✗"), path, err)
			}
			continue
		}

		results = append(results, result)

		if c.verbose && !c.quiet && (i+1)%10 == 0 {
			fmt.Printf("%s Processed %d/%d files\n", c.cyan("→"), i+1, len(inputPaths))
		}
	}

	return results
}

// processFilesParallel processes files in parallel
func (c *ExtractCommand) processFilesParallel(inputPaths []string, registry *PluginRegistry) []ExtractionResult {
	var wg sync.WaitGroup
	fileChan := make(chan string, len(inputPaths))
	resultChan := make(chan ExtractionResult, len(inputPaths))
	errorChan := make(chan error, len(inputPaths))

	var results []ExtractionResult

	// Start worker goroutines
	numWorkers := c.parallel
	if numWorkers > runtime.NumCPU() {
		numWorkers = runtime.NumCPU()
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for path := range fileChan {
				result, err := c.processSingleFile(path, registry)
				if err != nil {
					errorChan <- fmt.Errorf("%s: %v", path, err)
					continue
				}
				resultChan <- result
			}
		}(i)
	}

	// Send files to workers
	for _, path := range inputPaths {
		fileChan <- path
	}
	close(fileChan)

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Collect results
	for result := range resultChan {
		results = append(results, result)
	}

	// Report errors
	if !c.quiet {
		for err := range errorChan {
			fmt.Printf("%s %v\n", c.red("✗"), err)
		}
	}

	return results
}

// processSingleFile processes a single file
func (c *ExtractCommand) processSingleFile(path string, registry *PluginRegistry) (ExtractionResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return ExtractionResult{}, fmt.Errorf("failed to read file: %v", err)
	}

	// Determine the appropriate extractor based on file extension or language flag
	var extractor extractor.ExtractorPlugin
	
	if c.language != "" {
		// Use the specified language
		extractor = registry.GetExtractorByLanguage(c.language)
	} else {
		// Auto-detect based on file extension
		ext := strings.ToLower(filepath.Ext(path))
		extractor = registry.GetExtractorByExtension(ext)
	}

	if extractor == nil {
		// Use generic extractor as fallback
		extractor = registry.GetExtractorByLanguage("generic")
	}

	if extractor == nil {
		return ExtractionResult{}, fmt.Errorf("no suitable extractor found for file: %s", path)
	}

	// Initialize the extractor
	if err := extractor.Initialize(); err != nil {
		return ExtractionResult{}, fmt.Errorf("failed to initialize extractor: %v", err)
	}
	defer extractor.Cleanup()

	// Extract code blocks
	blocks := extractor.Extract(string(content))

	// Write extracted blocks to output directory
	var writtenFiles []string
	if !c.dryRun {
		writtenFiles = c.writeExtractedBlocks(blocks, path)
	}

	return ExtractionResult{
		SourceFile:    path,
		ExtractorName: extractor.Name(),
		CodeBlocks:    blocks,
		WrittenFiles:  writtenFiles,
	}, nil
}

// writeExtractedBlocks writes the extracted code blocks to the output directory
func (c *ExtractCommand) writeExtractedBlocks(blocks []extractor.CodeBlock, sourceFile string) []string {
	var writtenFiles []string

	for i, block := range blocks {
		// Create a unique filename based on the source file and block type
		baseName := strings.TrimSuffix(filepath.Base(sourceFile), filepath.Ext(sourceFile))
		var fileName string
		
		if block.Filename != "" {
			fileName = block.Filename
		} else {
			fileName = fmt.Sprintf("%s_%s_%d%s", baseName, block.Type, i, getExtensionForLanguage(block.Language))
		}

		// Create output path
		outputPath := filepath.Join(c.outputDir, fileName)

		// Ensure the directory exists
		dir := filepath.Dir(outputPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("%s Failed to create directory %s: %v\n", c.red("✗"), dir, err)
			continue
		}

		// Write the content
		if err := os.WriteFile(outputPath, []byte(block.Content), 0644); err != nil {
			fmt.Printf("%s Failed to write file %s: %v\n", c.red("✗"), outputPath, err)
			continue
		}

		writtenFiles = append(writtenFiles, outputPath)
		
		if c.verbose && !c.quiet {
			fmt.Printf("%s Wrote %s (%d bytes)\n", c.cyan("→"), fileName, len(block.Content))
		}
	}

	return writtenFiles
}

// getExtensionForLanguage returns the appropriate file extension for a language
func getExtensionForLanguage(lang string) string {
	switch strings.ToLower(lang) {
	case "go", "golang":
		return ".go"
	case "java":
		return ".java"
	case "python":
		return ".py"
	case "javascript", "js":
		return ".js"
	case "typescript", "ts":
		return ".ts"
	case "rust":
		return ".rs"
	case "dart":
		return ".dart"
	case "json":
		return ".json"
	case "xml":
		return ".xml"
	case "yaml", "yml":
		return ".yaml"
	case "markdown", "md":
		return ".md"
	case "text", "txt":
		return ".txt"
	default:
		return ".txt"
	}
}

// generateReport generates a detailed report of the extraction
func (c *ExtractCommand) generateReport(results []ExtractionResult) {
	fmt.Printf("\n%s Extraction Report\n", c.cyan("📊"))
	fmt.Printf("%s %s\n", c.cyan("┌"), strings.Repeat("─", 50))
	
	totalBlocks := 0
	totalFiles := 0
	for _, result := range results {
		totalBlocks += len(result.CodeBlocks)
		totalFiles += len(result.WrittenFiles)
	}

	fmt.Printf("%s Total source files processed: %s\n", c.cyan("│"), c.green(fmt.Sprintf("%d", len(results))))
	fmt.Printf("%s Total code blocks extracted: %s\n", c.cyan("│"), c.green(fmt.Sprintf("%d", totalBlocks)))
	fmt.Printf("%s Total files written: %s\n", c.cyan("│"), c.green(fmt.Sprintf("%d", totalFiles)))

	if c.verbose {
		fmt.Printf("%s Detailed breakdown:\n", c.cyan("│"))
		for _, result := range results {
			fmt.Printf("%s   %s: %d blocks -> %d files\n", 
				c.cyan("│"), 
				result.SourceFile, 
				len(result.CodeBlocks), 
				len(result.WrittenFiles))
		}
	}

	fmt.Printf("%s %s\n", c.cyan("└"), strings.Repeat("─", 50))
}

// ExtractionResult holds the result of extracting code blocks from a file
type ExtractionResult struct {
	SourceFile    string
	ExtractorName string
	CodeBlocks    []extractor.CodeBlock
	WrittenFiles  []string
}