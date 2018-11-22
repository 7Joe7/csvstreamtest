package dbconn

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// mysqlRetryWrapper serves to retry in case of invalid connection error
// this may happen when the connection was interrupted from the other side (e.g. restart of MySQL)
type mysqlRetryWrapper struct {
	*sql.DB
	db gorm.SQLCommon
}

func (w *mysqlRetryWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := w.db.Exec(query, args...)
	if err == mysql.ErrInvalidConn {
		result, err = w.db.Exec(query, args...)
	}
	return result, err
}

func (w *mysqlRetryWrapper) Prepare(query string) (*sql.Stmt, error) {
	stmt, err := w.db.Prepare(query)
	if err == mysql.ErrInvalidConn {
		stmt, err = w.db.Prepare(query)
	}
	return stmt, err
}

func (w *mysqlRetryWrapper) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := w.db.Query(query, args...)
	if err == mysql.ErrInvalidConn {
		rows, err = w.db.Query(query, args...)
	}
	return rows, err
}

func (w *mysqlRetryWrapper) QueryRow(query string, args ...interface{}) *sql.Row {
	return w.db.QueryRow(query, args...)
}
