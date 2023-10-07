package store

import (
	"bank/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func InitializeDB(config *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("dbname=%s user=%s password=%s host=localhost port=%d sslmode=disable",
																	config.DB.Name, config.DB.User, config.DB.Password, config.DB.Port)

	db, err := sql.Open(config.DB.Driver, connectionString)
	if err != nil {
		return nil, err			
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := CreateTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) error{
	stmt := `
						CREATE TABLE IF NOT EXISTS users (
								user_id SERIAL PRIMARY KEY,
								first_name VARCHAR(50) NOT NULL,
								last_name VARCHAR(50) NOT NULL,
								phone VARCHAR(11) UNIQUE NOT NULL,
								is_admin BOOLEAN NOT NULL DEFAULT FALSE,
								password VARCHAR(100) NOT NULL						
						);

						CREATE TABLE IF NOT EXISTS accounts (
								account_id SERIAL PRIMARY KEY,
								user_id INT UNIQUE NOT NULL,
								balance INT NOT NULL DEFAULT 100000,
								FOREIGN KEY (user_id) REFERENCES users(user_id)
						);

						CREATE TABLE IF NOT EXISTS transactions (
								transaction_id SERIAL PRIMARY KEY,
								account_id INT NOT NULL,
								receiver_id INT NOT NULL,
								amount INT NOT NULL,
								transaction_date TIMESTAMPTZ DEFAULT NOW(),
								FOREIGN KEY (account_id) REFERENCES accounts(account_id)
						);
						
						CREATE TABLE IF NOT EXISTS loans (
							loan_id SERIAL PRIMARY KEY,
							user_id INT NOT NULL,
							amount INT NOT NULL,
							status INT DEFAULT 0 CHECK (status >= 0 AND status <= 2),
							loan_date TIMESTAMPTZ DEFAULT NOW(),
							FOREIGN KEY (user_id) REFERENCES users(user_id)
						)
					`
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Println("Here1")
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte("Cheburek"), 12)
	if err != nil {
		return err
	}

	stmt = `
					INSERT INTO users (first_name, last_name, phone, is_admin, password)
					VALUES ($1, $2, $3, $4, $5)
					ON CONFLICT (phone)
					DO NOTHING;
				`

	_, err = db.Exec(stmt, "James", "Bond", "87779991100", true, string(password))
	if err != nil{
		fmt.Println("Here2")
		return err
	}

	stmt = `
					INSERT INTO accounts (user_id)
					VALUES ((SELECT user_id from users where phone = $1))
					ON CONFLICT (user_id)
					DO NOTHING;;
					`

	_, err = db.Exec(stmt, "87779991100")
	return err
}