// Package db - This package contains the code for connecting to the mysql datastore
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Database configuration
var (
	dbDriver   = "mysql"
	dbUser     = os.Getenv("MYSQL_USER")
	dbPassword = os.Getenv("MYSQL_PASSWORD")
	dbName     = os.Getenv("MYSQL_DATABASE")
	dbHostName = os.Getenv("MYSQL_HOSTNAME")
	dbPort     = os.Getenv("MYSQL_PORT")
)

// DB - is a connection object which is used to talk to db
var DB *sql.DB

/*
Init initializes the connection to the db and assigns it to an exported DB object, fatal if connection is failed.

	It reads all necessary values that are required to connect to db
*/
func Init() {
	// Open a database connection
	var err error
	DB, err = sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHostName, dbPort, dbName))
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to verify the connection
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
