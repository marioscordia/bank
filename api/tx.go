package api

import (
	"bank/model"
	"bank/oops"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) Transaction(c *fiber.Ctx) error{
	info := c.Context().UserValue("info").(jwt.MapClaims)
	userID := int(info["id"].(float64))

	var tx model.Tx
	if err := c.BodyParser(&tx); err != nil{
		return err
	}
	
	tx.SenderID = userID
	if err := h.service.Transact(c.Context(), &tx); err != nil{
		if errors.Is(err, oops.ErrFormInvalid) || errors.Is(err, oops.ErrNotEnough){
			return c.Status(fiber.StatusUnprocessableEntity).JSON(tx.FieldErrors)
		}
		
		return err
	}

	return c.SendString("Successfully processed!")
}