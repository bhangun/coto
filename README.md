# Coto CLI

A powerful CLI tool to recursively combine file contents from directories and subdirectories with an enhanced interactive mode and comprehensive filtering options.

## 🚀 Features

- **Interactive Mode**: When run without arguments, Coto enters an interactive mode prompting for all options
- **Recursive File Combination**: Combines files from directories and subdirectories
- **File Renaming**: Rename files based on patterns, prefixes, suffixes, or regular expressions
- **Multiple Output Formats**: Text, JSON, XML, and Markdown
- **Flexible Filtering**: Filter by file extensions, size, patterns, and more
- **Parallel Processing**: Process multiple files simultaneously for faster performance
- **Compression Support**: Optional GZIP compression for output
- **Configuration Files**: Load settings from JSON configuration files
- **Progress Indicators**: Real-time progress for large operations
- **Cross-Platform**: Works on Linux, macOS, and Windows

## 📦 Installation

### Automated Installation (Recommended)

The easiest way to install Coto is using the automated installation scripts:

**macOS/Linux (using curl):**
```bash
curl -sSL https://raw.githubusercontent.com/bhangun/coto/main/install.sh | bash
```

**Windows (using PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/bhangun/coto/main/install.ps1 | iex
```

These scripts will automatically:
- Detect your operating system and architecture
- Download the appropriate binary from the latest GitHub release
- Verify the checksum for security
- Install the binary to the correct location
- Add it to your PATH if needed

### Package Managers

For users who prefer package managers:

**Homebrew (macOS/Linux):**
```bash
brew install bhangun/coto/coto
```

**Chocolatey (Windows):**
```powershell
choco install coto
```

### From Source
```bash
git clone https://github.com/bhangun/coto.git
cd coto
make build
# Binary will be available at bin/coto
./bin/coto --help
```

### Manual Installation
```bash
# Build from source
make build
# The binary will be available at bin/coto
# Make it executable and move to your PATH
chmod +x bin/coto
sudo mv bin/coto /usr/local/bin/
```

## 🛠 Usage

### Main Command (File Combination)
Simply run Coto without any arguments to enter interactive mode:
```bash
./coto
```

#### Command Line Mode
```bash
# Basic usage
coto -i ./src -o combined.txt

# Filter by file extensions
coto -ext .go,.js,.py -o output.txt

# Multiple filters
coto -i ./src -ext .go,.md --min-size 100 --max-size 1000000

# JSON output with compression
coto --format json --compress --output output.json.gz

# Parallel processing
coto --parallel 4 --verbose

# Exclude patterns
coto --exclude "\.git|node_modules|\.DS_Store"

# Configuration file
coto --config config.json

# Dry run to see what would be processed
coto --dry-run --verbose
```

### Rename Command (New!)
The rename command allows you to rename files in a directory based on patterns, prefixes, suffixes, or regular expressions:

```bash
# Remove a prefix from filenames
coto rename -dir ./videos -prefix "Lk21.De-"

# Remove a suffix from filenames
coto rename -dir ./files -suffix "_backup"

# Remove a pattern anywhere in the filename
coto rename -dir ./files -pattern "_old_"

# Use regular expressions for complex renaming
coto rename -dir ./data -regex "^(\d+)_(.+)$" -replacement "$2"  # Remove leading numbers and underscore

# Dry run to preview changes without actually renaming
coto rename -dir ./photos -suffix ".bak" --dry-run

# Combine multiple operations
coto rename -dir ./files -prefix "old_" -suffix "_backup" -pattern "temp"
```

#### Rename Command Options

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--dir` | | Directory to rename files in (default: current directory) |
| `--prefix` | | Prefix to remove from filenames |
| `--suffix` | | Suffix to remove from filenames |
| `--pattern` | | Pattern to remove from filenames (anywhere in the name) |
| `--regex` | | Regular expression pattern to match |
| `--replacement` | | Replacement string for regex (use with --regex) |
| `--dry-run` | | Show what would be renamed without actually renaming |
| `--verbose` | | Show detailed progress |
| `--quiet` | | Suppress non-essential output |
| `--force` | | Force rename even if target file exists |
| `--help` | `-h` | Show help message |

### Extract Command
Extract code blocks from files:

```bash
# Extract from specific files
coto extract -input file1.txt,file2.txt -language javascript

# Extract from multiple files using glob patterns
coto extract -input "*.txt" -language python -output extracted/

# Parallel processing
coto extract -input "*.js" -parallel 4 -verbose

# Generate detailed report
coto extract -input code.txt -report
```

### Available Main Command Options

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--input` | `-i` | Input directory path (default: current directory) |
| `--output` | `-o` | Output file path (default: combined.txt) |
| `--ext` | | Comma-separated list of file extensions to include |
| `--exclude-hidden` | `-eh` | Exclude hidden files and directories (default: true) |
| `--max-size` | | Maximum file size in bytes (0 = unlimited) |
| `--min-size` | | Minimum file size in bytes |
| `--exclude` | | Regex pattern to exclude files |
| `--include` | | Regex pattern to include files |
| `--format` | | Output format: text, json, xml, markdown (default: text) |
| `--compress` | | Compress output with gzip |
| `--parallel` | | Number of files to process in parallel (default: 1) |
| `--dry-run` | | Show what would be processed without writing |
| `--quiet` | | Suppress non-essential output |
| `--verbose` | | Show detailed progress |
| `--config` | | Load configuration from JSON file |
| `--version` | `-v` | Show version information |
| `--help` | `-h` | Show help message |

## 📁 Sample Configuration File (config.json)

```json
{
  "input_dir": "./src",
  "output_file": "combined.txt",
  "extensions": [".go", ".js", ".py"],
  "exclude_hidden": true,
  "max_file_size": 1000000,
  "min_file_size": 0,
  "exclude_pattern": "\\.git|node_modules",
  "include_pattern": "",
  "output_format": "text",
  "compress": false,
  "parallel": 4,
  "quiet": false,
  "verbose": true,
  "dry_run": false
}
```

## 🚀 Deployment

Coto uses [JReleaser](https://jreleaser.org/) for automated releases and distribution to package managers:

- **Homebrew**: Releases are automatically published to the custom tap at `bhangun/coto`
- **Chocolatey**: Windows packages are automatically published to the Chocolatey Community Repository
- **GitHub Releases**: Binaries for all platforms are attached to each release with checksums
- **Automatic Signing**: All binaries are signed for security

### Homebrew Installation
After publication, users can install via Homebrew:
```bash
brew tap bhangun/coto
brew install coto
```

Or in a single command:
```bash
brew install bhangun/coto/coto
```

### Chocolatey Installation
On Windows, users can install via Chocolatey:
```powershell
choco install coto
```

The deployment workflow is triggered on tagged commits and handles:
1. Cross-compilation for all supported platforms (Linux AMD64/ARM64, macOS AMD64/ARM64, Windows AMD64/ARM64)
2. Binary signing and checksum generation
3. Publication to Homebrew and Chocolatey package managers
4. GitHub release creation with assets
5. Automatic formula/package updates

## 🏗 Development

```bash
# Clone repository
git clone https://github.com/bhangun/coto.git
cd coto

# Build
make build

# Run tests
make test

# Build for all platforms
make cross-compile

# Run demo
make demo
```

## 🎯 Use Cases

- **AI Context Gathering**: Combine source code files for LLM context
- **Documentation Aggregation**: Merge documentation files into a single document
- **Code Review Preparation**: Bundle related files for review
- **File Organization**: Rename files using patterns, prefixes, or suffixes (with `coto rename`)
- **Media File Management**: Clean up media filenames by removing unwanted prefixes/suffixes
- **Backup Operations**: Consolidate files from multiple directories
- **Data Analysis**: Combine structured data files for processing

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

MIT

## 🆘 Support

For support, please open an issue in the GitHub repository.