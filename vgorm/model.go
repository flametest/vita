package vgorm

import (
	"time"

	"gorm.io/gorm"
)

type BaseMysql struct {
	Id        uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Version   uint64    `gorm:"column:version"`
	CreatedAt time.Time `gorm:"<-:create;index;type:TIMESTAMP;default:CURRENT_TIMESTAMP not null;column:created_at"`
	UpdatedAt time.Time `gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP;column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:TIMESTAMP;column:deleted_at"`
}

type BasePostgres struct {
	ID        string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Version   uint64 `gorm:"column:version"`
	CreatedAt time.Time `gorm:"<-:create;index;type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP;not null;column:created_at"`
	UpdatedAt time.Time `gorm:"type:TIMESTAMPTZ;default:CURRENT_TIMESTAMP;not null;column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:TIMESTAMPTZ;column:deleted_at"`
}
