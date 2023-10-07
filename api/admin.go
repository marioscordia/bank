package api

import (
	"bank/model"
	"bank/oops"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)


func (h *Handler) Accounts(c *fiber.Ctx) error {
	accs, err := h.service.Accounts()
	if err != nil {
		return err
	}

	if len(accs) == 0 {
		return c.SendString("No users registered...yet!")
	}

	return c.JSON(accs)
}

func (h *Handler) Account(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return oops.NewError(fiber.StatusNotFound, fiber.ErrNotFound.Message)
	}

	acc, err := h.service.Account(id)
	if err != nil {
		if errors.Is(err, oops.ErrNoRecord){
			return oops.NewError(fiber.StatusNotFound, fiber.ErrNotFound.Message)
		}
		return err
	}

	return c.JSON(acc)
}

func (h *Handler) Loans(c *fiber.Ctx) error {
	loans, err := h.service.Loans()
	if err != nil {
		return err
	}

	if len(loans) == 0 {
		return c.SendString("No loan requests...yet!")
	}

	return c.JSON(loans)
}

func (h *Handler) LoanApproval(c *fiber.Ctx) error {
	var approval model.LoanDecision
	if err := c.BodyParser(&approval); err != nil {
		return err
	}

	if err := h.service.LoanApproval(c.Context(), &approval); err != nil {
		return err
	}

	return c.SendString("Done!")
}


