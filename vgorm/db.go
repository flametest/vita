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

	if err := db.Use(&OptimisticLockPlugin{}); err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	maxIdleConns := DefaultMaxIdleConns
	if config.MaxIdleConns != nil {
		maxIdleConns = *config.MaxIdleConns
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)

	maxOpenConns := DefaultMaxOpenConns
	if config.MaxOpenConns != nil {
		maxOpenConns = *config.MaxOpenConns
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)

	connMaxLife := DefaultConnMaxLife
	if config.ConnMaxLife != nil {
		connMaxLife = *config.ConnMaxLife
	}
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLife) * time.Second)

	return db, nil
}
