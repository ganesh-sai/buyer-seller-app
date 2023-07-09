package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/models"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSellerHandler(t *testing.T) {
	seller := models.Seller{
		Name:     "Test Seller",
		Location: "Test Location",
	}
	payload, err := json.Marshal(seller)
	if err != nil {
		t.Fatalf("Failed to marshal request payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/seller", bytes.NewReader(payload))

	recorder := httptest.NewRecorder()

	SellerHandler(recorder, req)

	res := recorder.Result()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Unexpected response status code: %d", res.StatusCode)
	}

	expectedResponse := fmt.Sprintf("Seller created with ID: ")
	if !strings.Contains(string(body), expectedResponse) {
		t.Errorf("Unexpected response body. Expected: %s, Got: %s", expectedResponse, string(body))
	}
}

func TestSellerHandler_InvalidPayload(t *testing.T) {
	invalidPayload := []byte(`{"invalid": "payload"}`)

	req := httptest.NewRequest(http.MethodPost, "/seller", bytes.NewReader(invalidPayload))

	recorder := httptest.NewRecorder()

	SellerHandler(recorder, req)

	res := recorder.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Unexpected response status code: %d", res.StatusCode)
	}
}

func TestSellerHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/seller", nil)

	recorder := httptest.NewRecorder()

	SellerHandler(recorder, req)

	res := recorder.Result()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected response status code: %d", res.StatusCode)
	}
}
