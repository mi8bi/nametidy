> ✨ A fast and simple CLI tool to clean and rename file names — powered by Go.

# [NameTidy](https://mi8bi.github.io/NameTidy/)

NameTidy is a fast and flexible command-line tool for cleaning and renaming file names.  
It supports operations such as filename cleanup, adding sequence numbers, and undoing changes — all with a simple and intuitive interface.

[![Build Status](https://github.com/mi8bi/NameTidy/actions/workflows/test.yml/badge.svg)](https://github.com/mi8bi/NameTidy/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mi8bi/NameTidy)](https://goreportcard.com/report/github.com/mi8bi/NameTidy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest Release](https://img.shields.io/github/v/release/mi8bi/NameTidy)](https://github.com/mi8bi/NameTidy/releases/latest)
[![GitHub Stars](https://img.shields.io/github/stars/mi8bi/NameTidy?style=social)](https://github.com/mi8bi/NameTidy/stargazers)

---

## 📽️ Demo

See NameTidy in action:

[![NameTidy_Demo002](https://asciinema.org/a/719898.svg)](https://asciinema.org/a/719898)

---

## Table of Contents

- [Download](#download)
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

## Download

You can download prebuilt binaries from the [GitHub Releases page](https://github.com/mi8bi/NameTidy/releases):

1. Go to the [Releases](https://github.com/mi8bi/NameTidy/releases) page on GitHub.
2. Find the latest release and download the binary file for your OS and architecture.
3. Extract the archive if needed, for example:

```bash
tar -xzvf NameTidy_0.1.0_linux_amd64.tar.gz
```

4. Move the executable to a directory in your PATH, for example:

```bash
mv NameTidy /usr/local/bin/
```

5. Run NameTidy --help to verify installation.

## Build

NameTidy is written in Go. To install it locally:

1. Make sure Go is installed: https://golang.org/dl/
2. Clone the repository:

```bash
git clone https://github.com/mi8bi/NameTidy.git
```

3. Build the project with Go:

```bash
cd NameTidy
go build
```


## Usage

Organize and standardize file names in your target directory using intuitive subcommands.

### Clean Up Filenames
Removes unwanted characters, converts spaces to underscores, and standardizes file names.

```bash
NameTidy clean -p ./test_dir
```

#### Example Output:

```
Renamed: ./test_dir/file (1).txt → ./test_dir/file_1.txt
Renamed: ./test_dir/hello world.txt → ./test_dir/hello_world.txt
History file path: ./test_dir/.NameTidy_History
```


### Undo Changes
Restores the most recent file renaming performed by NameTidy

```bash
NameTidy undo -p ./test_dir
```

#### Example Output:

```
Restored: ./test_dir/file_1.txt → ./test_dir/file (1).txt
Restored: ./test_dir/hello_world.txt → ./test_dir/hello world.txt
```


### Dry Run
Displays changes without modifying any files.

```bash
NameTidy clean -p ./test_dir -d
```

#### Example Output:

```
[DRY-RUN] ./test_dir/file (1).txt → ./test_dir/file_1.txt
[DRY-RUN] ./test_dir/hello world.txt → ./test_dir/hello_world.txt
```


### Verbose Logging
Enables detailed logs of the renaming process.

```bash
NameTidy clean -p ./test_dir -v
```

#### Example Output:

```
2025/03/30 17:39:08 [INFO] Starting file name cleanup...
Renamed: ./test_dir/file (1).txt → ./test_dir/file_1.txt
Renamed: ./test_dir/hello world.txt → ./test_dir/hello_world.txt
History file path: ./test_dir/.NameTidy_History
2025/03/30 17:39:08 [INFO] File name cleanup completed.
```

### Add Sequence Numbers
Adds numerical prefixes to file names. Use -n to set digit length, and -H for hierarchical mode.

```bash
NameTidy number -p ./test_dir -n 3
```

#### Example Output:

```
Renamed: ./test_dir/image.png → ./test_dir/001_image.png
Renamed: ./test_dir/photo.jpg → ./test_dir/002_photo.jpg
```


```bash
NameTidy number -p ./test_dir -n 3 -H
```

#### Example Output:

```
Renamed: ./test_dir/folder1/doc.txt → ./test_dir/folder1/001_doc.txt
Renamed: ./test_dir/folder1/note.pdf → ./test_dir/folder1/002_note.pdf
Renamed: ./test_dir/folder2/image.png → ./test_dir/folder2/001_image.png
```

## ⚙️ Options

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
