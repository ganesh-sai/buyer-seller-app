package main

import (
	"context"
	"flag"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/db"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/handlers"
	"github.com/ganesh-sai/buyer-seller-app/seller-service/pkg/logging"
	"io"
	"log/syslog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	ServiceName = "seller-service"
)

func main() {
	var (
		serviceName   = flag.String("service.name", ServiceName, "Name of service")
		sysLogAddress = flag.String("syslog.address", "localhost:514", "default location for the syslogger")
	)
	flag.Parse()
	if *serviceName == "" {
		serviceName = &ServiceName
	}

	var logOutput io.Writer
	if os.Getenv("ENABLE_DEV_MODE") == "true" {
		logOutput = os.Stdout
	} else {
		sysLogger, err := syslog.Dial("udp", *sysLogAddress, syslog.LOG_EMERG|syslog.LOG_LOCAL6, *serviceName)
		if err != nil {
			panic("failed to connect to syslog! application will now exit")
		}
		defer sysLogger.Close()
		logOutput = sysLogger
	}
	// Initialize the logger
	logging.Init(logging.Config{
		Output:   logOutput,
		Prefix:   *serviceName,
		LogLevel: logging.DEBUG,
	})
	logger := logging.GetLogger()

	// Initialize the datastore
	db.Init()

	// setup routes
	http.HandleFunc("/api/v1/product", handlers.ProductHandler)
	http.HandleFunc("/api/v1/product/search", handlers.SearchProducts)
	http.HandleFunc("/api/v1/seller/", handlers.SellerHandler)

	server := http.Server{Addr: ":8080"}
	logger.Debug("Starting Application")
	go func() {
		// Start Server
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("Server listen failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		logger.Error("Server shutdown failed: %v", err)
	}
	logger.Info("Server shutdown complete")
}
