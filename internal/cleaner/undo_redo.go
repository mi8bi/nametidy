package cleaner

import (
	"NameTidy/internal/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func Undo(db *gorm.DB, dirPath string, dryRun bool) error {
	var lastBatch RenameHistory
	// 最新の未返却操作を取得
	if err := db.Where("reverted = ? AND redone = ?", false, false).
		Order("created_at desc").First(&lastBatch).Error; err != nil {
		return errors.New("no operation to undo")
	}

	// 同じバッチIDを持つ履歴をすべて取得
	var histories []RenameHistory
	if err := db.Where("batch_id = ?", lastBatch.BatchID).Find(&histories).Error; err != nil {
		return err
	}

	// それぞれの履歴をやり直す
	for _, h := range histories {
		if utils.FileExists(h.NewPath) {
			if dryRun {
				fmt.Printf("[DRY-RUN] %s → %s\n", h.NewPath, h.OriginalPath)
			} else {
				if err := utils.RenameFile(h.NewPath, h.OriginalPath, false); err != nil {
					utils.Warn(fmt.Sprintf("Undo failed: %s → %s", h.NewPath, h.OriginalPath))
				}
			}
		}
	}

	// 履歴をrevertedとしてマーク
	if !dryRun {
		db.Model(&RenameHistory{}).
			Where("batch_id = ?", lastBatch.BatchID).
			Update("reverted", true).
			Update("operation_type", "undo") // "undo"としてマーク
	}

	return nil
}

func Redo(db *gorm.DB, dirPath string, dryRun bool) error {
	var lastUndone RenameHistory
	// 最新の戻された操作を取得
	if err := db.Where("reverted = ? AND redone = ?", true, false).
		Order("created_at desc").First(&lastUndone).Error; err != nil {
		return errors.New("no operation to redo")
	}

	// 同じバッチIDを持つ履歴をすべて取得
	var histories []RenameHistory
	if err := db.Where("batch_id = ?", lastUndone.BatchID).Find(&histories).Error; err != nil {
		return err
	}

	// それぞれの履歴を元に戻す
	for _, h := range histories {
		if utils.FileExists(h.OriginalPath) {
			if dryRun {
				fmt.Printf("[DRY-RUN] %s → %s\n", h.OriginalPath, h.NewPath)
			} else {
				if err := utils.RenameFile(h.OriginalPath, h.NewPath, false); err != nil {
					utils.Warn(fmt.Sprintf("Redo failed: %s → %s", h.OriginalPath, h.NewPath))
				}
			}
		}
	}

	// 履歴をredoneとしてマーク
	if !dryRun {
		db.Model(&RenameHistory{}).
			Where("batch_id = ?", lastUndone.BatchID).
			Update("reverted", false).
			Update("operation_type", "redo") // "redo"としてマーク
	}

	return nil
}
