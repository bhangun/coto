# Pecel Windows Installation Script
param(
    [string]$Action = "install",
    [string]$Version = "v0.1.0"
)

$ErrorActionPreference = "Stop"

# Configuration
$Repo = "bhangun/pecel"
$BinaryName = "pecel.exe"
$InstallDir = "$env:USERPROFILE\bin"
$TempDir = "$env:TEMP\pecel-install"

# Colors
$Green = "`e[32m"
$Yellow = "`e[33m"
$Red = "`e[31m"
$Reset = "`e[0m"

function Write-Info {
    Write-Host "${Green}[INFO]${Reset} $($args[0])"
}

function Write-Warn {
    Write-Host "${Yellow}[WARN]${Reset} $($args[0])"
}

function Write-Error {
    Write-Host "${Red}[ERROR]${Reset} $($args[0])"
    exit 1
}

# Determine platform
$Arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    "AMD64" { "amd64" }
    "ARM64" { "arm64" }
    default { "amd64" }
}

$Platform = "windows"

function Install-Binary {
    Write-Info "Attempting to download pecel $Version for $Platform/$Arch..."

    $DownloadUrl = "https://github.com/$Repo/releases/download/$Version/pecel-$Platform-$Arch.exe"

    # Create temp directory
    if (Test-Path $TempDir) {
        Remove-Item -Path $TempDir -Recurse -Force
    }
    New-Item -ItemType Directory -Path $TempDir -Force | Out-Null

    # Try to download binary
    try {
        Invoke-WebRequest -Uri $DownloadUrl -OutFile "$TempDir\$BinaryName"
        Write-Info "Successfully downloaded binary"
    } catch {
        Write-Warn "Failed to download binary: $_"
        Write-Info "Attempting to build from source..."

        # Check if Go is installed
        if (!(Get-Command "go" -ErrorAction SilentlyContinue)) {
            Write-Error "Go is required to build from source but is not installed.`nPlease install Go first or check the release at: https://github.com/$Repo/releases"
        }

        # Clone repo and build
        $RepoDir = "$TempDir\repo"
        git clone "https://github.com/$Repo.git" $RepoDir
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Failed to clone repository"
        }

        Set-Location $RepoDir
        .\make.bat build 2>$null
        if ($LASTEXITCODE -ne 0) {
            # If make.bat doesn't exist, try manual build
            go build -o bin/pecel.exe ./cmd/main
            if ($LASTEXITCODE -ne 0) {
                Write-Error "Build failed"
            }
        }

        Copy-Item "bin\pecel.exe" "$TempDir\$BinaryName" -Force
        Write-Info "Successfully built from source"
    }

    # Create install directory if it doesn't exist
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    # Install binary
    Copy-Item "$TempDir\$BinaryName" "$InstallDir\$BinaryName" -Force

    # Add to PATH if not already present
    $CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($CurrentPath -notlike "*$InstallDir*") {
        Write-Info "Adding $InstallDir to user PATH..."
        [Environment]::SetEnvironmentVariable("Path", "$CurrentPath;$InstallDir", "User")
        $env:Path += ";$InstallDir"
    }

    # Cleanup
    Remove-Item -Path $TempDir -Recurse -Force

    Write-Info "Installation completed!"
    Write-Info "Run 'pecel --help' to get started"
}

function Uninstall {
    $BinaryPath = "$InstallDir\$BinaryName"
    if (Test-Path $BinaryPath) {
        Write-Info "Removing $BinaryPath..."
        Remove-Item -Path $BinaryPath -Force
        
        # Remove from PATH
        $CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
        $NewPath = $CurrentPath -replace [regex]::Escape($InstallDir), "" -replace ";;", ";"
        [Environment]::SetEnvironmentVariable("Path", $NewPath.TrimEnd(';'), "User")
        
        Write-Info "Pecel has been uninstalled"
    } else {
        Write-Warn "Pecel is not installed"
    }
}

function Update {
    Write-Info "Checking for updates..."
    
    try {
        $LatestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        $LatestVersion = $LatestRelease.tag_name
        
        if ($LatestVersion -ne $Version) {
            Write-Info "New version available: $LatestVersion"
            $Version = $LatestVersion
            Install-Binary
        } else {
            Write-Info "You have the latest version ($Version)"
        }
    } catch {
        Write-Error "Failed to check for updates: $_"
    }
}

# Main execution
switch ($Action.ToLower()) {
    "install" {
        Install-Binary
    }
    "update" {
        Update
    }
    "uninstall" {
        Uninstall
    }
    default {
        Write-Error "Unknown action: $Action"
        Write-Host "Usage: .\install.ps1 [install|update|uninstall]"
        exit 1
    }
}