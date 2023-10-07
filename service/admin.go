package service

import (
	"bank/model"
	"bank/store"
	"context"
)

type Admin interface{
	Accounts() ([]*model.Account, error)
	Loans() ([]*model.Loan, error)
	LoanApproval(context.Context, *model.LoanDecision) error
}


type AdminService struct {
	admin store.AdminStore
}

func NewAdminService(store *store.Store) *AdminService{
	return &AdminService{
		admin: store.AdminStore,
	}
}

func (s *AdminService) Accounts() ([]*model.Account, error) {
	return s.admin.Accounts()
}

func (s *AdminService) Loans() ([]*model.Loan, error) {
	return s.admin.Loans()
}		

func (s *AdminService) LoanApproval(ctx context.Context, form *model.LoanDecision) error{
	if form.Decision == "approve"{
		return s.admin.LoanAccept(ctx, form.LoanID)
	}

	return s.admin.LoanReject(form.LoanID)
}
