#Requires -Version 5.0
<#
.SYNOPSIS
    Installs the nametidy utility for Windows.
.DESCRIPTION
    This script downloads the latest version of nametidy for the appropriate
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

$ReleaseBaseUrl = "https://github.com/mi8bi/nametidy/releases/latest/download"
$OsName = "windows"
$Architecture = ""

# --- Temporary Path ---
# Create a temporary directory for downloads and extraction
$TempDir = Join-Path -Path $env:TEMP -ChildPath "nametidy_Install_$($PID)_$(Get-Random)"
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
    Write-Host "nametidy Installer for Windows PowerShell"
    Write-Host "---------------------------------------"

    # --- Fetch Latest Release Information ---
    Write-Host "Fetching latest release information..."
    $LatestReleaseInfoUrl = "https://api.github.com/repos/mi8bi/nametidy/releases/latest"
    $TagName = ""
    $Version = ""

    try {
        # Using -UseBasicParsing for compatibility, though for modern PowerShell it's often not needed.
        $ReleaseInfo = Invoke-RestMethod -Uri $LatestReleaseInfoUrl -UseBasicParsing
        $TagName = $ReleaseInfo.tag_name
        if (-not $TagName) {
            throw "tag_name not found in the release information from $LatestReleaseInfoUrl."
        }
        Write-Host "Latest release tag: $TagName"

        # Remove 'v' prefix if it exists
        if ($TagName.StartsWith("v")) {
            $Version = $TagName.Substring(1)
        } else {
            # If no 'v' prefix, use the tag name as version directly.
            # This might be relevant if tag naming changes, though GitHub convention is often vX.Y.Z
            $Version = $TagName
        }

        if (-not $Version) {
            throw "Could not extract version from tag '$TagName'."
        }
        Write-Host "Detected version: $Version"
    }
    catch {
        # Specific error for this block, then re-throw to be caught by the main script's try-catch
        Write-Error "Error fetching or parsing release information: $($_.Exception.Message)"
        throw $_
    }

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
    # $OsName is already defined in Configuration section
    # $Architecture is determined above
    # $Version and $TagName are from the new block above
    $AssetName = "nametidy_${Version}_${OsName}_${Architecture}.zip"
    # Construct specific download URL using the tag
    $DownloadUrl = "https://github.com/mi8bi/nametidy/releases/download/${TagName}/${AssetName}"
    # Example: https://github.com/mi8bi/nametidy/releases/download/v0.1.0/nametidy_0.1.0_windows_amd64.zip
    $ArchiveFilePath = Join-Path -Path $TempDir -ChildPath $AssetName
    $ExtractedBinaryPath = Join-Path -Path $TempDir -ChildPath $BinaryName # Assuming it's at the root of the zip

    Write-Host "Download URL: $DownloadUrl"

    # --- Download ---
    Write-Host "Downloading nametidy release asset..."
    try {
        Invoke-WebRequest -Uri $DownloadUrl -OutFile $ArchiveFilePath -UseBasicParsing
        Write-Host "Download successful: $ArchiveFilePath"
    }
    catch {
        throw "Download failed from $DownloadUrl. Error: $($_.Exception.Message). Please check the URL, your internet connection, or if an asset for '$OsName/$Architecture' is available."
    }

    # --- Extraction ---
    Write-Host "Extracting archive from $ArchiveFilePath to $TempDir..."
    try {
        Expand-Archive -Path $ArchiveFilePath -DestinationPath $TempDir -Force
        Write-Host "Archive extraction completed."

        # Search for the binary, trying possible names
        # $BinaryName is "nametidy.exe" (defined in Configuration)
        $PossibleNames = @($BinaryName, "NameTidy.exe")
        $FoundBinaryInfo = $null

        foreach ($name_to_find in $PossibleNames) {
            Write-Host "Searching for binary '$name_to_find' in '$TempDir'..."
            # Get-ChildItem -File ensures we only get files, -Recurse searches subdirectories.
            # Using -Filter for efficiency if supported, otherwise -Include. For simple names, -Filter is fine.
            $foundFiles = Get-ChildItem -Path $TempDir -Recurse -File -Filter $name_to_find -ErrorAction SilentlyContinue

            if ($foundFiles) {
                # Take the first one if multiple are somehow found (e.g. in different subdirs or with same name)
                $FoundBinaryInfo = $foundFiles | Select-Object -First 1
                if ($FoundBinaryInfo) {
                    # Update $ExtractedBinaryPath which was initially set to $TempDir\$BinaryName
                    $ExtractedBinaryPath = $FoundBinaryInfo.FullName
                    Write-Host "Found executable binary at: $ExtractedBinaryPath"
                    break # Exit foreach loop as we found our binary
                }
            }
        }

        if (-not $FoundBinaryInfo) {
            Write-Host "Listing contents of '$TempDir' (top level):"
            Get-ChildItem -Path $TempDir -Depth 0 | ForEach-Object { Write-Host "  $($_.Name)" } # Depth 0 for top level
            Write-Host "Listing contents of '$TempDir' (recursive, files only, relative paths):"
            Get-ChildItem -Path $TempDir -Recurse -File | ForEach-Object { Write-Host "  $($_.FullName.Substring($TempDir.Length).TrimStart('\'))" }
            throw "Could not find '$($PossibleNames -join "' or '")' in the extracted files at '$TempDir'."
        }
        # $ExtractedBinaryPath is now updated to the actual found binary path.
        # The script will proceed to use this updated $ExtractedBinaryPath for installation.
        # If "NameTidy.exe" was found, $ExtractedBinaryPath points to it.
        # The $FinalInstallPath is $InstallDir\$BinaryName ("nametidy.exe").
        # So, Move-Item will effectively rename "NameTidy.exe" to "nametidy.exe" if that was what was found. This is desired.
    }
    catch {
        # This catch block handles errors from Expand-Archive or the 'throw' if binary not found.
        throw "Extraction or binary search failed. Error: $($_.Exception.Message). The archive might be corrupt, incompatible, or the binary is missing/not found."
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
