// Package models implement data storage and presentation models.
package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Users represents the data model of a specific user.
type Users struct {
	CreatedAt time.Time `gorm:"default:NOW()"`
	UID       uint      `gorm:"primaryKey;autoIncrement;unique"`
}

// Urls represents the data model of a specific shortened URL.
type Urls struct {
	CreatedAt     time.Time `gorm:"default:NOW()"`
	ShortUID      string    `gorm:"primaryKey"`
	CorrelationID string
	LongURL       string `gorm:"unique;index"`
	Statistics    uint
	UID           uint
	Removed       bool
}

// Migrate starts auto-migration of models in database.
func Migrate(sqlDB *gorm.DB) error {
	err := sqlDB.AutoMigrate(&Users{}, &Urls{})
	if err != nil {
		err = fmt.Errorf("error automigrate: %w", err)
	}

	sqlDB.Exec("ALTER TABLE urls ADD FOREIGN KEY(uid) REFERENCES users(uid);")

	return err
}
