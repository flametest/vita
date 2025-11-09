package vgorm

import (
	"fmt"

	"gorm.io/gorm/logger"
)

type Config struct {
	Dialect      Dialect `json:"dialect" yaml:"dialect"`
	Host         string  `json:"host" yaml:"host"`
	Port         string  `json:"port" yaml:"port"`
	Database     string  `json:"database" yaml:"database"`
	Username     string  `json:"username" yaml:"username"`
	Password     string  `json:"password" yaml:"password"`
	Debug        bool    `json:"debug" yaml:"debug"`
	MaxIdleConns *int    `json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns *int    `json:"maxOpenConns" yaml:"maxOpenConns"`
}

func (d *Config) DSN() string {
	if d.Dialect == DialectMySQL {
		return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			d.Username, d.Password, d.Host, d.Database)
	} else if d.Dialect == DialectPostgres {
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
			d.Host, d.Port, d.Username, d.Database)
	} else if d.Dialect == DialectSQLite3 {
		return fmt.Sprintf("file:%s?cache=shared&mode=memory", d.Database)
	} else {
		panic(fmt.Sprintf("unknown dialect: %s", d.Dialect))
	}
}

func (d *Config) LogMode() logger.LogLevel {
	if d.Debug {
		return logger.Info
	}
	return logger.Warn
}
