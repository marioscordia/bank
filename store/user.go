package store

import (
	"bank/model"
	"bank/oops"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface{
	CreateUser(*model.SignUp) error
	Authenticate(*model.SignIn) (*model.User, error)
	Account(int) (*model.Account, error)
	Request(*model.LoanForm) error
}

func NewUserStore(db *sql.DB) *UserModel {
	return &UserModel{
		DB: db,
	}
}

type UserModel struct{
	DB *sql.DB
}

func (m *UserModel) CreateUser(form *model.SignUp) error {
	stmt := `
					INSERT INTO users (first_name, last_name, phone, password)
					VALUES ($1, $2, $3, $4);`

	_, err := m.DB.Exec(stmt, form.Name, form.Surname, form.Phone, form.Password)
	if err!=nil {
		Err, ok := err.(*pq.Error)
		if ok && Err.Code == "23505"{
			return oops.ErrDuplicateEmail
		}
		return err
	}

	stmt = `
	INSERT INTO accounts (user_id)
	VALUES ((SELECT user_id from users where phone = $1))
	ON CONFLICT (user_id)
	DO NOTHING;;
	`

	_, err = m.DB.Exec(stmt, form.Phone)

	return err
}

func (m *UserModel) Authenticate(form *model.SignIn) (*model.User, error) {
	user := model.User{}
	var hashedPassword string

	stmt := `SELECT user_id, is_admin, password FROM users
					 WHERE phone = $1;`

	err := m.DB.QueryRow(stmt, form.Phone).Scan(&user.ID, &user.IsAdmin, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, oops.ErrInvalidCredentials
		}
		return nil, err
	} 

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(form.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, oops.ErrInvalidCredentials
		}
		return nil, err
	} 
	
	return &user, nil
}

func (m *UserModel) Account(userID int) (*model.Account, error){
	acc := &model.Account{}

	stmt := `SELECT u.first_name, u.last_name, u.phone, a.account_id, a.balance FROM users u
					 JOIN accounts a ON u.user_id = a.user_id
					 WHERE u.user_id = $1;
					`
	err := m.DB.QueryRow(stmt, userID).Scan(&acc.Name, &acc.Surname, &acc.Phone, &acc.ID, &acc.Balance)
	if err != nil {
		// Err, ok := err.(*pq.Error)
		// if ok && Err.Code == "P0002"{
		// 	return nil, oops.ErrNoRecord
		// }
		if errors.Is(err, sql.ErrNoRows){
			return nil, oops.ErrNoRecord
		}
		return nil, err
	}

	stmt = `SELECT receiver_id, amount FROM transactions
					WHERE account_id = $1;`

	rows, err := m.DB.Query(stmt, acc.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tx := &model.Transaction{}
		if err := rows.Scan(&tx.ReceiverID, &tx.Amount); err != nil {
			return nil, err
		}
		acc.Transactions = append(acc.Transactions, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return acc, nil
}

func (m *UserModel) Request(req *model.LoanForm) error{
	var enough bool
	stmt := `SELECT balance >= 5000 FROM accounts WHERE user_id = $1;`
	if err := m.DB.QueryRow(stmt, req.UserID).Scan(&enough); err != nil {
		return err
	}

	if enough{
		return oops.ErrNotAllowed
	}

	stmt = `INSERT INTO loans (user_id, amount)
					VALUES ($1, $2);`
	_, err := m.DB.Exec(stmt, req.UserID, req.Amount)
	
	return err
}