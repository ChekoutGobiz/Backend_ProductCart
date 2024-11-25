package middleware

import (
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// VerifyJWT memverifikasi token JWT yang diterima di header Authorization
func VerifyJWT(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	tokenString := c.Get("Authorization")

	// Token harus diawali dengan "Bearer "
	if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Ambil token tanpa kata "Bearer "
	tokenString = tokenString[7:]

	// Verifikasi token
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifikasi metode signature, pastikan menggunakan algoritma yang tepat
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Unexpected signing method")
			return nil, fiber.ErrUnauthorized
		}

		// Return JWT secret key dari env atau konfigurasi
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Println("Error verifying token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Lanjutkan ke handler berikutnya
	return c.Next()
}
