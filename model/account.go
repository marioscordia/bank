package model

type Account struct {
	UserInfo
	Balance int `json:"balance,omitempty"`
	Transactions []*Transaction `json:"transactions,omitempty"`
}

type Transaction struct {
	Amount int `json:"amount,omitempty"`
	ReceiverID int `json:"receiver_id,omitempty"`
}