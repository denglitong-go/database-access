// Package basic
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
package basic

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

func OpenDB() {
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
}

func PingDB() {
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")
}

func CloseDB() {
	if db == nil {
		log.Println("db is nil.")
		return
	}
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
	log.Println("Closed!")
}

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

// AlbumsByArtist queries for albums that have the specific artist name
func AlbumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AlbumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// AlbumById queries for the album with the specific ID.
func AlbumById(id int64) (Album, error) {
	var alb Album
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("AlbumById %d: no such album", id)
		}
		return alb, fmt.Errorf("AlbumById %d: %v", id, err)
	}
	return alb, nil
}

// AddAlbum adds the specific album to the database,
// returning the album ID of the new entry
func AddAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)",
		alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("AddAlbum: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddAlbum: %v", err)
	}

	return id, nil
}

// UpdateAlbum updates the specific album to the database,
// returning the affected rows of the update entry
func UpdateAlbum(alb Album) (int64, error) {
	result, err := db.Exec("UPDATE album SET title = ?, artist = ?, price = ? WHERE id = ?",
		alb.Title, alb.Artist, alb.Price, alb.ID)
	if err != nil {
		return 0, fmt.Errorf("UpdateAlbum: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("UpdateAlbum: %v", err)
	}

	return rowsAffected, nil
}
