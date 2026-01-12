// Package sqlsetpgxhelper provides a pgxhelper wrapper
// that integrates with github.com/istovpets/sqlset.
//
// It allows executing queries by their identifiers
// (queryID or setID.queryID) as well as by passing raw SQL text.
package sqlsetpgxhelper

import (
	"context"
	"strings"

	"github.com/istovpets/pgxhelper"
	"github.com/istovpets/sqlset"
)

// DBHelper is a wrapper around pgxhelper.DBHelper that uses sqlset.SQLSet
// to resolve query identifiers into SQL text.
type DBHelper struct {
	*pgxhelper.DBHelper
	sqlSet *sqlset.SQLSet
}

// New creates and returns a new DBHelper instance.
func New(sqlSet *sqlset.SQLSet, opts ...pgxhelper.Option) *DBHelper {
	return &DBHelper{
		DBHelper: pgxhelper.New(opts...),
		sqlSet:   sqlSet,
	}
}

// Get resolves queryRef (queryID, setID.queryID, or raw SQL),
// executes it as a single-row query, and scans the result into dest.
func (d *DBHelper) Get(ctx context.Context, dest any, queryRef string, args ...any) error {
	query, err := d.getQuery(queryRef)
	if err != nil {
		return err
	}
	return d.DBHelper.Get(ctx, dest, query, args...)
}

// Select resolves queryRef (queryID, setID.queryID, or raw SQL),
// executes the query, and scans the resulting rows into dest.
func (d *DBHelper) Select(ctx context.Context, dest any, queryRef string, args ...any) error {
	query, err := d.getQuery(queryRef)
	if err != nil {
		return err
	}
	return d.DBHelper.Select(ctx, dest, query, args...)
}

// Exec resolves queryRef (queryID, setID.queryID, or raw SQL)
// and executes the query, returning the number of affected rows.
func (d *DBHelper) Exec(ctx context.Context, queryRef string, args ...any) (int64, error) {
	query, err := d.getQuery(queryRef)
	if err != nil {
		return 0, err
	}

	return d.DBHelper.Exec(ctx, query, args...)
}

// getQuery resolves queryRef into SQL text.
// If queryRef looks like raw SQL, it is returned as-is.
// Otherwise, it is treated as a query identifier and resolved via sqlSet.
func (d *DBHelper) getQuery(queryRef string) (string, error) {
	q := strings.TrimSpace(queryRef)

	// Heuristic: presence of whitespace indicates raw SQL.
	if strings.ContainsAny(q, " \n\t") {
		return q, nil
	}

	return d.sqlSet.Get(q)
}
