package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// IsDirectory checks if the specified path is a directory
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// CleanFileName cleans up the file name by replacing unwanted characters
func CleanFileName(fileName string) string {
	// Remove file extension
	ext := filepath.Ext(fileName)
	baseName := fileName[:len(fileName)-len(ext)]

	// Replace non-alphanumeric characters (except dot) with an underscore
	reClean := regexp.MustCompile(`[^\w\d.]`)
	baseName = reClean.ReplaceAllString(baseName, "_")

	// Replace consecutive underscores with a single underscore
	reUnderscore := regexp.MustCompile(`_+`)
	baseName = reUnderscore.ReplaceAllString(baseName, "_")

	// Remove leading and trailing underscores
	baseName = strings.Trim(baseName, "_")

	// Restore file extension
	return baseName + ext
}

// RenameFile renames a file from oldPath to newPath
func RenameFile(oldPath, newPath string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[DRY-RUN] %s → %s\n", oldPath, newPath)
		return nil
	}

	err := os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("Rename failed: %v", err)
	}

	fmt.Printf("Renamed: %s → %s\n", filepath.Base(oldPath), filepath.Base(newPath))
	return nil
}

// FileExists checks if the specified file exists at the given path
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
