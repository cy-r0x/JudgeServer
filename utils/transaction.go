package utils

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// TxFunc is a function that runs within a transaction
type TxFunc func(*sqlx.Tx) error

// WithTransaction executes a function within a database transaction
// Automatically commits on success or rolls back on error
func WithTransaction(db *sqlx.DB, fn TxFunc) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Printf("Transaction panicked and rolled back: %v", p)
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Transaction rollback failed: %v", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				log.Printf("Transaction commit failed: %v", cmErr)
				err = cmErr
			}
		}
	}()

	err = fn(tx)
	return err
}

// Example usage:
// err := utils.WithTransaction(h.db, func(tx *sqlx.Tx) error {
//     _, err := tx.Exec("INSERT INTO ...")
//     if err != nil {
//         return err
//     }
//     _, err = tx.Exec("UPDATE ...")
//     return err
// })
