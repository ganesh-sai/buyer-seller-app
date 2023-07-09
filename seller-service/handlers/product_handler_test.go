package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/db"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/models"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/pkg/logging"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	dbUser     = os.Getenv("MYSQL_USER")
	dbPassword = os.Getenv("MYSQL_PASSWORD")
	dbName     = os.Getenv("MYSQL_DATABASE")
	dbHostName = os.Getenv("MYSQL_HOSTNAME")
	dbPort     = os.Getenv("MYSQL_PORT")
)

func TestProductHandler(t *testing.T) {
	err := setupTestDatabase()
	if err != nil {
		t.Fatalf("Unable to run tests")
	}
	product := models.Product{
		SellerID:    1,
		ProductName: "Sample Product",
		Price:       10.0,
		Quantity:    5,
	}
	payload, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Failed to marshal request payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(payload))

	recorder := httptest.NewRecorder()

	ProductHandler(recorder, req)
	res := recorder.Result()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	t.Log(string(body))
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Unexpected response status code: %d", res.StatusCode)
	}

	var response struct {
		ID int `json:"id"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if response.ID == 0 {
		t.Error("Product ID should not be zero")
	}
}

func TestProductHandler_InvalidPayload(t *testing.T) {
	invalidPayload := []byte(`{"invalid": "payload"`)
	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(invalidPayload))
	recorder := httptest.NewRecorder()
	ProductHandler(recorder, req)
	res := recorder.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Unexpected response status code: %d", res.StatusCode)
	}
}

func TestProductHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/product", nil)
	recorder := httptest.NewRecorder()
	ProductHandler(recorder, req)
	res := recorder.Result()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected response status code: %d", res.StatusCode)
	}
}
func TestProductHandler_SaveErrorInvalidSellerId(t *testing.T) {
	product := models.Product{
		SellerID:    1123,
		ProductName: "Sample Product",
		Price:       10.0,
		Quantity:    5,
	}
	payload, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Failed to marshal request payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(payload))
	recorder := httptest.NewRecorder()

	ProductHandler(recorder, req)

	res := recorder.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Unexpected response status code: %d", res.StatusCode)
	}
}

func dropTestDatabase(tables ...string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHostName, dbPort, dbName))
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}
	// DISABLE Foreign Key Checks
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		fmt.Println("Failed to disable foreign key checks:", err)
		return err
	}
	for _, table := range tables {
		stmt := fmt.Sprintln("TRUNCATE TABLE " + table + ";")

		_, err := db.Exec(stmt)
		if err != nil {
			fmt.Printf("Failed to truncate table %s: %v\n", table, err)
		} else {
			fmt.Printf("Successfully truncated table %s\n", table)
		}
	}
	// Enable Foreign key checks
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		fmt.Println("Failed to enable foreign key checks:", err)
	}
	return nil
}

func TestMain(m *testing.M) {
	// Run tests
	logOutput := os.Stdout
	logging.Init(logging.Config{
		Output:   logOutput,
		Prefix:   "test",
		LogLevel: logging.DEBUG,
	})
	db.Init()
	exitCode := m.Run()

	// Clean up
	err := dropTestDatabase("products", "sellers")
	if err != nil {
		fmt.Printf("Failed to drop test database: %v\n", err)
	}
	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

func setupTestDatabase() error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHostName, dbPort, dbName))
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	insertSellerQuery := `
		INSERT INTO sellers (name, location)
		VALUES (?, ?)
	`
	sellers := []struct {
		Name     string
		Location string
	}{
		{Name: "Seller A", Location: "IND"},
		{Name: "Seller B", Location: "US"},
		{Name: "Seller C", Location: "IND"},
		{Name: "Seller D", Location: "UK"},
		{Name: "Seller E", Location: "US"},
	}
	for _, seller := range sellers {
		_, err = db.Exec(insertSellerQuery, seller.Name, seller.Location)
		if err != nil {
			return fmt.Errorf("failed to insert seller: %v", err)
		}
	}

	// Insert products with different sellers
	insertProductQuery := `
		INSERT INTO products (seller_id, product_name, price, quantity)
		VALUES (?, ?, ?, ?)
	`
	productNames := []string{
		"Smartphone", "Laptop", "Tablet", "Smartwatch", "Headphones",
		"Wireless Earbuds", "Gaming Console", "VR Headset", "Fitness Tracker", "Bluetooth Speaker",
		"Drone", "Camera", "Power Bank", "External Hard Drive", "Robot Vacuum",
		"Smart Home Hub", "Smart TV", "Wireless Router", "Wireless Mouse", "E-book Reader",
	}
	for i := 0; i < 30; i++ {
		sellerID := (i % 5) + 1 // Rotate between seller IDs 1 to 5
		productName := productNames[i%len(productNames)]
		price := float64((i + 1) * 10)
		quantity := i + 1

		_, err = db.Exec(insertProductQuery, sellerID, productName, price, quantity)
		if err != nil {
			return fmt.Errorf("failed to insert product: %v", err)
		}
	}
	return nil
}
