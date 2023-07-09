// Package handlers - This package contains the code for handling the code for
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ganesh-sai/buyer-seller-app/seller-service/models"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/utils"
)

// SellerHandler handles the creation of a seller,
//
// POST is the only operation supported by this API. The below is json input
//
// Input:
//
//	{
//	  "name": "Seller Name",
//	  "location": "Seller Location"
//	}
//
// Output:
//
//	Seller created with ID: 1
//
// Returns a saved record id back as JSON response or error if any
func SellerHandler(w http.ResponseWriter, r *http.Request) {
	defer utils.PanicHandler(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var seller models.Seller
	err = json.Unmarshal(body, &seller)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if seller.Name == "" || seller.Location == "" {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Name and location cannot be empty", http.StatusBadRequest)
		return
	}

	err = seller.Save()
	if err != nil {
		http.Error(w, "Failed to save seller", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Seller created with ID: %d", seller.ID)
}
