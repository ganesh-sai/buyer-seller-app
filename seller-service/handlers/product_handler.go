package handlers

import (
	"encoding/json"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/models"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/pkg/logging"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/utils"
	"io"
	"net/http"
)

// ProductHandler handles the creation of a product
//
// POST is the only operation supported by this API. The below is json input
//
// Input:
//
//		{
//	 		"sellerId": 1,
//	 		"productName": "Product Name",
//	 		"price": 10.0,
//	 		"quantity": 5
//	 	}
//
// Output:
//
//	{
//		"id": 1
//	}
//
// Returns a saved record id back as JSON response or error if any
func ProductHandler(w http.ResponseWriter, r *http.Request) {
	// Set up panic handler
	defer utils.PanicHandler(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse request body
	var product *models.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		logging.GetLogger().Debugf("%v", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the product
	err = product.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save the product in the database
	err = product.Save()
	if err != nil {
		http.Error(w, "Failed to save product", http.StatusInternalServerError)
		return
	}

	// Respond with the created product ID
	response := struct {
		ID int `json:"id"`
	}{
		ID: product.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
