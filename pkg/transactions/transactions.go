package transactions

import (
	"database/sql"
	"log"
)

// Transaction is an interface that models the standard transaction in
// `database/sql`.
// To ensure `TxFn` funcs cannot commit or rollback a transaction (which is
// handled by `WithTransaction`), those methods are not included here.
type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// TxFn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(Transaction) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func WithTransaction(db *sql.DB, fn TxFn) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and re-panic
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("could not rollback %v", rbErr)
			}
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("could not rollback %v", rbErr)
			}
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
