#!/bin/bash

# Script to download and install the latest nametidy binary

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Global Variables ---
TMP_DIR="" # Initialize TMP_DIR, will be set in main
BINARY_NAME="nametidy"
INSTALL_PATH="/usr/local/bin"
INSTALLED_BINARY_PATH="${INSTALL_PATH}/${BINARY_NAME}"
EXTRACTED_BINARY_PATH="" # Will be set after extraction

# --- Helper Functions ---
# Function to clean up temporary directory on exit
cleanup() {
    if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
        # Check if the binary still exists in TMP_DIR (i.e., wasn't moved)
        # This check is mainly for interactive scenarios; if script exits due to error, it's cleaned.
        if [ -f "$EXTRACTED_BINARY_PATH" ]; then
            echo "Note: The extracted binary '$EXTRACTED_BINARY_PATH' was not moved to $INSTALL_PATH."
            echo "Cleaning up temporary directory, including the binary."
        else
            echo "Cleaning up temporary directory: $TMP_DIR"
        fi
        rm -rf "$TMP_DIR"
    fi
}

# Trap EXIT, ERR, and INT signals to call cleanup function
# ERR is for set -e
# INT is for Ctrl+C
trap cleanup EXIT ERR INT

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to download using curl or wget
download_file() {
    local url="$1"
    local output_path="$2"

    if command_exists curl; then
        echo "Attempting to download with curl..."
        # -L: follow redirects, -S: show errors, -s: silent, -f: fail fast (HTTP errors => exit code 22)
        if curl -LSsf -o "$output_path" "$url"; then
            echo "Download successful: $output_path"
            return 0
        else
            echo "curl download failed. HTTP error or other issue."
            return 1 # curl with -f already returns non-zero on server errors
        fi
    elif command_exists wget; then
        echo "Attempting to download with wget..."
        # -q: quiet, -O: output file
        if wget -q -O "$output_path" "$url"; then
            echo "Download successful: $output_path"
            return 0
        else
            echo "wget download failed."
            return 1
        fi
    else
        echo "Error: Neither curl nor wget found. Please install one of them and try again."
        exit 1 # This will trigger ERR trap and then EXIT trap
    fi
}

# --- Main Script ---
echo "nametidy Installer"
echo "------------------"

# Create a temporary directory for download
# mktemp -d will create a directory with a unique name
TMP_DIR=$(mktemp -d 2>/dev/null || mktemp -d -t 'nametidy-install.XXXXXXXXXX')
if [ ! -d "$TMP_DIR" ]; then
    echo "Error: Failed to create temporary directory."
    exit 1
fi
# Set EXTRACTED_BINARY_PATH now that TMP_DIR is known
EXTRACTED_BINARY_PATH="${TMP_DIR}/${BINARY_NAME}"
echo "Temporary directory created: $TMP_DIR"


# Determine OS and Architecture
OS=""
ARCH=""

echo "Detecting OS and architecture..."
case "$(uname -s)" in
    Linux*)     OS="linux";;
    Darwin*)    OS="darwin";;
    *)          echo "Error: Unsupported OS: $(uname -s)"; exit 1;;
esac

case "$(uname -m)" in
    x86_64)     ARCH="amd64";;
    arm64)      ARCH="arm64";;
    aarch64)    ARCH="arm64";;
    *)          echo "Error: Unsupported architecture: $(uname -m)"; exit 1;;
esac

echo "Detected OS: $OS"
echo "Detected Architecture: $ARCH"

# Fetch the latest release tag and version
echo "Fetching latest release information..."
LATEST_RELEASE_INFO_URL="https://api.github.com/repos/mi8bi/nametidy/releases/latest"

