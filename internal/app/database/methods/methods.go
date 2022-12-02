package methods

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/alaleks/shortener/internal/app/database"
	"github.com/alaleks/shortener/internal/app/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrIsExist = errors.New("such an entry exists in the database")

type Database struct {
	DB *gorm.DB
}

func OpenDB(dsn string) Database {
	var dBase Database

	db, err := database.Connect(dsn)
	if err != nil {
		return dBase
	}

	dBase.DB = db

	return dBase
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

func (d Database) GetShortUID(longURL string) string {
	var uri models.Urls

	d.DB.Where("long_url = ?", longURL).Find(&uri)

	return uri.ShortUID
}

func (d Database) AddURL(userID, shortUID, longURL string) (string, error) {
	userIDtoInt, _ := strconv.Atoi(userID)

	uri := models.Urls{
		ShortUID: shortUID, LongURL: longURL,
		CreatedAt: time.Now(), UID: uint(userIDtoInt),
	}

	res := d.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&uri)

	if res.RowsAffected == 0 {
		return d.GetShortUID(longURL), ErrIsExist
	}

	return uri.ShortUID, nil
}

func (d Database) AddURLBatch(userID, shortUID, corID, longURL string) string {
	userIDtoInt, _ := strconv.Atoi(userID)

	uri := models.Urls{
		ShortUID: shortUID, LongURL: longURL,
		CreatedAt: time.Now(), UID: uint(userIDtoInt),
		CorrelationID: corID,
	}

	res := d.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&uri)

	if res.RowsAffected == 0 {
		return d.GetShortUID(longURL)
	}

	return uri.ShortUID
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
	urls := d.GetUrlsUser(uid)

	usersURL := make([]struct {
		ShotrURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}, 0, len(urls))

	if len(urls) == 0 {
		return usersURL
	}

	for _, item := range urls {
		usersURL = append(usersURL, struct {
			ShotrURL    string `json:"short_url"`
			OriginalURL string `json:"original_url"`
		}{
			ShotrURL:    item.ShortUID,
			OriginalURL: item.LongURL,
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
