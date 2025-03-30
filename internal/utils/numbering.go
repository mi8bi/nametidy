package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// AddNumbering adds a sequence number to the file name
func AddNumbering(path string, digits int, index int) (string, error) {
	dir, file := filepath.Split(path)
	newName := GenerateNumberedName(file, digits, index)

	// Create the new file path
	newPath := filepath.Join(dir, newName)
	return newPath, nil
}

// GenerateNumberedName generates a numbered name for the file
func GenerateNumberedName(baseName string, digits int, index int) string {
	indexStr := fmt.Sprintf("%0*d", digits, index)
	ext := filepath.Ext(baseName)
	fileName := baseName[:len(baseName)-len(ext)]
	return fmt.Sprintf("%s_%s%s", indexStr, fileName, ext)
}

// ProcessDirectory adds sequence numbers to the files in the directory
func ProcessDirectory(dirPath string, digits int, hierarchical bool) error {
	counts := make(map[string]int)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			var dirKey string
			if hierarchical {
				dirKey = filepath.Dir(path)
			} else {
				dirKey = "global"
			}

			// Increment the count
			counts[dirKey]++
			count := counts[dirKey]

			// Add the sequence number to the file
			newPath, err := AddNumbering(path, digits, count)
			if err != nil {
				return err
			}

			// Rename the file
			if err := os.Rename(path, newPath); err != nil {
				return fmt.Errorf("failed to rename file: %v", err)
			}
			fmt.Printf("Renamed: %s → %s\n", path, newPath)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error occurred while processing the directory: %v", err)
	}
	return nil
}