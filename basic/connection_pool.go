// Package basic
// The sql.DB database handle is safe for concurrent use by multiple goroutines
// (meaning the handle is what other languages might call "thread-safe").
// Some other database access libraries are based on connections that can only be used
// for one operation at a time. To bridge that gap, each sql.DB manages a pool
// of active connections to the underlying database, creating new ones as needed for parallelism.
//
// The connection pool is suitable for most data access needs.
// When you call an sql.DB Query or Exec method, the sql.DB implementation
// retrieves an available connection from the pool or, if needed, creates one.
// The package returns the connection to the pool when it's no longer needed.
// This supports a high level of parallelism for database access.
package basic

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

func ShowConnectionPoolProperties() {
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

	log.Println("get connection pool properties:", db.Stats())

	// Setting the maximum number of open connections
	// Pass this limit, new database operations will wait for an existing operation to finish,
	// at which time sql.DB will create another connection.
	// By default, sql.DB creates a new connection any time all the existing connections are
	// in use when a connection is needed.
	// Keep in mind that setting a limit makes database usage similar to acquiring a lock
	// or semaphore, with the result that your application can deadlock
	// waiting for a new database connection.
	db.SetMaxOpenConns(10)

	// Setting the maximum number of idle connections
	// When an SQL operation finishes on a given database connection, it's not typically
	// shut down immediately: the application may need again soon, and keeping the open
	// connection around avoids having to reconnect to the database for the next operation.
	// By default, sql.DB keeps 2 idle connections at any given moment.
	// Raising the limit can avoid frequent reconnects in programms with significant parallelism.
	db.SetMaxIdleConns(3)

	// Setting the maximum amount a time a connection can be idle
	// This causes the sql.DB to close connections that have been idle for longer than
	// the given duration.
	// By default, when an idle connection is added to the connection pool, it remains there
	// until it's needed again. When using DB.SetMaxIdleConns to increase the number of
	// allowed idle connections during bursts of parallel activity, also using
	// DB.SetConnMaxIdleTime can arrange to release those connections later when the system is quiet.
	db.SetConnMaxIdleTime(30 * time.Second)

	// Setting the maximum lifetime of connections
	// By default, a connection can be used and reused for an arbitrary long amount of time.
	// In some systems, such as those using a load-balanced database server, it can be helpful
	// to ensure that the application never uses a particular connection
	// for too long without reconnecting.
	db.SetConnMaxLifetime(300 * time.Second)

	// Using dedicated connections
	// The most common example is transactions, which typically start with a BEGIN command,
	// end with a COMMIT or ROLLBACK command, and include all the commands issued
	// on the connection between those commands in the overall transaction. For this use case,
	// use the sql package's transaction support.
	// For other use cases where a sequence of individual operations
	// must all execute on the same connection, the sql package provides dedicated connection.
	// DB.Conn obtains a dedicated connection, an sql.Conn.
	// The sql.Conn has method BeginTx, ExecContext, PingContext, PrepareContext, QueryContext,
	// and QueryRowContext that behave like the equivalent methods on DB but
	// only use the dedicated connection.
	// When finished with the dedicated connection, your code must release it using Conn.Close.
}
