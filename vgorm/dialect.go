package vgorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dialect string

const (
	DialectPostgres = "postgres"
	DialectMySQL    = "mysql"
	DialectSQLite3  = "sqlite3"
)

func (d Dialect) String() string {
	return string(d)
}

func NewDialector(config *Config) gorm.Dialector {
	if config.Dialect == DialectPostgres {
		return postgres.New(postgres.Config{})
	} else if config.Dialect == DialectMySQL {
		return mysql.Open(config.DSN())
	} else if config.Dialect == DialectSQLite3 {
		return sqlite.Open(config.DSN())
	} else {
		panic("unsupport dialect")
	}
}
