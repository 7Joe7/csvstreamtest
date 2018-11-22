package dbconn

import "github.com/jinzhu/gorm"

// Transactor handles DB transaction
type Transactor interface {
	Begin() (*gorm.DB, error)
	Rollback(tr *gorm.DB) error
	Commit(tr *gorm.DB) error
}
