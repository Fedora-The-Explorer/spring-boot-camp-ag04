package sqlite

import (
	"context"
	"database/sql"
)

// DatabaseExecutor implements methods for handling db data
type DatabaseExecutor interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}
