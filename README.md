# NameTidy

NameTidy is a CLI tool that helps users easily rename and organize files. It provides functionalities such as file name cleanup, renaming, numbering, and undoing actions.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Options](#options)
- [License](#license)

## Installation

NameTidy is developed in Go. You can install it by following these steps:

1. Ensure that Go is installed on your system.
2. Clone this repository:
   ```bash
   git clone https://github.com/mi8bi/NameTidy.git
   ```
3. Build the project with Go:
   ```bash
   cd NameTidy
   go build
   ```

## Usage

You can easily organize file names within a specified directory using `NameTidy`. The following commands allow you to perform various actions:

### Clean Up
This command cleans up file names by removing unwanted characters or formatting, converting them to a standard format.

```bash
NameTidy clean -p ./test_dir
```

### Undo Rename (Undo)
This command undoes the most recent rename operation.

```bash
NameTidy undo -p ./test_dir
```

### Dry Run
This command shows the intended changes without actually modifying any files.

```bash
NameTidy clean -p ./test_dir -d
```

### Verbose Mode
This command provides detailed log output during processing.

```bash
NameTidy clean -p ./test_dir -v
```

### Numbering Files (Numbered)
This command adds a sequence number to file names. You can specify the number of digits or apply hierarchical numbering based on the directory structure.

```bash
NameTidy number -p ./test_dir -n 3 
NameTidy number -p ./test_dir -n 3 -H
```

## Options

- `clean`: Cleans up the file names.
- `undo`: Undoes the most recent rename operation.
- `-d`: Displays the intended changes without applying them.
- `-v`: Provides detailed log output.
- `number`: Adds sequence numbers to the file names.
  - `-n <digits>`: Specifies the number of digits for the sequence number.
  - `-n <digits> -H`: Adds hierarchical numbering based on directory structure.

## License

This project is licensed under the MIT License. For more details, see the [LICENSE](LICENSE) file.
