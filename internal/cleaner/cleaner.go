package cleaner

import (
	"NameTidy/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const HISTORY_FILE = ".NameTidy_History"

// Clean cleans up file names within the specified directory (recursively processes subdirectories)
func Clean(dirPath string, dryRun bool) error {
	entries := make(map[string]string)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore history file and directories
		if info.IsDir() || filepath.Base(path) == HISTORY_FILE {
			return nil
		}

		oldName := info.Name()
		newName := utils.CleanFileName(oldName)

		if oldName != newName {
			newPath := filepath.Join(filepath.Dir(path), newName)
			if dryRun {
				fmt.Printf("[DRY-RUN] %s → %s\n", path, newPath)
			} else {
				if err := os.Rename(path, newPath); err != nil {
					utils.Error("Rename failed", err)
					return err
				}
				fmt.Printf("Renamed: %s → %s\n", path, newPath)
				entries[path] = newPath
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	if !dryRun {
		if err := saveHistory(dirPath, entries); err != nil {
			utils.Error("Failed to save history", err)
			return err
		}
	}

	return nil
}

// saveHistory saves the rename history
func saveHistory(dirPath string, history map[string]string) error {
	data, err := json.Marshal(history)
	if err != nil {
		return err
	}

	historyFilePath := filepath.Join(dirPath, HISTORY_FILE)
	fmt.Printf("History file path: %s\n", historyFilePath)

	return os.WriteFile(historyFilePath, data, 0644)
}

// NumberFiles adds numbering to file names
func NumberFiles(dirPath string, digits int, hierarchical bool, dryRun bool) error {
	// Manage count per directory
	counts := make(map[string]int)

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore history file and directories
		if info.IsDir() || filepath.Base(path) == HISTORY_FILE {
			return nil
		}

		var dirKey string
		if hierarchical {
			// Reset count for each directory
			dirKey = filepath.Dir(path)
		} else {
			// Global count for all files
			dirKey = "global"
		}

		// Increment count
		counts[dirKey]++
		count := counts[dirKey]

		// Generate new file name
		newPath, err := utils.AddNumbering(path, digits, count)
		if err != nil {
			return err
		}

		if dryRun {
			fmt.Printf("[DRY-RUN] %s → %s\n", path, newPath)
		} else {
			if err := os.Rename(path, newPath); err != nil {
				return fmt.Errorf("failed to rename the file: %v", err)
			}
			fmt.Printf("Renamed: %s → %s\n", path, newPath)
		}
		return nil
	})
}

// Undo undoes the most recent rename operation
func Undo(dirPath string, dryRun bool) error {
	history, err := loadHistory(dirPath)
	if err != nil {
		return err
	}

	for oldPath, newPath := range history {
		if oldPath == HISTORY_FILE || newPath == HISTORY_FILE {
			continue // Ignore history file
		}

		if utils.FileExists(newPath) {
			if dryRun {
				fmt.Printf("[DRY-RUN] %s → %s\n", newPath, oldPath)
			} else {
				err := utils.RenameFile(newPath, oldPath, false)
				if err != nil {
					utils.Error("Failed to undo the rename", err)
				}
			}
		}
	}

	if !dryRun {
		return os.Remove(filepath.Join(dirPath, HISTORY_FILE))
	}

	return nil
}

// loadHistory loads the rename history
func loadHistory(dirPath string) (map[string]string, error) {
	data, err := os.ReadFile(filepath.Join(dirPath, HISTORY_FILE))
	if err != nil {
		return nil, errors.New("rename history not found")
	}

	history := make(map[string]string)
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, errors.New("failed to read history data")
	}

	return history, nil
}
