package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Users struct {
	CreatedAt time.Time
	UID       uint `gorm:"primaryKey;autoIncrement;unique"`
}

type Urls struct {
	CreatedAt     time.Time
	ShortUID      string `gorm:"primaryKey"`
	CorrelationID string
	LongURL       string `gorm:"unique;index"`
	Statistics    uint
	UID           uint
	Removed       bool
}

func Migrate(sqlDB *gorm.DB) error {
	err := sqlDB.AutoMigrate(&Users{}, &Urls{})
	if err != nil {
		err = fmt.Errorf("error automigrate: %w", err)
	}

	sqlDB.Exec("ALTER TABLE urls ADD FOREIGN KEY(uid) REFERENCES users(uid);")

	return err
}
