package store

import (
	"bank/model"
	"bank/oops"
	"context"
	"database/sql"
	"errors"
)

type TxStore interface {
	Transaction(context.Context, *model.Tx) error
}

type TxModel struct {
	DB *sql.DB
}

func NewTxStore(db *sql.DB) *TxModel {
	return &TxModel{
		DB: db,
	}
}

func (m *TxModel) Transaction(ctx context.Context, form *model.Tx) error{
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	var userID int
	stmt := `SELECT user_id FROM users WHERE phone = $1;`
	err = tx.QueryRowContext(ctx, stmt, form.Phone).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return oops.ErrNoRecord
		}
		return err
	}

	var enough bool
	stmt = `SELECT balance >= $1 from accounts where user_id = $2;`
	err = tx.QueryRowContext(ctx, stmt, form.Amount, form.SenderID).Scan(&enough)
	if err != nil {
		return err
	}

	if !enough {
		return oops.ErrNotEnough
	}

	stmt = `UPDATE accounts SET balance = balance + $1 WHERE user_id = $2;`
	_, err = tx.ExecContext(ctx, stmt, form.Amount, userID)
	if err != nil {
		return err
	}

	stmt = `UPDATE accounts SET balance = balance - $1 WHERE user_id = $2;`
	_, err = tx.ExecContext(ctx, stmt, form.Amount, form.SenderID)
	if err != nil {
		return err
	}

	var senderAcc, recAcc int
	stmt = `SELECT account_id from accounts where user_id = $1;`
	err = tx.QueryRowContext(ctx, stmt, form.SenderID).Scan(&senderAcc)
	if err != nil {
		return err
	}
	stmt = `SELECT account_id from accounts where user_id = $1;`
	err = tx.QueryRowContext(ctx, stmt, userID).Scan(&recAcc)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO transactions (account_id, receiver_id, amount)
					VALUES ($1, $2, $3);`
	_, err = tx.ExecContext(ctx, stmt, senderAcc, recAcc, form.Amount)
	if err != nil{
		return err
	}

	if err = tx.Commit(); err != nil{
		return err
	}

	return nil
}