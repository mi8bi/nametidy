#!/bin/bash

# Script to download and install the latest NameTidy binary

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
echo "NameTidy Installer"
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

# Construct the download URL
RELEASE_BASE_URL="https://github.com/mi8bi/NameTidy/releases/latest/download"
ASSET_NAME="NameTidy_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="${RELEASE_BASE_URL}/${ASSET_NAME}"

echo "Constructed download URL: $DOWNLOAD_URL"

# Define download path
ARCHIVE_PATH="${TMP_DIR}/${ASSET_NAME}"

# Download the archive
echo "Downloading NameTidy release asset..."
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

# Extracting only the 'nametidy' binary from the archive to the temp directory's root
if ! tar -xzf "$ARCHIVE_PATH" -C "$TMP_DIR" "$BINARY_NAME"; then
    echo "Error: Failed to extract '$BINARY_NAME' from '$ARCHIVE_PATH'."
    echo "The archive might be corrupted, the binary name/path inside the archive is unexpected,"
    echo "or the binary for your system ($OS/$ARCH) might not be correctly packaged as '$BINARY_NAME'."
    exit 1
fi
echo "Binary '$BINARY_NAME' extracted to '$EXTRACTED_BINARY_PATH'"

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
