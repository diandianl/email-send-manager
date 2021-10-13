package gormx

import (
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"email-send-manager/pkg/errors"
)

// Config GORM Config
type Config struct {
	Debug        bool
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TablePrefix  string
}

// New Create gorm.DB instance
func New(c *Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(c.DBType) {
	case "sqlite3":
		dialector = sqlite.Open(c.DSN)
	default:
		return nil, errors.Errorf("unsupported db driver '%s'", c.DBType)
	}

	gconfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.TablePrefix,
			SingularTable: true,
		},
	}

	db, err := gorm.Open(dialector, gconfig)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)

	return db, nil
}
