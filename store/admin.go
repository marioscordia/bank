package store

import (
	"bank/model"
	"context"
	"database/sql"
	"fmt"
)

type AdminStore interface{
	Accounts() ([]*model.Account, error)
	Loans() ([]*model.Loan, error)
	LoanAccept(context.Context, int) error
	LoanReject(int) error
}

type AdminModel struct {
	DB *sql.DB
}

func NewAdminStore(db *sql.DB) *AdminModel{
	return &AdminModel{
		DB: db,
	}
}

func (m *AdminModel) Accounts() ([]*model.Account, error) {
	accounts := []*model.Account{}
	
	stmt := `SELECT u.first_name, u.last_name, u.phone, a.account_id, a.balance FROM users u
					 JOIN accounts a ON u.user_id = a.user_id
					 WHERE u.is_admin = false;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err		
	}
	defer rows.Close()

	for rows.Next() {
		acc := &model.Account{}
		if err := rows.Scan(&acc.Name, &acc.Surname, &acc.Phone, &acc.ID, &acc.Balance); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (m *AdminModel) Loans() ([]*model.Loan, error){
	loans := []*model.Loan{}
	
	stmt := `SELECT u.first_name, u.last_name, u.user_id, l.amount, l.loan_date, l.loan_id from users u
					 JOIN loans l ON l.user_id = u.user_id
					 WHERE l.status = 0;`
	
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		l := &model.Loan{}
		err = rows.Scan(&l.Name, &l.Surname, &l.ID, &l.Amount, &l.Date, &l.LoanID)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		loans = append(loans, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}

func (m *AdminModel) LoanAccept(ctx context.Context, loanID int) error{
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `UPDATE loans SET status = 1 WHERE loan_id = $1;`
	_, err = tx.ExecContext(ctx, stmt, loanID)
	if err != nil {
		return err
	}

	stmt = `SELECT user_id, amount FROM loans WHERE loan_id = $1;`
	var userID, amount int
	if err := tx.QueryRowContext(ctx, stmt, loanID).Scan(&userID, &amount); err != nil {
		return err
	}

	stmt = `UPDATE accounts SET balance = balance + $1
					WHERE user_id = $2;`
	_, err = tx.ExecContext(ctx, stmt, amount, userID)
	if err != nil {
		return err
	}
	
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *AdminModel) LoanReject(loanID int) error{
	stmt := `UPDATE loans SET status = 2 WHERE loan_id = $1;`

	_, err := m.DB.Exec(stmt, loanID)
	
	return err
}