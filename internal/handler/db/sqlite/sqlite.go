package sqlite

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
)

//go:embed schema.sql
var schema []byte

func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open the sqlite database: %w", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the sqlite database: %w", err)
	}

	return db, nil
}
