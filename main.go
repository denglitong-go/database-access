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
	// root:12345678@tcp(127.0.0.1:3306)/recordings?allowNativePasswords=false&checkConnLiveness=false&maxAllowedPacket=0
	log.Println("config data source name", cfg.FormatDSN())
	db, err = sql.Open(dbDriverName, cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Closed!")
	}()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Albums found: %v\n", albums)
}

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

// albumsByArtist queries for albums that have the specific artist name
func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}
