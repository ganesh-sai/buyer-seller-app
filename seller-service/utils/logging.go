// Package utils - contains utility methods
package utils

import (
	"github.com/ganesh-sai/buyer-seller-app/seller-service/pkg/logging"
	"net/http"
)

// PanicHandler handles panics and logs them
func PanicHandler(w http.ResponseWriter) {
	if err := recover(); err != nil {
		logging.GetLogger().Error("Panic: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
