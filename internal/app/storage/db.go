package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/alaleks/shortener/internal/app/config"
	"github.com/alaleks/shortener/internal/app/service"
	"github.com/alaleks/shortener/internal/app/storage/models"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	MaxIdleConns = 25
	MaxOpenConns = 25
	MaxLifetime  = time.Hour
)

var (
	ErrAlreadyExists = errors.New("such an entry exists in the database")
	ErrDBConnection  = errors.New("failed to check database connection")
	ErrInvalidData   = errors.New("data invalid")
)

type DB struct {
	db   *gorm.DB
	conf config.Configurator
	mu   *sync.RWMutex
}

func NewDB(conf config.Configurator) *DB {
	return &DB{conf: conf, mu: &sync.RWMutex{}}
}

func (d *DB) Init() error {
	db, err := gorm.Open(postgres.Open(d.conf.GetDSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	err = models.Migrate(db)

	if err != nil {
		return fmt.Errorf("error automigrating tables to database: %w", err)
	}

	d.db = db

	sqlDB, err := db.DB()

	if err == nil {
		sqlDB.SetMaxIdleConns(MaxIdleConns)
		sqlDB.SetMaxOpenConns(MaxOpenConns)
		sqlDB.SetConnMaxLifetime(MaxLifetime)
	}

	return nil
}

func (d *DB) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("SQL instance error: %w", err)
	}

	err = sqlDB.Close()

	if err != nil {
		return fmt.Errorf("db connection closed error: %w", err)
	}

	return err
}

func (d *DB) Ping() error {
	if d.db == nil {
		return ErrDBConnection
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("SQL instance error: %w", err)
	}

	err = sqlDB.Ping()

	if err != nil {
		return fmt.Errorf("ping db error: %w", err)
	}

	return err
}

func PingDB(sqlDB *sql.DB) error {
	err := sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("ping db error: %w", err)
	}

	return err
}

func (d *DB) Add(longURL, userID string) (string, error) {
	if d.Ping() != nil {
		return "", ErrDBConnection
	}

	userIDtoInt, err := strconv.Atoi(userID)

	uri := models.Urls{
		ShortUID: service.GenUID(d.conf.GetSizeUID()), LongURL: longURL,
		CreatedAt: time.Now(),
	}

	if err == nil {
		uri.UID = uint(userIDtoInt)
	}

	rowsAffected := writeURL(d.db, uri)

	if rowsAffected == 0 {
		return getShortUID(d.db, longURL, d.conf.GetBaseURL()), ErrAlreadyExists
	}

	return d.conf.GetBaseURL() + uri.ShortUID, nil
}

func writeURL(db *gorm.DB, uri models.Urls) int {
	res := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&uri)

	return int(res.RowsAffected)
}

func getShortUID(db *gorm.DB, longURL, baseURL string) string {
	var uri models.Urls

	db.Where("long_url = ?", longURL).Find(&uri)

	return baseURL + uri.ShortUID
}

func (d *DB) AddBatch(longURL, userID, corID string) string {
	if d.Ping() != nil {
		return ""
	}

	userIDtoInt, err := strconv.Atoi(userID)

	uri := models.Urls{
		ShortUID: service.GenUID(d.conf.GetSizeUID()), LongURL: longURL,
		CreatedAt: time.Now(),
	}

	if err == nil {
		uri.UID = uint(userIDtoInt)
	}

	rowsAffected := writeURLBatch(d.db, uri)

	if rowsAffected == 0 {
		return getShortUID(d.db, longURL, d.conf.GetBaseURL())
	}

	return d.conf.GetBaseURL() + uri.ShortUID
}

func writeURLBatch(db *gorm.DB, uri models.Urls) int {
	res := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&uri)

	return int(res.RowsAffected)
}

func (d *DB) Update(uid string) {
	if d.Ping() != nil {
		return
	}

	var url models.Urls

	res := d.db.Where("short_uid = ?", uid).First(&url)

	if res.Error != nil || res.RowsAffected == 0 {
		return
	}

	d.mu.Lock()
	url.Statistics++
	d.db.Save(&url)
	d.mu.Unlock()
}

func (d *DB) GetURL(uid string) (string, error) {
	if d.Ping() != nil {
		return "", ErrDBConnection
	}

	var url models.Urls

	res := d.db.Where("short_uid = ?", uid).First(&url)

	if res.RowsAffected == 0 {
		return url.LongURL, ErrUIDNotValid
	}

	if url.Removed {
		return url.LongURL, ErrShortURLDeleted
	}

	return url.LongURL, nil
}

func (d *DB) Stat(uid string) (Statistics, error) {
	if d.Ping() != nil {
		return Statistics{}, ErrDBConnection
	}

	var uri models.Urls

	res := d.db.Where("short_uid = ?", uid).First(&uri)

	stat := Statistics{
		ShortURL:  d.conf.GetBaseURL() + uid,
		LongURL:   uri.LongURL,
		CreatedAt: uri.CreatedAt.Format("02.01.2006 15:04:05"),
		Usage:     uri.Statistics,
	}

	if res.RowsAffected == 0 {
		return stat, ErrUIDNotValid
	}

	return stat, nil
}

func (d *DB) Create() uint {
	if d.Ping() != nil {
		return 0
	}

	user := models.Users{CreatedAt: time.Now()}

	d.db.Create(&user)

	return user.UID
}

/*
func getUser(db *gorm.DB, uid int) models.Users {
	var user models.Users

	db.Where("uid = ?", uid).First(&user)

	return user
}
*/

func getUrlsUser(db *gorm.DB, uid uint) []models.Urls {
	var urls []models.Urls

	db.Where("uid = ?", uid).Find(&urls)

	return urls
}

func (d *DB) GetUrlsUser(userID string) ([]URLUser, error) {
	if d.Ping() != nil {
		return []URLUser{}, ErrDBConnection
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		return []URLUser{}, ErrUserIDNotValid
	}

	urls := getUrlsUser(d.db, uint(uid))

	usersURL := make([]URLUser, 0, len(urls))

	if len(urls) == 0 {
		return usersURL, ErrUserUrlsEmpty
	}

	for _, item := range urls {
		usersURL = append(usersURL, URLUser{
			ShortURL:    item.ShortUID,
			OriginalURL: item.LongURL,
		})
	}

	return usersURL, nil
}

func (d *DB) DelUrls(userID string, shortsUID ...string) error {
	if d.Ping() != nil {
		return ErrDBConnection
	}

	if len(shortsUID) == 0 || userID == "" {
		return ErrInvalidData
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		return ErrUserIDNotValid
	}

	res := d.db.Model(models.Urls{}).
		Where("short_uid IN ? AND uid = ?", shortsUID, uid).
		Updates(models.Urls{Removed: true})

	return res.Error
}
