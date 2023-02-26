// This tutorial introduce the basics of accessing a relational database with Go
// and its database/sql package in the standard library.
// The database/sql package includes types and functions for connecting to databases,
// execution transactions, canceling an operation in progress, and more.
// For more details on using the database/sql package, see https://go.dev/doc/database.
// Sections:
// 	1.Set up a database.
// 	2.Import the database driver.
// 	3.Get a database handle and connect.
// 	4.Query for multiple rows.
// 	5.Query for a single row.
//  6.Add data.
package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var (
	dbName       = "recordings"
	dbAddr       = "127.0.0.1:3306"
	dbUser       = "DB_USER"
	dbPassword   = "DB_PASSWORD"
	dbDriverName = "mysql"
	db           *sql.DB
)

// main
// DB_USER=root DB_PASSWORD=12345678 go run .
func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv(dbUser),
		Passwd: os.Getenv(dbPassword),
		Net:    "tcp",
		Addr:   dbAddr,
		DBName: dbName,
	}
	// Get a database handle
	var err error
	db, err = sql.Open(dbDriverName, cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Closed!")
	}()
}
