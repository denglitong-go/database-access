// With database/sql package, your code opens a database handle
// that represents a connection pool, then executes data access operations with the handle,
// calling a Close method only when needed to free resources, such as
// those held by retrieved rows or a prepared statement.
package basic

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

func ShowOpenDBHanlde() {
	// When opening a database handle, you follow these high-level steps:
	// 1. Locate a driver
	// 2. Open a database handle
	// 3. Confirm a connection

	// Specify connection properties:
	cfg := mysql.Config{
		User:   dbName,
		Passwd: dbPassword,
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	// Get a database handle
	dbHandle, err := sql.Open(dbDriverName, cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	// Release resources
	defer dbHandle.Close()

	// Openning with a Connector
	dbConnector, err := mysql.NewConnector(&cfg)
	// Get a database handle
	dbHandle = sql.OpenDB(dbConnector)

	// Confirming a connection
	if err = dbHandle.Ping(); err != nil {
		log.Fatal(err)
	}

	// Freeing resources
	var artist = ""
	rows, err := dbHandle.Query("SELECT * FROM album WHERE artist = ?", artist)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}
