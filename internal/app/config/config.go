package config

import (
	"strings"
	"sync"


	"github.com/koding/multiconfig"
)

var (
	// C 全局配置(需要先执行MustLoad，否则拿不到配置)
	C    = new(Config)
	once sync.Once
)

// MustLoad 加载配置
func MustLoad(fpaths ...string) {
	once.Do(func() {
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
			&multiconfig.EnvironmentLoader{Prefix: "ESM", CamelCase: true},
		}

		for _, fpath := range fpaths {
			if strings.HasSuffix(fpath, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "json") {
				loaders = append(loaders, &multiconfig.JSONLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "yaml") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: fpath})
			}
		}

		m := multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}
		m.MustLoad(C)
	})
}

// Config 配置参数
type Config struct {
	RunMode      string `default:"debug"`
	HTTP         HTTP
	Log          Log
	Gorm         Gorm
	Sqlite3      Sqlite3
}

// IsDebugMode 是否是debug模式
func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

// Log 日志配置参数
type Log struct {
	Level         int `default:"4"`
	Format        string `default:"text"`
	Output        string `default:"stdout"`
	OutputFile    string
}

// HTTP http配置参数
type HTTP struct {
	Host               string `default:"127.0.0.1"`
	Port               int `default:"9527"`
	ShutdownTimeout    int `default:"30"`
	MaxReqLoggerLength int `default:"1024"`
	MaxResLoggerLength int
}

// Gorm gorm配置参数
type Gorm struct {
	Debug             bool `default:"false"`
	DBType            string `default:"sqlite3"`
	MaxLifetime       int `default:"7200"`
	MaxOpenConns      int `default:"150"`
	MaxIdleConns      int `default:"50"`
	TablePrefix       string `default:"tb_"`
	EnableAutoMigrate bool `default:"true"`
}

// Sqlite3 sqlite3配置参数
type Sqlite3 struct {
	Path string `default:"sqlite.db"`
}

// DSN 数据库连接串
func (a Sqlite3) DSN() string {
	return a.Path
}