# Attempt to fetch tag_name using curl and grep/sed
# This extracts the value of "tag_name": "vX.Y.Z"
TAG_NAME=$(curl -LsS "$LATEST_RELEASE_INFO_URL" | grep '"tag_name":' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')

if [ -z "$TAG_NAME" ]; then
    echo "Error: Could not fetch the latest release tag_name from $LATEST_RELEASE_INFO_URL."
    echo "Please check your internet connection and if the repository/releases are accessible."
    # Try to see if jq is available for a more robust parsing as a fallback
    if command_exists jq; then
        echo "Attempting with jq..."
        TAG_NAME=$(curl -LsS "$LATEST_RELEASE_INFO_URL" | jq -r .tag_name)
         if [ -z "$TAG_NAME" ] || [ "$TAG_NAME" == "null" ]; then
            echo "Error: jq also failed to fetch or parse tag_name."
            exit 1
         fi
    else
        echo "jq not available. Could not parse release information automatically."
        exit 1
    fi
fi

echo "Latest release tag: $TAG_NAME"

# Remove 'v' prefix from tag to get version
VERSION=$(echo "$TAG_NAME" | sed 's/^v//')
if [ -z "$VERSION" ]; then
    echo "Error: Could not extract version from tag '$TAG_NAME'."
    exit 1
fi
echo "Detected version: $VERSION"

# Construct the download URL using the fetched tag and version
# Note: For Windows, the script might use .zip. This script is .sh, so .tar.gz is appropriate.
# The ASSET_NAME now includes the VERSION.
RELEASE_DOWNLOAD_URL_BASE="https://github.com/mi8bi/nametidy/releases/download"
ASSET_NAME="nametidy_${VERSION}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="${RELEASE_DOWNLOAD_URL_BASE}/${TAG_NAME}/${ASSET_NAME}"
# Example: https://github.com/mi8bi/nametidy/releases/download/v0.1.0/nametidy_0.1.0_linux_amd64.tar.gz

echo "Constructed download URL: $DOWNLOAD_URL"

# Define download path
ARCHIVE_PATH="${TMP_DIR}/${ASSET_NAME}"

# Download the archive
echo "Downloading nametidy release asset..."
if ! download_file "$DOWNLOAD_URL" "$ARCHIVE_PATH"; then
    echo "Error: Download failed from $DOWNLOAD_URL."
    echo "Please check the URL and your internet connection."
    echo "It's possible that a release asset for your OS/architecture ($OS/$ARCH) is not available."
    exit 1
fi

# Extract the binary
echo "Extracting the binary..."
if ! command_exists tar; then
    echo "Error: 'tar' command not found. Please install tar and try again."
    exit 1
fi

echo "Extracting archive contents..."
if ! tar -xzf "$ARCHIVE_PATH" -C "$TMP_DIR"; then
    echo "Error: Failed to extract archive '$ARCHIVE_PATH' to '$TMP_DIR'."
    echo "The archive might be corrupted."
    exit 1
fi
echo "Archive extracted to '$TMP_DIR'."

# Search for the binary, first 'nametidy', then 'nametidy'
# We look for an executable file.
FOUND_BINARY_PATH=""
# $BINARY_NAME is "nametidy" as defined in Global Variables
POSSIBLE_NAMES=("$BINARY_NAME" "nametidy")

for name_to_find in "${POSSIBLE_NAMES[@]}"; do
    echo "Searching for executable binary '$name_to_find' in '$TMP_DIR'..."
    # Use find to locate the executable file. -print -quit ensures we only get one if multiple exist.
    # Redirect stderr to /dev/null for -executable if user doesn't own files, though less likely in TMP_DIR.
    # -path '*/.DS_Store' -prune -o handles macOS specific files if any present in archive
    found_path=$(find "$TMP_DIR" -path '*/.DS_Store' -prune -o -type f -name "$name_to_find" -print -quit 2>/dev/null)

    if [ -n "$found_path" ]; then
        # Found a file with the name, now check if it's executable or can be made executable
        if [ -x "$found_path" ]; then
             FOUND_BINARY_PATH="$found_path"
             echo "Found executable binary at: $FOUND_BINARY_PATH"
             break
        else
            echo "Found file '$found_path' but it is not executable. Attempting to make it executable..."
            chmod +x "$found_path"
            if [ -x "$found_path" ]; then
                echo "Made '$found_path' executable."
                FOUND_BINARY_PATH="$found_path"
                echo "Found executable binary at: $FOUND_BINARY_PATH"
                break
            else
                # If chmod failed, this path is not viable. Continue search.
                echo "Could not make '$found_path' executable. Searching further..."
            fi
        fi
    fi
done

if [ -z "$FOUND_BINARY_PATH" ]; then
    echo "Error: Could not find the '$BINARY_NAME' or 'nametidy' executable binary in the extracted files at '$TMP_DIR'."
    echo "Archive contents (top level):"
    ls -Al "$TMP_DIR"
    # More detailed listing if needed:
    # echo "Archive contents (full):"
    # find "$TMP_DIR" -ls
    exit 1
fi

# Update EXTRACTED_BINARY_PATH to the actual found path
EXTRACTED_BINARY_PATH="$FOUND_BINARY_PATH"
# The rest of the script uses $EXTRACTED_BINARY_PATH for chmod, mv, etc.
# Note: The global BINARY_NAME variable is still "nametidy".
# If "nametidy" was found, INSTALLED_BINARY_PATH (".../nametidy") would be different from the found name.
# The script later does `mv "$EXTRACTED_BINARY_PATH" "$INSTALLED_BINARY_PATH"`.
# This means if "nametidy" is found, it will be moved and renamed to "nametidy" in the install directory. This is the desired behavior.
echo "Using binary found at '$EXTRACTED_BINARY_PATH'"

# Clean up downloaded archive
rm "$ARCHIVE_PATH"

# Check if binary exists and prompt for overwrite
if [ -f "$INSTALLED_BINARY_PATH" ]; then
    echo "Warning: '$BINARY_NAME' already exists at $INSTALLED_BINARY_PATH."
    # Disable set -e for this interactive block
    set +e
    read -r -p "Do you want to overwrite it? (y/N): " overwrite_choice
    set -e
    case "$overwrite_choice" in
        y|Y ) echo "Proceeding with overwrite...";;
        * )   echo "Overwrite cancelled by user. Exiting."; exit 0;; # Graceful exit, trap will clean TMP_DIR
    esac
