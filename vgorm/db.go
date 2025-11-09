package vgorm

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDB(config *Config) (*gorm.DB, error) {
	dialector := NewDialector(config)
	db, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		Logger:                                   logger.Default.LogMode(config.LogMode()),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if config.MaxIdleConns != nil {
		sqlDB.SetMaxIdleConns(*config.MaxIdleConns)
	}

	if config.MaxOpenConns != nil {
		sqlDB.SetMaxOpenConns(*config.MaxOpenConns)
	}
	return db, nil
}
