package models

import (
	"github.com/ganesh-sai/buyer-seller-app/seller-service/db"
)

// Seller represents a seller
type Seller struct {
	ID       int
	Name     string
	Location string
}

// Save saves the seller in the database using a transaction and returns the inserted object and last inserted ID
func (s *Seller) Save() error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO sellers (name, location)
		VALUES (?, ?)
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(s.Name, s.Location)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Retrieve the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	s.ID = int(id)
	return nil
}
