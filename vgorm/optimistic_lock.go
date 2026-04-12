package vgorm

import (
	"reflect"

	"github.com/flametest/vita/verrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OptimisticLockPlugin struct{}

func (p *OptimisticLockPlugin) Name() string {
	return "optimistic_lock"
}

func (p *OptimisticLockPlugin) Initialize(db *gorm.DB) error {
	if err := db.Callback().Update().Before("gorm:update").Register("optimistic_lock:before_update", p.beforeUpdate); err != nil {
		return err
	}
	return db.Callback().Update().After("gorm:update").Register("optimistic_lock:after_update", p.afterUpdate)
}

func (p *OptimisticLockPlugin) beforeUpdate(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}

	versionField, exists := db.Statement.Schema.FieldsByDBName["version"]
	if !exists {
		return
	}

	switch db.Statement.ReflectValue.Kind() {
	case reflect.Struct:
		val, isZero := versionField.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
		if isZero {
			return
		}
		currentVersion := val.(uint64)
		db.Statement.AddClause(clause.Where{
			Exprs: []clause.Expression{
				clause.Eq{
					Column: clause.Column{Table: db.Statement.Table, Name: "version"},
					Value:  currentVersion,
				},
			},
		})
		db.Statement.SetColumn("version", currentVersion+1)
	}
}

func (p *OptimisticLockPlugin) afterUpdate(db *gorm.DB) {
	if db.Error != nil || db.Statement.Schema == nil {
		return
	}
	if _, exists := db.Statement.Schema.FieldsByDBName["version"]; !exists {
		return
	}
	if db.RowsAffected == 0 {
		db.Error = verrors.ErrOptimisticLock
	}
}
