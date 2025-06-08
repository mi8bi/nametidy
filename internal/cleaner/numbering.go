package cleaner

import (
	"nametidy/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

func NumberFiles(db *gorm.DB, dirPath string, digits int, hierarchical bool, dryRun bool) error {
	counts := make(map[string]int)
	batchID := fmt.Sprintf("number-%d", time.Now().UnixNano())
	histories := []RenameHistory{}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		var dirKey string
		if hierarchical {
			dirKey = filepath.Dir(path)
		} else {
			dirKey = "global"
		}
		counts[dirKey]++
		count := counts[dirKey]

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
			histories = append(histories, RenameHistory{
				OriginalPath: path,
				NewPath:      newPath,
				Operation:    "number",
				BatchID:      batchID,
				CreatedAt:    time.Now(),
			})
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
