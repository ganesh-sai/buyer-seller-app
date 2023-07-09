package handlers

import (
	"github.com/ganesh-sai/buyer-seller-app/seller-service/models"
	"net/http"
	"strconv"
)

// SearchProducts will search for products matching this API
// supports based-on-input GET call
//
// Query Parameters: below are list of available query params
//
//	`productName` (optional): Product name for filtering products
//	`desiredQty` (optional): Desired quantity for filtering products
//	`location` (optional): Location for filtering products
//	`minPrice` (optional): Minimum price for filtering products
//	`maxPrice` (optional): Maximum price for filtering products
//	`sortBy` (optional): Field to sort the products (available:  "price", "productName", "sellerId", "productId")
//	`page` (optional): Page number for pagination
//	`perPage` (optional): Number of products per page
//
// Returns:
//
//	[ {
//		  "id": 1,
//		  "sellerId": 1,
//		  "productName": "Product Name",
//		  "price": 10.0,
//		  "quantity": 5
//		},
//		{
//		  "id": 2,
//		  "sellerId": 1,
//		  "productName": "Another Product",
//		  "price": 20.0,
//		  "quantity": 3
//		},
//		... ]
func SearchProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	productName := r.URL.Query().Get("productName")
	desiredQty := r.URL.Query().Get("desiredQty")
	location := r.URL.Query().Get("location")
	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")
	sortBy := r.URL.Query().Get("sortBy")
	page := r.URL.Query().Get("page")
	perPage := r.URL.Query().Get("perPage")
	desiredQuantity, err := strconv.Atoi(desiredQty)
	if err != nil {
		desiredQuantity = 1
	}
	page1, err := strconv.ParseUint(page, 10, 64)
	if err != nil {
		page1 = 1
	}
	perPage1, err := strconv.ParseUint(perPage, 10, 64)
	if err != nil {
		perPage1 = 10
	}
	minimumPrice, err := strconv.ParseFloat(minPrice, 64)
	if err != nil {
		minimumPrice = 0
	}
	maximumPrice, err := strconv.ParseFloat(maxPrice, 64)
	if err != nil {
		maximumPrice = 0
	}
	var productRequest = models.NewProductRequest(productName, desiredQuantity, location, minimumPrice, maximumPrice, sortBy, page1, perPage1)

	resp, err := productRequest.SearchProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}
