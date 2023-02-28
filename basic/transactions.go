package basic

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Typically, a transaction workflow includes:
// 1. Beginning the transaction.
// 2. Performing a set of database operations.
// 3. If no error occurs, committing the transaction to make database changes.
// 4. If an error occurs, rolling back the transaction to leave the database unchanged.
// 	DB.Begin or Db.BeginTx
// 	Tx.Exec or Tx.ExecContext
// 	Tx.Commit
// 	Tx.Rollback

// CreateOrder creates an order for an album and returns the new order ID.
func CreateOrder(ctx context.Context, albumId, quantity, customerID int) (orderID int64, err error) {
	// Create a helper function for preparing failure results.
	fail := func(err error) (int64, error) {
		return 0, fmt.Errorf("CreateOrder: %v", err)
	}

	// Get a Tx with context for making transaction requests.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	// Defer a rollback in case anything fails.
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			
		}
	}(tx)

	var enough bool

	row := tx.QueryRowContext(ctx,
		"SELECT (quantity >= ?) from album where id = ?", quantity, albumId)
	if err := row.Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			return fail(fmt.Errorf("no such album"))
		}
		return fail(err)
	}
	if !enough {
		return fail(fmt.Errorf("not enough inventory"))
	}

	// Update the album inventory to remove the quantity in the order.
	_, err = tx.ExecContext(ctx,
		"UPDATE album SET quantity = quantity - ? WHERE id = ?",
		quantity, albumId)
	if err != nil {
		return fail(err)
	}

	// Create a enw row in the album_order table.
	result, err := tx.ExecContext(ctx,
		"INSERT INTO album_order (album_id, customer_id, quantity, date) VALUES (?, ?, ?, ?)",
		albumId, customerID, quantity, time.Now())
	if err != nil {
		return fail(err)
	}
	// Get the ID of the order item just created
	orderID, err = result.LastInsertId()
	if err != nil {
		return fail(err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fail(err)
	}

	return orderID, nil
}
