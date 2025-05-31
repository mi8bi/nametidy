#Requires -Version 5.0
<#
.SYNOPSIS
    Installs the NameTidy utility for Windows.
.DESCRIPTION
    This script downloads the latest version of NameTidy for the appropriate
    architecture, extracts it, installs it to $env:USERPROFILE\bin, and
    attempts to add this directory to the user's PATH environment variable.
.NOTES
    Author: AI Assistant
    Version: 1.0
#>

# Stop on first error
$ErrorActionPreference = 'Stop'

# --- Configuration ---
$BinaryName = "nametidy.exe"
$InstallBaseDir = $env:USERPROFILE
$InstallRelativePath = "bin" # Relative to InstallBaseDir
$InstallDir = Join-Path -Path $InstallBaseDir -ChildPath $InstallRelativePath

$ReleaseBaseUrl = "https://github.com/mi8bi/NameTidy/releases/latest/download"
$OsName = "windows"
$Architecture = ""

# --- Temporary Path ---
# Create a temporary directory for downloads and extraction
$TempDir = Join-Path -Path $env:TEMP -ChildPath "NameTidy_Install_$($PID)_$(Get-Random)"
try {
    if (Test-Path $TempDir) {
        Write-Warning "Temporary directory $TempDir already exists. Removing."
        Remove-Item -Path $TempDir -Recurse -Force
    }
    New-Item -Path $TempDir -ItemType Directory -Force | Out-Null
    Write-Host "Temporary directory created: $TempDir"
}
catch {
    Write-Error "Failed to create temporary directory '$TempDir'. Error: $($_.Exception.Message)"
    exit 1
}

# --- Main Script Logic ---
try {
    Write-Host "NameTidy Installer for Windows PowerShell"
    Write-Host "---------------------------------------"

    # --- Determine Architecture ---
    Write-Host "Detecting system architecture..."
    switch ($env:PROCESSOR_ARCHITECTURE) {
        "AMD64" { $Architecture = "amd64" }
        "ARM64" { $Architecture = "arm64" }
        default {
            throw "Unsupported architecture: $($env:PROCESSOR_ARCHITECTURE). Only AMD64 and ARM64 are supported."
        }
    }
    Write-Host "Detected Architecture: $Architecture"

    # --- Construct Download URL ---
    $AssetName = "NameTidy_${OsName}_${Architecture}.zip"
    $DownloadUrl = "$ReleaseBaseUrl/$AssetName"
    $ArchiveFilePath = Join-Path -Path $TempDir -ChildPath $AssetName
    $ExtractedBinaryPath = Join-Path -Path $TempDir -ChildPath $BinaryName # Assuming it's at the root of the zip

    Write-Host "Download URL: $DownloadUrl"

    # --- Download ---
    Write-Host "Downloading NameTidy release asset..."
    try {
        Invoke-WebRequest -Uri $DownloadUrl -OutFile $ArchiveFilePath -UseBasicParsing
        Write-Host "Download successful: $ArchiveFilePath"
    }
    catch {
        throw "Download failed from $DownloadUrl. Error: $($_.Exception.Message). Please check the URL, your internet connection, or if an asset for '$OsName/$Architecture' is available."
    }

    # --- Extraction ---
    Write-Host "Extracting $BinaryName from $ArchiveFilePath..."
    try {
        # Expand-Archive extracts all files. We assume nametidy.exe is at the root of the zip.
        Expand-Archive -Path $ArchiveFilePath -DestinationPath $TempDir -Force
        if (-not (Test-Path $ExtractedBinaryPath)) {
            throw "$BinaryName not found in the extracted files at $TempDir. The archive might not contain '$BinaryName' at its root."
        }
        Write-Host "Extraction successful. $BinaryName is at $ExtractedBinaryPath"
    }
    catch {
        throw "Extraction failed. Error: $($_.Exception.Message). The archive might be corrupt or incompatible."
    }

    # --- Installation ---
    Write-Host "Installing $BinaryName..."

    # Create installation directory if it doesn't exist
    if (-not (Test-Path $InstallDir)) {
        Write-Host "Creating installation directory: $InstallDir"
        New-Item -Path $InstallDir -ItemType Directory -Force | Out-Null
    }

    $FinalInstallPath = Join-Path -Path $InstallDir -ChildPath $BinaryName

    # Check if binary already exists and prompt for overwrite
    if (Test-Path $FinalInstallPath) {
        Write-Warning "$BinaryName already exists at $FinalInstallPath."
        $overwriteChoice = Read-Host "Do you want to overwrite it? (Y/N)"
        if ($overwriteChoice -ne 'Y' -and $overwriteChoice -ne 'y') {
            Write-Host "Overwrite cancelled by user. Exiting."
            # No error, graceful exit. Cleanup will still run.
            return
        }
        Write-Host "Proceeding with overwrite..."
    }

    # Move the binary
    try {
        Move-Item -Path $ExtractedBinaryPath -Destination $FinalInstallPath -Force
        Write-Host "$BinaryName successfully installed to $FinalInstallPath."
    }
    catch {
        throw "Failed to move $BinaryName to $FinalInstallPath. Error: $($_.Exception.Message). Check permissions."
    }

    # --- PATH Update ---
    Write-Host "Attempting to add $InstallDir to your User PATH..."
    try {
        $currentUserPath = [System.Environment]::GetEnvironmentVariable('Path', 'User')
        $pathParts = $currentUserPath -split ';' | Where-Object { $_ -ne '' } # Remove empty entries

        if ($pathParts -notcontains $InstallDir) {
            $newPath = ($pathParts + $InstallDir) -join ';'
            [System.Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
            Write-Host "$InstallDir added to User PATH. Changes will apply to new PowerShell sessions/Windows logon."
            Write-Host "You might need to restart your current PowerShell session or log out and back in."
        } else {
            Write-Host "$InstallDir is already in the User PATH."
        }
    }
    catch {
        Write-Warning "Failed to automatically add $InstallDir to User PATH. Error: $($_.Exception.Message)"
        Write-Warning "You may need to add it manually through System Properties (Environment Variables)."
    }

    Write-Host "Installation complete!"
    Write-Host "Please open a new PowerShell session or terminal to use $BinaryName."

}
catch {
    # Catch any error that occurred in the main try block
    Write-Error "An error occurred during installation: $($_.Exception.Message)"
    Write-Error "Script execution aborted."
    # Exit with an error code if desired, e.g., exit 1, though script will terminate anyway.
}
finally {
    # --- Cleanup ---
    if (Test-Path $TempDir) {
        Write-Host "Cleaning up temporary directory: $TempDir"
        Remove-Item -Path $TempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
    Write-Host "Installer finished."
}
