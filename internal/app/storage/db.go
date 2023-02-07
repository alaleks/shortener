package storage

import (
	"database/sql"
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
	maxIdleConns = 100
	maxOpenConns = 200
	maxLifetime  = (15 * time.Minute)
)

// DB represents a database instance
type DB struct {
	db   *gorm.DB
	conf config.Configurator
	mu   sync.RWMutex
}

// NewDB creates a pointer of DB instance.
func NewDB(conf config.Configurator) *DB {
	return &DB{
		conf: conf,
		mu:   sync.RWMutex{},
	}
}

// Init initialize a new database instance.
func (d *DB) Init() error {
	db, err := gorm.Open(postgres.Open(d.conf.GetDSN()), &gorm.Config{
		CreateBatchSize:        1000,
		SkipDefaultTransaction: true,
	})
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
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
		sqlDB.SetConnMaxLifetime(maxLifetime)
	}

	return nil
}

// Close performs closing the database connection.
func (d *DB) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("SQL instance error: %w", err)
	}

	err = sqlDB.Close()

	if err != nil {
		return fmt.Errorf("db connection closed error: %w", err)
	}

	return nil
}

// Ping performs checking the database connection.
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

	return nil
}

func pingDB(sqlDB *sql.DB) error {
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping db error: %w", err)
	}

	return nil
}

// AddOld (Deprecated) performs adding URL to the DB.
func (d *DB) AddOld(longURL, userID string) (string, error) {
	userIDtoInt, err := strconv.Atoi(userID)

	uri := models.Urls{
		ShortUID:  service.GenUID(d.conf.GetSizeUID()),
		LongURL:   longURL,
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

// Delete performs removing data from shortens URLs by UID.
func (d *DB) Delete(shortURL string) error {
	res := d.db.Where("short_uid = ?", shortURL).
		Delete(&models.Urls{})

	return res.Error
}

// Add performs adding URL to the DB.
func (d *DB) Add(longURL, userID string) (string, error) {
	var uri models.Urls
	res := d.db.Where("long_url = ?", longURL).FirstOrInit(&uri)
	if res.RowsAffected > 0 {
		return d.conf.GetBaseURL() + uri.ShortUID, ErrAlreadyExists
	}

	uri.LongURL = longURL
	uri.ShortUID = service.GenUID(d.conf.GetSizeUID())
	userIDtoInt, err := strconv.Atoi(userID)
	if err == nil {
		uri.UID = uint(userIDtoInt)
	}

	d.db.Create(&uri)

	return d.conf.GetBaseURL() + uri.ShortUID, nil
}

// AddBatch performs adding URL to the DB (batch insert).
func (d *DB) AddBatch(longURL, userID, corID string) string {
	userIDtoInt, err := strconv.Atoi(userID)

	uri := models.Urls{
		ShortUID: service.GenUID(d.conf.GetSizeUID()),
		LongURL:  longURL,
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

// UpdateOld changes short link usage statistics.
func (d *DB) UpdateOld(uid string) {
	var url models.Urls

	res := d.db.Where("short_uid = ?", uid).First(&url)

	if res.Error != nil || res.RowsAffected == 0 {
		return
	}

	d.mu.Lock()
	url.Statistics++
	d.mu.Unlock()

	d.db.Save(&url)
}

// Update changes short link usage statistics.
func (d *DB) Update(uid string) {
	d.db.Model(&models.Urls{}).
		Where("short_uid = ?", uid).
		UpdateColumn("statistics", gorm.Expr("statistics + ?", 1))
}

// GetURL returns the original url by its id.
func (d *DB) GetURL(uid string) (string, error) {
	var url models.Urls

	res := d.db.Where("short_uid = ?", uid).First(&url)

	if res.RowsAffected == 0 {
		return url.LongURL, ErrUIDNotValid
	}

	if url.Removed {
		return url.LongURL, ErrShortURLRemoved
	}

	return url.LongURL, nil
}

// Stat returns short link uses statistics by its short UID.
func (d *DB) Stat(uid string) (Statistics, error) {
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

// Create performs adding new user in DB.
func (d *DB) Create() uint {
	user := models.Users{}
	d.db.Create(&user)

	return user.UID
}

func getUrlsUser(db *gorm.DB, uid uint) []models.Urls {
	var urls []models.Urls

	db.Where("uid = ?", uid).Find(&urls)

	return urls
}

// GetUrlsUser performs getting shorts URLs from DB for current user.
func (d *DB) GetUrlsUser(userID string) ([]struct {
	ShortUID string `json:"short_url"`
	LongURL  string `json:"original_url"`
}, error,
) {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return []struct {
			ShortUID string `json:"short_url"`
			LongURL  string `json:"original_url"`
		}{}, ErrUserIDNotValid
	}

	var urls []struct {
		ShortUID string `json:"short_url"`
		LongURL  string `json:"original_url"`
	}

	d.db.Model(models.Urls{}).
		Where("uid = ?", uid).Find(&urls)

	if len(urls) == 0 {
		return urls, ErrUserUrlsEmpty
	}

	return urls, nil
}

// GetUrlsUserOld (Deprecated) performs getting shorts URLs from DB for current user.
func (d *DB) GetUrlsUserOld(userID string) ([]struct {
	ShortUID string `json:"short_url"`
	LongURL  string `json:"original_url"`
}, error,
) {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		return []struct {
			ShortUID string `json:"short_url"`
			LongURL  string `json:"original_url"`
		}{}, ErrUserIDNotValid
	}

	urls := getUrlsUser(d.db, uint(uid))

	usersURL := make([]struct {
		ShortUID string `json:"short_url"`
		LongURL  string `json:"original_url"`
	}, 0, len(urls))

	if len(urls) == 0 {
		return usersURL, ErrUserUrlsEmpty
	}

	for _, item := range urls {
		usersURL = append(usersURL, struct {
			ShortUID string `json:"short_url"`
			LongURL  string `json:"original_url"`
		}{
			ShortUID: item.ShortUID,
			LongURL:  item.LongURL,
		})
	}

	return usersURL, nil
}

// DelUrls marks as deleted urls added by a specific user.
func (d *DB) DelUrls(userID string, shortsUID ...string) error {
	if len(shortsUID) == 0 || userID == "" {
		return ErrInvalidData
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		return ErrUserIDNotValid
	}

	res := d.db.Model(models.Urls{}).
		Where("short_uid IN ? AND uid = ?", shortsUID, uid).
		Updates(models.Urls{
			Removed: true,
		})

	return res.Error
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

func writeURLBatch(db *gorm.DB, uri models.Urls) int {
	res := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&uri)

	return int(res.RowsAffected)
}
