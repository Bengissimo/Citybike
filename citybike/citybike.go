//Package citybike creates sqlite database
// and includes database related methods and functions
package citybike

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Citybike struct {
	db *sql.DB
}

func New(path string) (*Citybike, error) {
	sqldb, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db := &Citybike{
		db: sqldb,
	}
	return db, nil
}

