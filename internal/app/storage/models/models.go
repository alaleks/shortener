// Package models implement data storage and presentation models.
package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Users view model of users
type Users struct {
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
	UID       uint      `gorm:"primaryKey;autoIncrement;unique"`
}

// Urls representation model for short URL.
type Urls struct {
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
	ShortUID      string    `gorm:"primaryKey"`
	CorrelationID string
	LongURL       string `gorm:"unique;index"`
	Statistics    uint
	UID           uint
	Removed       bool
}

// Migrate starts auto-migration of models.
func Migrate(sqlDB *gorm.DB) error {
	err := sqlDB.AutoMigrate(&Users{}, &Urls{})
	if err != nil {
		err = fmt.Errorf("error automigrate: %w", err)
	}

	sqlDB.Exec("ALTER TABLE urls ADD FOREIGN KEY(uid) REFERENCES users(uid);")

	return err
}
