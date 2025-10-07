package cleaner

import (
	"fmt"

	"gorm.io/gorm"
)

// ClearHistory deletes all rename history records in the database
func ClearHistory(db *gorm.DB) error {
	result := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&RenameHistory{})
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("Deleted %d total history entries.\n", result.RowsAffected)
	return nil
}
