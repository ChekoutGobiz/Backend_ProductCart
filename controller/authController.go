package controllers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// Login user dan menghasilkan token JWT
func Login(c *fiber.Ctx) error {
	// Ambil data login dari request body
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid login data",
		})
	}

	// Untuk demo, kita hanya memeriksa username dan password statis
	if loginData.Username != "admin" || loginData.Password != "password" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Generate token JWT
	claims := jwt.MapClaims{
		"username": loginData.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating token",
		})
	}

	// Kirimkan token JWT sebagai response
	return c.JSON(fiber.Map{
		"token": tokenString,
	})
}
