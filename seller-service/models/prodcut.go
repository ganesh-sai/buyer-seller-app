// Package models - responsible for storing retrieving the data from mysql
package models

import (
	"database/sql"
	"errors"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/db"
)

// Product represents a product
type Product struct {
	ID          int     `json:"ID,omitempty"`
	SellerID    int     `json:"sellerId"`
	ProductName string  `json:"productName"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

// Save saves the product in the database
func (p *Product) Save() error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT  INTO products (seller_id, product_name, price, quantity) VALUES (?,?,?,?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(p.SellerID, p.ProductName, p.Price, p.Quantity)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// Validate validates the product
// Checks if the sellerId is present or not in DB..
// If not, it will through an error
func (p *Product) Validate() error {
	// Check if SellerID is valid (existing seller)
	_, err := GetSellerByID(p.SellerID)
	if err != nil {
		return errors.New("invalid SellerID")
	}

	return nil
}

// GetSellerByID retrieves a seller by ID from the database
// Args:
//
//	id int: seller id
//
// Returns:
//
//	Seller: Seller Object
//	error: root cause of error
func GetSellerByID(id int) (Seller, error) {
	var seller Seller

	row := db.DB.QueryRow(`
		SELECT id, name, location
		FROM sellers
		WHERE id = ?
	`, id)

	err := row.Scan(&seller.ID, &seller.Name, &seller.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return seller, errors.New("seller not found")
		}
		return seller, err
	}

	return seller, nil
}
