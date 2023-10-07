package model

import (
	"bank/validator"
	"time"
)

type Tx struct {
	SenderID int `json:"-"` 
	Phone string `json:"phone,omitempty"`
	Amount int	`json:"amount,omitempty"`
	validator.Validator
}

type LoanForm struct {
	UserID int `json:"-"`
	Amount int	`json:"amount,omitempty"`
	validator.Validator
}

type Loan struct {
	UserInfo
	Amount int `json:"amount,omitempty"`
	Date time.Time `json:"date,omitempty"`
	LoanID int `json:"loan_id,omitempty"`
}

type LoanDecision struct {
	LoanID int `json:"loan_id,omitempty"`
	Decision string `json:"decision,omitempty"`
}


