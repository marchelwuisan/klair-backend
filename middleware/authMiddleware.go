package middleware

import (
	helper "golanglearn/helpers"
	"golanglearn/responses"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Auth validates token and authorizes users
func Authentication(c *fiber.Ctx) error {
	clientToken := c.Locals("token").(string)
	if clientToken == "" {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error getting token"})
	}

	claims, err := helper.ValidateToken(clientToken)
	if err != "" {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error generating claims"})

	}

	c.Set("email", claims.Email)
	c.Set("first_name", claims.FirstName)
	c.Set("last_name", claims.LastName)
	c.Set("uid", claims.Uid)

	c.Next()

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success"})
}