fi

# Make the binary executable
echo "Making '$BINARY_NAME' executable..."
if ! chmod +x "$EXTRACTED_BINARY_PATH"; then
    echo "Error: Failed to make binary executable. Please check permissions for '$EXTRACTED_BINARY_PATH'."
    exit 1
fi

# Move the binary to the installation path
echo "Attempting to install '$BINARY_NAME' to $INSTALL_PATH..."

# Create install directory if it doesn't exist (relevant for /usr/local/bin but good practice)
# This usually requires sudo as well.
if [ ! -d "$INSTALL_PATH" ]; then
    echo "Installation directory $INSTALL_PATH does not exist. Attempting to create it..."
    if command_exists sudo; then
        if ! sudo mkdir -p "$INSTALL_PATH"; then
            echo "Error: Failed to create $INSTALL_PATH even with sudo."
            echo "Please create $INSTALL_PATH manually and try again."
            exit 1
        fi
        echo "$INSTALL_PATH created."
    else
        echo "Error: $INSTALL_PATH does not exist and sudo is not available to create it."
        echo "Please create $INSTALL_PATH manually and try again."
        exit 1
    fi
fi


# Try direct move first
if mv "$EXTRACTED_BINARY_PATH" "$INSTALLED_BINARY_PATH" 2>/dev/null; then
    echo "Successfully installed '$BINARY_NAME' to $INSTALLED_BINARY_PATH."
else
    echo "Direct move failed (likely due to permissions)."
    if command_exists sudo; then
        echo "Attempting to move with sudo..."
        if sudo mv "$EXTRACTED_BINARY_PATH" "$INSTALLED_BINARY_PATH"; then
            echo "Successfully installed '$BINARY_NAME' to $INSTALLED_BINARY_PATH with sudo."
        else
            echo "Error: 'sudo mv' failed. You might not have sudo privileges or entered the wrong password."
            echo "The extracted binary is available at: $EXTRACTED_BINARY_PATH"
            echo "Please move it manually to a directory in your PATH."
            echo "Example: sudo mv '$EXTRACTED_BINARY_PATH' '$INSTALLED_BINARY_PATH'"
            set +e # Disable exit on error to allow prompt
            read -r -p "Press Enter to finish the script (this will remove the temporary file), or Ctrl+C to abort and move it manually now."
            set -e
            exit 1 # Indicate an issue occurred, trap will clean up.
        fi
    else
        echo "Error: 'sudo' command not found."
        echo "The extracted binary is available at: $EXTRACTED_BINARY_PATH"
        echo "Please move it manually to a directory in your PATH."
        echo "Example: mv '$EXTRACTED_BINARY_PATH' '$INSTALLED_BINARY_PATH' (you might need to be root or use sudo if available another way)"
        set +e
        read -r -p "Press Enter to finish the script (this will remove the temporary file), or Ctrl+C to abort and move it manually now."
        set -e
        exit 1 # Indicate an issue occurred, trap will clean up.
    fi
fi

# Check if installation path is in PATH
if [[ ":$PATH:" != *":${INSTALL_PATH}:"* ]]; then
    echo ""
    echo "Warning: Installation directory '$INSTALL_PATH' is not in your PATH."
    echo "You may need to add it to your shell configuration file (e.g., ~/.bashrc, ~/.zshrc, or /etc/profile):"
    echo "  export PATH=\"\$PATH:$INSTALL_PATH\""
    echo "Then, source the file (e.g., 'source ~/.bashrc') or open a new terminal."
fi

echo ""
echo "Installation of '$BINARY_NAME' complete!"
echo "You can now try running '$BINARY_NAME' from your terminal."

# Cleanup is handled by the trap on EXIT, ERR, INT.
# Setting EXTRACTED_BINARY_PATH to empty if it was successfully moved, so cleanup doesn't show confusing message.
EXTRACTED_BINARY_PATH=""
exit 0
