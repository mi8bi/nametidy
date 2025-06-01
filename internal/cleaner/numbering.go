package cleaner

import (
	"NameTidy/internal/utils"
	"fmt"
	"os"
	"path/filepath" // Still needed for dirKey logic within the callback

	"gorm.io/gorm"
)

// NumberFiles adds sequence numbers to file names in a directory.
// It uses the generic ProcessFiles function with a specific numbering callback.
func NumberFiles(db *gorm.DB, dirPath string, digits int, hierarchical bool, dryRun bool) error {
	// counts map needs to be managed by the calling function and captured by the closure.
	counts := make(map[string]int)

	processFileFunc := func(filePath string, info os.FileInfo) (string, string, error) {
		oldName := info.Name()
		var dirKey string
		if hierarchical {
			dirKey = filepath.Dir(filePath)
		} else {
			dirKey = "_global_" // Use a distinct key for non-hierarchical to avoid collisions with actual dir names
		}
		counts[dirKey]++
		count := counts[dirKey]

		// utils.AddNumbering is expected to return the new full path, or just new name.
		// Based on its usage in the original code: `newPath, err := utils.AddNumbering(path, digits, count)`
		// it seems to return the new full path.
		// However, ProcessFiles expects oldName (filename only) and newName (filename only).
		// Let's assume there's a utils.FormatNumberedFileName(oldName, digits, count) or similar
		// that returns just the new file name.
		// For now, I will adapt assuming utils.AddNumbering can be used to derive the new name.
		// This might need adjustment if utils.AddNumbering strictly returns a full path.

		// Let's try to construct the new name based on the old name and numbering.
		// This is a common pattern: prefix-number.extension or number-name.extension
		// The original utils.AddNumbering took the full path.
		// We need a utility that gives us the new name, or we derive it here.

		// Simplification: Assume utils.GenerateNumberedFileName(oldName, digits, count, extension)
		// For now, let's mimic what utils.AddNumbering might have done to get the new name.
		// This part is tricky without knowing the exact behavior of utils.AddNumbering
		// and how it forms the new name.

		// Let's assume utils.AddNumbering can take the old name and return a new name
		// or that we have a similar utility.
		// If utils.AddNumbering returns a full path, we'd do:
		// tempNewPath, err := utils.AddNumbering(filePath, digits, count)
		// if err != nil { return oldName, oldName, err }
		// newName = filepath.Base(tempNewPath)

		// Given the original code: newPath, err := utils.AddNumbering(path, digits, count)
		// The new ProcessFiles expects (oldName, newName, error) where names are file names, not paths.
		// So, the callback should be:
		tempNewFullPath, err := utils.AddNumbering(filePath, digits, count)
		if err != nil {
			// If AddNumbering fails, return oldName for both to signify no change, and the error
			return oldName, oldName, fmt.Errorf("failed to determine new numbered name for %s: %w", oldName, err)
		}
		newName := filepath.Base(tempNewFullPath)
		return oldName, newName, nil
	}

	return ProcessFiles(db, dirPath, "number", dryRun, processFileFunc)
}
