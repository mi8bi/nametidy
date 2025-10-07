package utils

import (
	"fmt"
	"path/filepath"
)

// AddNumbering adds a sequence number to the file name
func AddNumbering(path string, digits int, index int) (string, error) {
	dir, file := filepath.Split(path)
	newName := generateNumberedName(file, digits, index)

	// Create the new file path
	newPath := filepath.Join(dir, newName)
	return newPath, nil
}

// generateNumberedName generates a numbered name for the file
func generateNumberedName(baseName string, digits int, index int) string {
	indexStr := fmt.Sprintf("%0*d", digits, index)
	ext := filepath.Ext(baseName)
	fileName := baseName[:len(baseName)-len(ext)]
	return fmt.Sprintf("%s_%s%s", indexStr, fileName, ext)
}
