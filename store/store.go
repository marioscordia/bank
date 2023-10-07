package store

import "database/sql"

type Store struct {
	UserStore
	TxStore
	AdminStore
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		NewUserStore(db),
		NewTxStore(db),
		NewAdminStore(db),
	}
}