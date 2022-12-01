package methods

import (
	"fmt"
	"strconv"
	"time"

	"github.com/alaleks/shortener/internal/app/database"
	"github.com/alaleks/shortener/internal/app/database/models"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDB(dsn string) Database {
	var d Database

	db, err := database.Connect(dsn)

	if err != nil {
		return d
	}

	d.DB = db

	return d
}

func (d Database) Close() error {
	if d.DB == nil {
		return nil
	}

	dbInstance, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("error getting database instance: %w", err)
	}

	err = dbInstance.Close()
	if err != nil {
		err = fmt.Errorf("error close database: %w", err)
	}

	return err
}

func (d Database) GetUser(uid int) models.Users {
	var user models.Users

	d.DB.Where("uid = ?", uid).First(&user)

	return user
}

func (d Database) AddUser() uint {
	user := models.Users{CreatedAt: time.Now()}

	d.DB.Create(&user)

	return user.UID
}

func (d Database) GetUrlsUser(uid int) []models.Urls {
	var urls []models.Urls

	d.DB.Where("uid = ?", uid).Find(&urls)

	return urls
}

func (d Database) AddURL(userID, shortUID, longURL string) {
	userIDtoInt, _ := strconv.Atoi(userID)

	d.DB.Create(&models.Urls{
		ShortUID: shortUID, LongURL: longURL,
		CreatedAt: time.Now(), UID: uint(userIDtoInt),
	})
}

func (d Database) GetOriginalURL(shortUID string) string {
	var url models.Urls

	d.DB.Where("short_uid = ?", shortUID).First(&url)

	return url.LongURL
}

func (d Database) UpdateStat(shortUID string) {
	var url models.Urls

	d.DB.Where("short_uid = ?", shortUID).First(&url)

	url.Statistics++

	d.DB.Save(&url)
}

func (d Database) GetUrlsUserHandler(uid int) []struct {
	ShotrURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
} {
	var usersURL []struct {
		ShotrURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	urls := d.GetUrlsUser(uid)

	if len(urls) == 0 {
		return usersURL
	}

	for _, v := range urls {
		usersURL = append(usersURL, struct {
			ShotrURL    string `json:"short_url"`
			OriginalURL string `json:"original_url"`
		}{
			ShotrURL:    v.ShortUID,
			OriginalURL: v.LongURL,
		})
	}

	return usersURL
}

func (d Database) GetStat(shortUID string) struct {
	ShortURL  string `json:"shorturl"`
	LongURL   string `json:"longurl"`
	Usage     uint   `json:"usage"`
	CreatedAt string `json:"createdAt"`
} {
	var (
		url  models.Urls
		stat struct {
			ShortURL  string `json:"shorturl"`
			LongURL   string `json:"longurl"`
			Usage     uint   `json:"usage"`
			CreatedAt string `json:"createdAt"`
		}
	)

	d.DB.Where("short_uid = ?", shortUID).First(&url)

	if url.LongURL == "" {
		return stat
	}

	stat.ShortURL = url.ShortUID
	stat.LongURL = url.LongURL
	stat.Usage = url.Statistics
	stat.CreatedAt = url.CreatedAt.Format("02.01.2006 15:04:05")

	return stat
}
