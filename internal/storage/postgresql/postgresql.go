package postgresql

import "database/sql"

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage/postgresql.New"

	db, err := sql.Open("postgres", "./url-shortner.db")
	
}
