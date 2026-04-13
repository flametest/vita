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
	switch config.Dialect {
	case DialectPostgres:
		return postgres.Open(config.DSN())
	case DialectMySQL:
		return mysql.Open(config.DSN())
	case DialectSQLite3:
		return sqlite.Open(config.DSN())
	default:
		panic("unsupport dialect")
	}
}
