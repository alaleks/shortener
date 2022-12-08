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

var ErrAlreadyExists = errors.New("such an entry exists in the database")

type Database struct {
	SDB *gorm.DB
}

func OpenDB(dsn string) (Database, error) {
	db, err := database.Connect(dsn)

	return Database{SDB: db}, err
}

func (d Database) Close() error {
	if d.SDB == nil {
		return nil
	}

	dbInstance, err := d.SDB.DB()
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

	d.SDB.Where("uid = ?", uid).First(&user)

	return user
}

func (d Database) AddUser() uint {
	user := models.Users{CreatedAt: time.Now()}

	d.SDB.Create(&user)

	return user.UID
}

func (d Database) GetUrlsUser(uid int) []models.Urls {
	var urls []models.Urls

	d.SDB.Where("uid = ?", uid).Find(&urls)

	return urls
}

func (d Database) GetShortUID(longURL string) string {
	var uri models.Urls

	d.SDB.Where("long_url = ?", longURL).Find(&uri)

	return uri.ShortUID
}

func (d Database) AddURL(userID, shortUID, longURL string) (string, error) {
	userIDtoInt, err := strconv.Atoi(userID)

	uri := models.Urls{
		ShortUID: shortUID, LongURL: longURL,
		CreatedAt: time.Now(),
	}

	if err == nil {
		uri.UID = uint(userIDtoInt)
	}

	rowsAffected := writeURL(d.SDB, uri)

	if rowsAffected == 0 {
		return d.GetShortUID(longURL), ErrAlreadyExists
	}

	return uri.ShortUID, nil
}

func writeURL(db *gorm.DB, uri models.Urls) int {
	res := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&uri)

	return int(res.RowsAffected)
}

func (d Database) AddURLBatch(userID, shortUID, corID, longURL string) string {
	userIDtoInt, err := strconv.Atoi(userID)

	uri := models.Urls{
		ShortUID: shortUID, LongURL: longURL,
		CreatedAt: time.Now(),
	}

	if err == nil {
		uri.UID = uint(userIDtoInt)
	}

	rowsAffected := writeURLBatch(d.SDB, uri)

	if rowsAffected == 0 {
		return d.GetShortUID(longURL)
	}

	return uri.ShortUID
}

func writeURLBatch(db *gorm.DB, uri models.Urls) int {
	res := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&uri)

	return int(res.RowsAffected)
}

func (d Database) GetOriginalURL(shortUID string) string {
	var url models.Urls

	d.SDB.Where("short_uid = ?", shortUID).First(&url)

	return url.LongURL
}

func (d Database) UpdateStat(shortUID string) {
	var url models.Urls

	d.SDB.Where("short_uid = ?", shortUID).First(&url)

	url.Statistics++

	d.SDB.Save(&url)
}

func (d Database) GetUrlsUserHandler(uid int) []struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
} {
	urls := d.GetUrlsUser(uid)

	usersURL := make([]struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}, 0, len(urls))

	if len(urls) == 0 {
		return usersURL
	}

	for _, item := range urls {
		usersURL = append(usersURL, struct {
			ShortURL    string `json:"short_url"`
			OriginalURL string `json:"original_url"`
		}{
			ShortURL:    item.ShortUID,
			OriginalURL: item.LongURL,
		})
	}

	return usersURL
}

func (d Database) GetStat(shortUID string) struct {
	ShortURL  string `json:"shorturl"`
	LongURL   string `json:"longurl"`
	CreatedAt string `json:"createdAt"`
	Usage     uint   `json:"usage"`
} {
	var (
		url  models.Urls
		stat struct {
			ShortURL  string `json:"shorturl"`
			LongURL   string `json:"longurl"`
			CreatedAt string `json:"createdAt"`
			Usage     uint   `json:"usage"`
		}
	)

	d.SDB.Where("short_uid = ?", shortUID).First(&url)

	if url.LongURL == "" {
		return stat
	}

	stat.ShortURL = url.ShortUID
	stat.LongURL = url.LongURL
	stat.Usage = url.Statistics
	stat.CreatedAt = url.CreatedAt.Format("02.01.2006 15:04:05")

	return stat
}
