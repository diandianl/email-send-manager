package app

import (
	"os"
	"path/filepath"

	"gorm.io/gorm"

	"email-send-manager/internal/app/config"
	"email-send-manager/internal/app/dao"
	"email-send-manager/pkg/errors"
	"email-send-manager/pkg/gormx"
)

// InitGormDB 初始化gorm存储
func InitGormDB() (*gorm.DB, func(), error) {
	cfg := config.C.Gorm
	db, err := NewGormDB()
	if err != nil {
		return nil, nil, err
	}

	cleanFunc := func() {}

	if cfg.EnableAutoMigrate {
		err = dao.AutoMigrate(db)
		if err != nil {
			return nil, cleanFunc, err
		}
	}

	return db, cleanFunc, nil
}

// NewGormDB 创建DB实例
func NewGormDB() (*gorm.DB, error) {
	cfg := config.C
	var dsn string
	switch cfg.Gorm.DBType {
	case "sqlite3":
		dsn = cfg.Sqlite3.DSN()
		_ = os.MkdirAll(filepath.Dir(dsn), 0777)
	default:
		return nil, errors.Errorf("unsupported db driver '%s'", cfg.Gorm.DBType)
	}

	return gormx.New(&gormx.Config{
		Debug:        cfg.Gorm.Debug,
		DBType:       cfg.Gorm.DBType,
		DSN:          dsn,
		MaxIdleConns: cfg.Gorm.MaxIdleConns,
		MaxLifetime:  cfg.Gorm.MaxLifetime,
		MaxOpenConns: cfg.Gorm.MaxOpenConns,
		TablePrefix:  cfg.Gorm.TablePrefix,
	})
}
