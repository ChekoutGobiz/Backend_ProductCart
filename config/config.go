package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Konfigurasi Fiber App
var GoBiz = fiber.Config{
	Prefork:       true,
	CaseSensitive: true,
	StrictRouting: true,
	ServerHeader:  "GoBiz",
	AppName:       "Gibizyuhu",
}

// Konfigurasi CORS (Fiber middleware cors.New)
var Cors = cors.New(cors.Config{
	AllowOrigins: "http://127.0.0.1:5501", // Sesuaikan dengan origin yang diizinkan
	AllowMethods: "GET,POST,PUT,DELETE",
	AllowHeaders: "Content-Type,Authorization",
})

// IPPort bisa disesuaikan dengan environment atau hardcoded
var IPPort = ":8080"
