package vgorm

import (
	"fmt"

	"gorm.io/gorm/logger"
)

const (
	DefaultMaxIdleConns = 10
	DefaultMaxOpenConns = 25
	DefaultConnMaxLife  = 300 // seconds
)

type Config struct {
	Dialect      Dialect `json:"dialect" yaml:"dialect"`
	Host         string  `json:"host" yaml:"host"`
	Port         string  `json:"port" yaml:"port"`
	Database     string  `json:"database" yaml:"database"`
	Schema       string  `json:"schema" yaml:"schema"`
	Username     string  `json:"username" yaml:"username"`
	Password     string  `json:"password" yaml:"password"`
	Debug        bool    `json:"debug" yaml:"debug"`
	MaxIdleConns *int    `json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns *int    `json:"maxOpenConns" yaml:"maxOpenConns"`
	ConnMaxLife  *int    `json:"connMaxLife" yaml:"connMaxLife"`
}

func (d *Config) DSN() string {
	switch d.Dialect {
	case DialectMySQL:
		return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			d.Username, d.Password, d.Host, d.Database)
	case DialectPostgres:
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=disable",
			d.Host, d.Port, d.Username, d.Password, d.Database, d.Schema)
	case DialectSQLite3:
		return fmt.Sprintf("file:%s?cache=shared&mode=memory", d.Database)
	default:
		panic(fmt.Sprintf("unknown dialect: %s", d.Dialect))
	}
}

func (d *Config) LogMode() logger.LogLevel {
	if d.Debug {
		return logger.Info
	}
	return logger.Warn
}
