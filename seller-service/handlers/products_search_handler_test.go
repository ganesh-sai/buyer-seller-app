package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchProductsHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/search/products?productName=SmartPhone&desiredQty=1&location=IND&minPrice=10&maxPrice=100&sortBy=price&page=1&perPage=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	SearchProducts(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	respBody, err := io.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	expectedResp := []byte(`[{"ID":1,"sellerId":1,"productName":"Smartphone","price":10,"quantity":1}]`)
	if !bytes.Equal(respBody, expectedResp) {
		t.Errorf("Unexpected response body. Expected: %s, Got: %s", expectedResp, respBody)
	}
}
