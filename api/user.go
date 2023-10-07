package api

import (
	"bank/model"
	"bank/oops"
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) Register(c *fiber.Ctx) error{
	var signup model.SignUp

	if err := c.BodyParser(&signup); err != nil{
		return err
	}

	if err := h.service.Register(&signup); err != nil{
		if errors.Is(err, oops.ErrFormInvalid) || errors.Is(err, oops.ErrDuplicateEmail){
			return c.Status(fiber.StatusUnprocessableEntity).JSON(signup.FieldErrors)
		}
		return err
	}

	return c.SendString("Registered successfully!")
}

func (h *Handler) Authenticate(c *fiber.Ctx) error{
	var signin model.SignIn
	if err := c.BodyParser(&signin); err != nil{
		return err
	}

	user, err := h.service.Authenticate(&signin)
	if err != nil{
			if errors.Is(err, oops.ErrFormInvalid){
				return c.Status(fiber.StatusUnprocessableEntity).JSON(signin.FieldErrors)
			}else if errors.Is(err, oops.ErrInvalidCredentials){
				return c.Status(fiber.StatusUnprocessableEntity).JSON(signin.NonFieldErrors)
			}
			return err
	}

	token, err := CreateToken(user)
	if err != nil{
		return err
	}

	return c.SendString(token)
}

func CreateToken(u *model.User) (string, error) {
	claims := jwt.MapClaims{
		"id": u.ID,
		"admin": u.IsAdmin,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(key))
	if err != nil{
		return "", err
	}

	return tokenStr, nil
}

func (h *Handler) MyAccount(c *fiber.Ctx) error{
	info := c.Context().UserValue("info").(jwt.MapClaims)
	userID := int(info["id"].(float64))

	acc, err := h.service.Account(userID)
	if err != nil{
		return err
	}

	return c.JSON(acc)
}

func (h *Handler) Request(c *fiber.Ctx) error{
	info := c.Context().UserValue("info").(jwt.MapClaims)
	userID := int(info["id"].(float64))

	var loan *model.LoanForm

	if err := c.BodyParser(&loan); err != nil{
		return err
	}

	loan.UserID = userID
	if err := h.service.Request(loan); err != nil{
		if errors.Is(err, oops.ErrNotAllowed){
			return c.Status(fiber.StatusUnprocessableEntity).JSON(loan.NonFieldErrors)
		}
		return err
	}

	return c.SendString("Your request has been admitted. Wait for the approval!")
}
