package vgorm

import (
	"fmt"
	"runtime/debug"

	log "github.com/flametest/vita/vlog"
	"gorm.io/gorm"
)

type BaseRepo interface {
	DoInTx(fn TxnFunc) (err error)
}

type baseRepoImpl struct {
	db *gorm.DB
}

func NewBaseRepo(db *gorm.DB) BaseRepo {
	return &baseRepoImpl{db: db}
}

func (r *baseRepoImpl) DoInTx(fn TxnFunc) (err error) {
	var tx = NewTxImpl(r.db).Begin()
	if tx.Error() != nil {
		return tx.Error()
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("recovery from panic: %s", debug.Stack())
			switch x := r.(type) {
			case string:
				err = fmt.Errorf(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("unknown panic: %+v", x)
			}
		}
		if err != nil {
			if e := tx.Rollback().Error(); e != nil {
				err = fmt.Errorf("tx rollback err '%s', caused by other err '%s'", e.Error(), err.Error())
			}
		} else {
			err = tx.Commit().Error()
		}
	}()
	return fn(tx)
}
