package service

import (
	"bank/store"
)

type Service struct {
	User
	Tx
	Admin
}

func NewService(store *store.Store) *Service {
	return &Service{
		NewUserService(store),
		NewTxService(store),
		NewAdminService(store),
	}
}

