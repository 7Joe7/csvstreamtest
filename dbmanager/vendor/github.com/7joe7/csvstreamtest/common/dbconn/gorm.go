package dbconn

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// GORMOpen build DB URL, open a connection and set our default expected behaviour before returning the conn
func GORMOpen(username, password, host, database string, port int) (*gorm.DB, error) {
	// initialize DB connectors
	URL := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
	)
	dbSQL, err := sql.Open("mysql", URL)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to database")
	}
	DBConn, err := gorm.Open("mysql", &mysqlRetryWrapper{DB: dbSQL, db: dbSQL})
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize mysql connection")
	}
	DBConn.SetLogger(log.New(ioutil.Discard, "", 0))

	return DBConn, nil
}

// GORMTransactor is a tool to handle db transaction using gorm
type GORMTransactor struct {
	log *zerolog.Logger
	db  *gorm.DB
}

// NewGORMTransactor is GORMTransactor constructor
func NewGORMTransactor(log *zerolog.Logger, db *gorm.DB) *GORMTransactor {
	return &GORMTransactor{
		log: log,
		db:  db,
	}
}

// Begin performs gorm begin function
func (tra *GORMTransactor) Begin() (*gorm.DB, error) {
	db := tra.db.Begin()
	return db, db.Error
}

// Rollback performs gorm rollback function
func (tra *GORMTransactor) Rollback(tr *gorm.DB) error {
	err := tr.Rollback().Error
	if err != nil {
		tra.log.Error().Err(err).Msg("could not rollback transaction !")
	}
	return err
}

// Commit performs gorm commit function
func (tra *GORMTransactor) Commit(tr *gorm.DB) error {
	err := tr.Commit().Error
	if err != nil {
		tra.log.Error().Err(err).Msg("could not commit transaction !")
	}
	return err
}
