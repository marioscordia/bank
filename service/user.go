package service

import (
	"bank/model"
	"bank/oops"
	"bank/store"
	"bank/validator"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User interface{
	Register(*model.SignUp) error
	Authenticate(*model.SignIn)(*model.User, error)
	Account(int) (*model.Account, error)
	Request(*model.LoanForm) error
}

type UserService struct {
	user store.UserStore
}

func NewUserService(store *store.Store) *UserService {
	return &UserService{
		user: store.UserStore,
	}
}

func (s *UserService) Register(form *model.SignUp) error{
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Surname), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Phone), "phone", "This field cannot be blank")
	form.CheckField(validator.CheckPhone(form.Phone), "phone", "This field must be a valid phone number")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.Confirm), "confirm", "This field cannot be blank")
	form.CheckField(validator.ConfirmPassword(form.Password, form.Confirm), "confirm", "Passwords do not match")

	if !form.Valid(){
		return oops.ErrFormInvalid
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 12)
	if err != nil {
		return err
	}
	form.Password = string(hashedPassword)

	err = s.user.CreateUser(form)
	if err != nil {
		if errors.Is(err, oops.ErrDuplicateEmail) {
			form.AddFieldError("email", "User already exists")
		}
		return err
	}

	return nil
}

func (s *UserService) Authenticate(form *model.SignIn) (*model.User, error) {
	form.CheckField(validator.NotBlank(form.Phone), "email", "This field cannot be blank")
	form.CheckField(validator.CheckPhone(form.Phone), "phone", "This field must be a valid phone number")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid(){
		return nil, oops.ErrFormInvalid
	}

	user, err := s.user.Authenticate(form)
	if err != nil {
		if errors.Is(err, oops.ErrInvalidCredentials){
			form.AddNonFieldError("Invalid credentials")
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) Account(userID int) (*model.Account, error){
	return s.user.Account(userID)
}

func (s *UserService) Request(form *model.LoanForm) error{
	form.CheckField(validator.CheckLoan(form.Amount), "amount", "min amount - 5000, max amount - 50000")

	err := s.user.Request(form)
	if err != nil{
		if errors.Is(err, oops.ErrNotAllowed){
			form.AddNonFieldError("To request for loan, your account balance must be less than 5000")
		}
		return err
	}

	return nil
}