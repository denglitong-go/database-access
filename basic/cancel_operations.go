package basic

import (
	"context"
	"log"
	"time"
)

// You can manage in-progress operations by using Go context.Context.
// A Context is a standard Go data value that can report whether the overall operation
// it represents has been canceled and is no longer needed.
// By passing a context.Context across function calls and services in your application,
// those can stop working early and return an error when their processing in longer needed.
// Ref: https://blog.golang.org/context

// Cancel operations, for example:
// End long-running operations, including database operations that are taking too long to complete
// Propagate cancellation requests from elsewhere, such as when a client closes a connection

// QueryWithTimeout cancelling database operations after a timeout
func QueryWithTimeout(ctx context.Context) {
	// Create a Context with a timeout.
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Pass the timeout Context with a query.
	rows, err := db.QueryContext(queryCtx, "SELECT * FROM album")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Handle returned rows.
}
