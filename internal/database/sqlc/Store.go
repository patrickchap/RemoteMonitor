package database

import "database/sql"

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
}

func NewStore(conn *sql.DB) Store {
	return &SQLStore{
		Queries: New(conn),
	}
}
