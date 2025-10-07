package cleaner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DB_FILE = ".name_tidy_history.db"

type RenameHistory struct {
	ID            uint `gorm:"primaryKey"`
	OriginalPath  string
	NewPath       string
	Operation     string // "clean" または "number"
	BatchID       string
	CreatedAt     time.Time
	Reverted      bool
	Redone        bool
	OperationType string // "undo" または "redo"
}

func GetDB() (*gorm.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(homeDir, DB_FILE)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&RenameHistory{}); err != nil {
		return nil, err
	}
	return db, nil
}

func SaveRenameHistory(db *gorm.DB, entries map[string]string, operation string) error {
	batchID := fmt.Sprintf("%d", time.Now().UnixNano())
	var records []RenameHistory
	for oldPath, newPath := range entries {
		records = append(records, RenameHistory{
			OriginalPath: oldPath,
			NewPath:      newPath,
			Operation:    operation,
			BatchID:      batchID,
			CreatedAt:    time.Now(),
			Reverted:     false,
			Redone:       false,
		})
	}
	return db.Create(&records).Error
}

func GetLastUndoableBatch(db *gorm.DB) (string, error) {
	var last RenameHistory
	if err := db.Where("reverted = ? AND redone = ?", false, false).
		Order("created_at desc").First(&last).Error; err != nil {
		return "", errors.New("no operation to undo")
	}
	return last.BatchID, nil
}

func GetLastRedoableBatch(db *gorm.DB) (string, error) {
	var last RenameHistory
	if err := db.Where("reverted = ? AND redone = ?", true, false).
		Order("created_at desc").First(&last).Error; err != nil {
		return "", errors.New("no operation to redo")
	}
	return last.BatchID, nil
}

func GetHistoriesByBatch(db *gorm.DB, batchID string) ([]RenameHistory, error) {
	var records []RenameHistory
	if err := db.Where("batch_id = ?", batchID).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
