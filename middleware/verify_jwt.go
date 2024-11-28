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
	log.Println("Received token:", tokenString)

	// Token harus diawali dengan "Bearer "
	if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
		log.Println("Invalid token format")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token format",
		})
	}

	// Ambil token tanpa kata "Bearer "
	tokenString = tokenString[7:]

	// Verifikasi token
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metode signing adalah HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Unexpected signing method")
			return nil, fiber.ErrUnauthorized
		}

		// Kembalikan secret key untuk validasi
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
