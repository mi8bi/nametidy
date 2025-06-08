package cleaner

import (
	"nametidy/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

func Clean(db *gorm.DB, dirPath string, dryRun bool) error {
	batchID := fmt.Sprintf("clean-%d", time.Now().UnixNano())
	histories := []RenameHistory{}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
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
				histories = append(histories, RenameHistory{
					OriginalPath: path,
					NewPath:      newPath,
					Operation:    "clean",
					BatchID:      batchID,
					CreatedAt:    time.Now(),
				})
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	if !dryRun && len(histories) > 0 {
		return db.Create(&histories).Error
	}
	return nil
}