package cleaner

import (
	"NameTidy/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

// osRenameFunc is a variable that holds the function to rename files.
// It defaults to os.Rename but can be changed for testing.
var osRenameFunc = os.Rename

// ProcessFiles is a generic function to process files in a directory.
// It walks through the directory, applies the processFileFunc to each file,
// and handles renaming and history logging.
func ProcessFiles(
	db *gorm.DB,
	dirPath string,
	operation string, // e.g., "clean", "number"
	dryRun bool,
	// processFileFunc takes the full file path and its os.FileInfo,
	// returns oldName (just filename), newName (just filename), or an error.
	processFileFunc func(filePath string, info os.FileInfo) (oldName string, newName string, err error),
) error {
	batchID := fmt.Sprintf("%d", time.Now().UnixNano())
	var histories []RenameHistory

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Propagate errors from walking (e.g., permission issues)
		}

		// Skip directories
		if info.IsDir() {
			// For hierarchical numbering, we might need to process directories.
			// However, the core file processing and renaming applies to files.
			// The decision to descend into a directory is handled by filepath.Walk itself.
			// If a directory itself needs to be "processed" (e.g. its name changed),
			// that would be a different type of operation not covered by this specific ProcessFiles.
			return nil
		}

		// Apply the custom processing function to get old and new names
		oldName, newName, processErr := processFileFunc(path, info)
		if processErr != nil {
			utils.Error(fmt.Sprintf("Error processing file %s: %v", path, processErr), processErr)
			return nil // Continue with other files
		}

		if oldName != newName {
			originalPath := path
			newPath := filepath.Join(filepath.Dir(path), newName)

			if dryRun {
				utils.Info(fmt.Sprintf("[Dry Run] Would rename %s to %s", originalPath, newPath))
			} else {
				err := osRenameFunc(originalPath, newPath) // Use the function variable
				if err != nil {
					utils.Error(fmt.Sprintf("Error renaming %s to %s: %v", originalPath, newPath, err), err)
					return nil // Continue with other files, but log this error
				}
				utils.Info(fmt.Sprintf("Renamed %s to %s", originalPath, newPath))
				histories = append(histories, RenameHistory{
					OriginalPath:  originalPath,
					NewPath:       newPath,
					Operation:     operation,
					BatchID:       batchID,
					CreatedAt:     time.Now(),
					Reverted:      false,
					Redone:        false,
					OperationType: "", // OperationType is for undo/redo itself
				})
			}
		}
		return nil
	})

	if err != nil { // Error from filepath.Walk itself
		return fmt.Errorf("error walking the path %s: %w", dirPath, err)
	}

	if !dryRun && len(histories) > 0 {
		result := db.Create(&histories)
		if result.Error != nil {
			return fmt.Errorf("failed to save rename history: %w", result.Error)
		}
		utils.Info(fmt.Sprintf("Saved %d rename operations to history with batch ID %s.", len(histories), batchID))
	} else if !dryRun && len(histories) == 0 {
		utils.Info("No files needed renaming.")
	}

	return nil
}
