package models

import (
	"encoding/json"
	"errors"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/db"
)

// ProductRequest represents the fields that are used to filter the records
type ProductRequest struct {
	ProductName string  `json:"productName"`
	DesiredQty  int     `json:"desiredQty"`
	Location    string  `json:"location"`
	MinPrice    float64 `json:"minPrice"`
	MaxPrice    float64 `json:"maxPrice"`
	SortBy      string  `json:"sortBy"`
	Page        uint64  `json:"page"`
	PerPage     uint64  `json:"perPage"`
}

// NewProductRequest Creates a new ProductRequest
func NewProductRequest(productName string, desiredQty int, location string, minPrice float64, maxPrice float64, sortBy string, page uint64, perPage uint64) *ProductRequest {
	return &ProductRequest{ProductName: productName, DesiredQty: desiredQty, Location: location, MinPrice: minPrice, MaxPrice: maxPrice, SortBy: sortBy, Page: page, PerPage: perPage}
}

// SearchProducts - This frames the sql queries based on the fields supplied to ProductRequest
// It filters the matching records and returns them and error if any
//
// Returns:
//
//	[]byte, error
func (p *ProductRequest) SearchProducts() ([]byte, error) {
	if p.MinPrice > p.MaxPrice {
		return nil, errors.New("Minimum price cannot be greater than maximum price")
	}

	offset := (p.Page - 1) * p.PerPage

	var query = "SELECT p.* FROM products AS p INNER JOIN sellers AS s ON p.seller_id = s.id WHERE 1=1"
	var args []interface{}

	if p.ProductName != "" {
		query += " AND p.product_name LIKE ? "
		args = append(args, "%"+p.ProductName+"%")
	}
	if p.DesiredQty > 0 {
		query += " AND p.quantity >= ? "
		args = append(args, p.DesiredQty)
	}

	if p.Location != "" {
		query += " AND s.location LIKE ? "
		args = append(args, "%"+p.Location+"%")
	}
	if p.MinPrice > 0 {
		query += " AND p.price >= ?"
		args = append(args, p.MinPrice)
	}
	if p.MaxPrice > 0 {
		query += " AND p.price <= ?"
		args = append(args, p.MaxPrice)
	}

	// Sort by the specified field
	switch p.SortBy {
	case "price":
		query += " ORDER BY p.price"
	case "productName":
		query += " ORDER BY p.product_name"
	case "sellerId":
		query += " ORDER BY p.seller_id"
	case "productId":
		query += " ORDER BY p.id"
	}

	// Add pagination to the query
	query += " LIMIT ? OFFSET ?"
	args = append(args, p.PerPage, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the result set and create a list of products
	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.SellerID, &product.ProductName, &product.Price, &product.Quantity)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	// Convert the products to JSON
	response, err := json.Marshal(&products)
	if err != nil {
		return nil, err
	}
	return response, nil
}
