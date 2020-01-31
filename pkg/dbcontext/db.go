// Package dbcontext provides DB transaction support for transactions tha span method calls of multiple
// repositories and services.
package dbcontext

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

// DB represents a DB connection that can be used to run SQL queries.
type DB struct {
	db *dbx.DB
}

// TransactionFunc represents a function that will start a transaction and run the given function.
type TransactionFunc func(ctx context.Context, f func(ctx context.Context) error) error

type contextKey int

const (
	txKey contextKey = iota
)

// New returns a new DB connection that wraps the given dbx.DB instance.
func New(db *dbx.DB) *DB {
	return &DB{db}
}

// DB returns the dbx.DB wrapped by this object.
func (db *DB) DB() *dbx.DB {
	return db.db
}

// With returns a Builder that can be used to build and execute SQL queries.
// With will return the transaction if it is found in the given context.
// Otherwise it will return a DB connection associated with the context.
func (db *DB) With(ctx context.Context) dbx.Builder {
	if tx, ok := ctx.Value(txKey).(*dbx.Tx); ok {
		return tx
	}
	return db.db.WithContext(ctx)
}

// Transactional starts a transaction and calls the given function with a context storing the transaction.
// The transaction associated with the context can be accesse via With().
func (db *DB) Transactional(ctx context.Context, f func(ctx context.Context) error) error {
	return db.db.TransactionalContext(ctx, nil, func(tx *dbx.Tx) error {
		return f(context.WithValue(ctx, txKey, tx))
	})
}

// TransactionHandler returns a middleware that starts a transaction.
// The transaction started is kept in the context and can be accessed via With().
func (db *DB) TransactionHandler() routing.Handler {
	return func(c *routing.Context) error {
		return db.db.TransactionalContext(c.Request.Context(), nil, func(tx *dbx.Tx) error {
			ctx := context.WithValue(c.Request.Context(), txKey, tx)
			c.Request = c.Request.WithContext(ctx)
			return c.Next()
		})
	}
}
