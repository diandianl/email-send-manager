package dao

import (
	"strings"

	"github.com/google/wire"
	"gorm.io/gorm"

	"email-send-manager/internal/app/config"
	"email-send-manager/internal/app/dao/customer"
	"email-send-manager/internal/app/dao/record"
	"email-send-manager/internal/app/dao/setting"
	"email-send-manager/internal/app/dao/template"
	"email-send-manager/internal/app/dao/util"
) // end

// RepoSet repo injection
var RepoSet = wire.NewSet(
	util.TransSet,
	customer.CustomerSet,
	template.TemplateSet,
	record.RecordSet,
	setting.SettingSet,
) // end

// Define repo type alias
type (
	TransRepo    = util.Trans
	CustomerRepo = customer.CustomerRepo
	TemplateRepo = template.TemplateRepo
	RecordRepo   = record.RecordRepo
	SettingRepo  = setting.SettingRepo
) // end

// Auto migration for given models
func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return db.AutoMigrate(
		new(customer.Customer),
		new(template.Template),
		new(record.Record),
		new(setting.Setting),
	) // end
}
