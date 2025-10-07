package cleaner

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}

	if err := db.AutoMigrate(&RenameHistory{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return db
}

func insertDummyHistories(db *gorm.DB, count int) error {
	for i := 0; i < count; i++ {
		record := RenameHistory{
			OriginalPath: "dummy/original.txt",
			NewPath:      "dummy/renamed.txt",
			Operation:    "clean",
			BatchID:      "test-batch",
			CreatedAt:    time.Now(),
			Reverted:     false,
			Redone:       false,
		}
		if err := db.Create(&record).Error; err != nil {
			return err
		}
	}
	return nil
}

func TestClearHistory(t *testing.T) {
	db := setupTestDB(t)

	// Insert dummy history
	if err := insertDummyHistories(db, 3); err != nil {
		t.Fatalf("failed to insert dummy histories: %v", err)
	}

	var count int64
	db.Model(&RenameHistory{}).Count(&count)
	if count != 3 {
		t.Fatalf("expected 3 history records, got %d", count)
	}

	// Run ClearHistory
	if err := ClearHistory(db); err != nil {
		t.Fatalf("ClearHistory failed: %v", err)
	}

	db.Model(&RenameHistory{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected 0 history records after clear, got %d", count)
	}
}
