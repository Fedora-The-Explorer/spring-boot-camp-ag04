package bootstrap

import (
	"database/sql"
	"elProfessor/cmd/config"
)

// Sqlite bootstraps the db connection.
func Sqlite() *sql.DB {
	db, err := sql.Open("sqlite3", config.Cfg.SqliteDatabase)
	if err != nil {
		panic(err)
	}

	return db
}
