package cmd

import (
	"NameTidy/internal/cleaner"
	"NameTidy/internal/utils"
	"fmt"

	"gorm.io/gorm"
)

// handleCommonInitializations performs common initialization tasks for commands.
// It initializes the logger and the database.
// If checkDir is true, it also checks if dirPath is a valid directory.
func handleCommonInitializations(verbose bool, dirPath string, checkDir bool) (*gorm.DB, error) {
	// Initialize logger
	utils.InitLogger(verbose)

	// Check if directory exists, if required
	if checkDir {
		if !utils.IsDirectory(dirPath) {
			// Return an error that can be handled by utils.Error
			return nil, fmt.Errorf("the specified directory '%s' does not exist or is not a directory", dirPath)
		}
	}

	// Initialize DB
	db, err := cleaner.GetDB()
	if err != nil {
		// Return an error that can be handled by utils.Error
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	return db, nil
}
