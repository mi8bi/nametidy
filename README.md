> ✨ A fast and simple CLI tool to clean and rename file names — powered by Go.

# [nametidy](https://mi8bi.github.io/nametidy/)

nametidy is a fast and flexible command-line tool for cleaning and renaming file names.
It supports operations such as filename cleanup, adding sequence numbers, and undoing changes — all with a simple and intuitive interface.

[![Build Status](https://github.com/mi8bi/nametidy/actions/workflows/test.yml/badge.svg)](https://github.com/mi8bi/nametidy/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mi8bi/nametidy)](https://goreportcard.com/report/github.com/mi8bi/nametidy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest Release](https://img.shields.io/github/v/release/mi8bi/nametidy)](https://github.com/mi8bi/nametidy/releases/latest)
[![GitHub Stars](https://img.shields.io/github/stars/mi8bi/nametidy?style=social)](https://github.com/mi8bi/nametidy/stargazers)

---

## 📽️ Demo

See nametidy in action:

[![nametidy_Demo002](https://asciinema.org/a/719898.svg)](https://asciinema.org/a/719898)

---

## Table of Contents

- [Automated Installation](#automated-installation)
  - [For Linux/macOS (using `install.sh`)](#for-linuxmacos-using-installsh)
  - [For Windows Command Prompt (using `install.cmd`)](#for-windows-command-prompt-using-installcmd)
  - [For Windows PowerShell (using `install.ps1`)](#for-windows-powershell-using-installps1)
- [Manual Installation](#manual-installation)
- [Build](#build)
- [Usage](#usage)
  - [Clean Up Filenames](#clean-up-filenames)
  - [Undo Changes](#undo-changes)
  - [Dry Run Mode](#dry-run-mode)
  - [Verbose Logging](#verbose-logging)
  - [Add Sequence Numbers](#add-sequence-numbers)
- [Options](#options)
- [License](#license)

---

## Automated Installation

You can use the following scripts to automate the installation of nametidy. These scripts will detect your system's architecture and download the appropriate binary.

### For Linux/macOS (using `install.sh`)

1.  **Download the script (it will be saved as `install.sh` in your current directory):**
    ```bash
    # Using curl:
    curl -LO https://raw.githubusercontent.com/mi8bi/nametidy/main/scripts/install.sh
    # Or using wget:
    # wget https://raw.githubusercontent.com/mi8bi/nametidy/main/scripts/install.sh

2.  **Make it executable:**
    ```bash
    chmod +x install.sh
    ```

3.  **Run the installer:**
    ```bash
    ./install.sh
    ```
    This script installs `nametidy` to `/usr/local/bin`. It may require `sudo` privileges if your user doesn't have write access to this directory.

### For Windows Command Prompt (using `install.cmd`)

1.  **Download the script:**
    You can download `install.cmd` directly from the repository (e.g., save it to your `Downloads` folder):
    [https://raw.githubusercontent.com/mi8bi/nametidy/main/scripts/install.cmd](https://raw.githubusercontent.com/mi8bi/nametidy/main/scripts/install.cmd)
    (Right-click the link and select "Save link as..." or "Save As...")

2.  **Run the installer:**
    Open Command Prompt (`cmd.exe`). Navigate to the directory where you saved `install.cmd` (e.g., `Downloads`) and run it:
    ```cmd
    cd C:\Users\YourUser\Downloads
    install.cmd
    ```
    Or, if `install.cmd` is in your current directory:
    ```cmd
    install.cmd
    ```
    This script installs `nametidy.exe` to `%USERPROFILE%\bin` and attempts to add this directory to your User PATH environment variable. PATH changes will apply to new Command Prompt sessions.

### For Windows PowerShell (using `install.ps1`)

1.  **Download the script:**
    You can download `install.ps1` directly from the repository (e.g., save it to your `Downloads` folder):
    [https://raw.githubusercontent.com/mi8bi/nametidy/main/scripts/install.ps1](https://raw.githubusercontent.com/mi8bi/nametidy/main/scripts/install.ps1)
    (Right-click the link and select "Save link as..." or "Save As...")

2.  **Run the installer:**
    Open PowerShell. Navigate to the directory where you saved `install.ps1` (e.g., `Downloads`). You may need to adjust your execution policy to run the script.
    Example:
    ```powershell
    cd C:\Users\YourUser\Downloads
    # Then run one of the following:
    ```
    To run the script for the current session without changing global policy:
    ```powershell
    PowerShell -ExecutionPolicy Bypass -File .\install.ps1
    ```
    Alternatively, from within that PowerShell prompt:
    ```powershell
    # Temporarily bypass execution policy for the current process
    Set-ExecutionPolicy Bypass -Scope Process -Force
    .\install.ps1
    ```
    This script installs `nametidy.exe` to `$env:USERPROFILE\bin` and attempts to add this directory to your User PATH environment variable persistently. PATH changes will apply to new PowerShell sessions or after restarting Windows.

---

## Manual Installation

You can download prebuilt binaries from the [GitHub Releases page](https://github.com/mi8bi/nametidy/releases):

1. Go to the [Releases](https://github.com/mi8bi/nametidy/releases) page on GitHub.
2. Find the latest release and download the binary file for your OS and architecture (e.g., `nametidy_windows_amd64.zip` or `nametidy_linux_amd64.tar.gz`).
3. Extract the archive. For `.tar.gz` files on Linux/macOS:
   ```bash
   tar -xzvf nametidy_VERSION_OS_ARCH.tar.gz
   ```
   For `.zip` files on Windows, you can use File Explorer's built-in "Extract All..." option.
4. Move the extracted `nametidy` (or `nametidy.exe` on Windows) executable to a directory in your system's PATH.
   - For Linux/macOS, a common location is `/usr/local/bin/`:
     ```bash
     sudo mv nametidy /usr/local/bin/
     ```
   - For Windows, you might choose a directory like `%USERPROFILE%\bin\` and ensure this directory is added to your PATH environment variable.
5. Run `nametidy --help` (or `nametidy.exe --help` on Windows) to verify the installation.

## Build

nametidy is written in Go. To install it locally:

1. Make sure Go is installed: https://golang.org/dl/
2. Clone the repository:

```bash
git clone https://github.com/mi8bi/nametidy.git
```

3. Build the project with Go:

```bash
cd nametidy
go build
```


## Usage

Organize and standardize file names in your target directory using intuitive subcommands.

### Clean Up Filenames
Removes unwanted characters, converts spaces to underscores, and standardizes file names.

```bash
nametidy clean -p ./test_dir
```

#### Example Output:

```
Renamed: ./test_dir/file (1).txt → ./test_dir/file_1.txt
Renamed: ./test_dir/hello world.txt → ./test_dir/hello_world.txt
History file path: ./test_dir/.nametidy_history
```


### Undo Changes
Restores the most recent file renaming performed by nametidy

```bash
nametidy undo -p ./test_dir
```

#### Example Output:

```
Restored: ./test_dir/file_1.txt → ./test_dir/file (1).txt
Restored: ./test_dir/hello_world.txt → ./test_dir/hello world.txt
```


### Dry Run
Displays changes without modifying any files.

```bash
nametidy clean -p ./test_dir -d
```

#### Example Output:

```
[DRY-RUN] ./test_dir/file (1).txt → ./test_dir/file_1.txt
[DRY-RUN] ./test_dir/hello world.txt → ./test_dir/hello_world.txt
```


### Verbose Logging
Enables detailed logs of the renaming process.

```bash
nametidy clean -p ./test_dir -v
```

#### Example Output:

```
2025/03/30 17:39:08 [INFO] Starting file name cleanup...
Renamed: ./test_dir/file (1).txt → ./test_dir/file_1.txt
Renamed: ./test_dir/hello world.txt → ./test_dir/hello_world.txt
History file path: ./test_dir/.nametidy_history
2025/03/30 17:39:08 [INFO] File name cleanup completed.
```

### Add Sequence Numbers
Adds numerical prefixes to file names. Use -n to set digit length, and -H for hierarchical mode.

```bash
nametidy number -p ./test_dir -n 3
```

#### Example Output:

```
Renamed: ./test_dir/image.png → ./test_dir/001_image.png
Renamed: ./test_dir/photo.jpg → ./test_dir/002_photo.jpg
```


```bash
nametidy number -p ./test_dir -n 3 -H
```

#### Example Output:

```
Renamed: ./test_dir/folder1/doc.txt → ./test_dir/folder1/001_doc.txt
Renamed: ./test_dir/folder1/note.pdf → ./test_dir/folder1/002_note.pdf
Renamed: ./test_dir/folder2/image.png → ./test_dir/folder2/001_image.png
```

## Options

| Option / Command      | Description |
|-----------------------|-------------|
| `clean`               | Cleans up file names (e.g., removes symbols, replaces spaces). |
| `number`              | Adds sequence numbers to file names. |
| `undo`                | Reverts the most recent operation. |
| `-p <path>`           | (Required) Target directory to process. |
| `-n <digits>`         | Sets the number of digits for sequence numbers (e.g., `-n 3` → 001, 002). |
| `-H`                  | Enables hierarchical numbering by folder. |
| `-d`                  | Dry run mode — preview changes without applying them. |
| `-v`                  | Verbose output — shows logs during execution. |

## License

This project is licensed under the MIT License. For more details, see the [LICENSE](LICENSE) file.

# GitHub Topics
[cli](https://github.com/topics/cli) [golang](https://github.com/topics/golang) [file-management](https://github.com/topics/file-management) [rename-files](https://github.com/topics/rename-files)
