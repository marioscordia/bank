package api

import (
	"bank/oops"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) LogRequest(c *fiber.Ctx) error {
	h.infoLog.Printf("%s - %s %s %s\n", c.IP(), c.Protocol(), c.Method(), c.OriginalURL())
	return c.Next()
}

func (h *Handler) JWTauthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		h.errorLog.Println("No X-Api-Token header")
		return oops.NewError(fiber.StatusForbidden, fiber.ErrForbidden.Message) 
	}

	claims, err := ParseToken(token) 
	if err != nil {
		oops.ErrorLog(err)
		return oops.NewError(fiber.StatusForbidden, fiber.ErrForbidden.Message) 
	}

	exp := claims["exp"].(float64)
	if float64(time.Now().Unix()) >= exp{
		h.errorLog.Println("Token is expired!")
		return oops.NewError(fiber.StatusForbidden, fiber.ErrForbidden.Message)
	}

	c.Context().SetUserValue("info", claims)
	return c.Next()
}

func (h *Handler) IsAdmin(c *fiber.Ctx) error {
	info := c.Context().UserValue("info").(jwt.MapClaims)
	
	IsAdmin := info["admin"].(bool)
	if !IsAdmin{
		h.errorLog.Println("Only admin is allowed")
		return oops.NewError(fiber.StatusForbidden, fiber.ErrForbidden.Message)
	}

	c.Context().SetUserValue("info", info)
	return c.Next()
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid{
		return nil, fmt.Errorf("token is invalid")
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok{
		// fmt.Println(claims)
		return claims, nil
	}

	return nil, fmt.Errorf("unauthenticated")
}