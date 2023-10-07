package api

import (
	"bank/oops"
	"bank/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

var (
	conf = fiber.Config{
	// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {

			if e, ok := err.(oops.Error); ok {
				return ctx.Status(e.Code).JSON(map[string]string{"msg": e.Msg,})
			}
			oops.ErrorLog(err)
			return ctx.Status(http.StatusInternalServerError).JSON(map[string]string{
				"msg": http.StatusText(http.StatusInternalServerError),
			})
		},
	}
)

type Handler struct{
	infoLog *log.Logger
	errorLog *log.Logger
	service *service.Service
	app *fiber.App
}

func NewHandler(infoLog, errorLog *log.Logger, service *service.Service) *Handler{
	return &Handler{
		infoLog: infoLog,
		errorLog: errorLog,
		service: service,
	}
}

func (h *Handler) RunServer(port string) error {
	h.app = fiber.New(conf)

	h.app.Use(h.LogRequest)

	user := h.app.Group("/user")
	tx := h.app.Group("/operation", h.JWTauthentication)
	admin := h.app.Group("/admin", h.JWTauthentication, h.IsAdmin)

	user.Post("/signup", h.Register)
	user.Post("/signin", h.Authenticate)
	user.Get("/account", h.JWTauthentication, h.MyAccount)
	user.Post("/request", h.JWTauthentication, h.Request)

	tx.Post("/transact", h.Transaction)

	admin.Get("/accounts", h.Accounts)
	admin.Get("/account/:id", h.Account)
	admin.Get("/loans", h.Loans)
	admin.Post("/loan", h.LoanApproval)
	
	fmt.Printf("Working on localhost%s\n", port)
	return h.app.Listen(port)
}