package vgorm

import "gorm.io/gorm"

type Tx interface {
	Begin() Tx
	Commit() Tx
	Rollback() Tx
	Error() error
}

type txImpl struct {
	db *gorm.DB
}

func NewTxImpl(db *gorm.DB) Tx {
	return &txImpl{db: db}
}
func (t *txImpl) Begin() Tx {
	t.db = t.db.Begin()
	return t
}

func (t *txImpl) Commit() Tx {
	t.db = t.db.Commit()
	return t
}

func (t *txImpl) Rollback() Tx {
	t.db = t.db.Rollback()
	return t
}

func (t *txImpl) Error() error {
	return t.db.Error
}
