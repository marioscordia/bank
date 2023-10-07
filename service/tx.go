package service

import (
	"bank/model"
	"bank/oops"
	"bank/store"
	"bank/validator"
	"context"
	"errors"
)

type Tx interface{
	Transact(context.Context, *model.Tx) error
}

type TxService struct {
	tx store.TxStore
}

func NewTxService(store *store.Store) *TxService {
	return &TxService{
		tx: store.TxStore,
	}
}

func (s *TxService) Transact(ctx context.Context, form *model.Tx) error{
	form.CheckField(validator.NotBlank(form.Phone), "phone", "This field can not be blank")
	form.CheckField(validator.CheckPhone(form.Phone), "phone", "Phone number must be written in correct format")
	form.CheckField(validator.NumCheck(form.Amount), "amount", "Amount must be between 100 and 20000")

	if !form.Valid(){
		return oops.ErrFormInvalid
	}

	if err := s.tx.Transaction(ctx, form); err != nil {
		if errors.Is(err, oops.ErrNoRecord){
			form.AddFieldError("phone", "No user with this phone number")
		}else if errors.Is(err, oops.ErrNotEnough){
			form.AddFieldError("amount", "Not enough money on your account")
		}

		return err
	}

	return nil
}

