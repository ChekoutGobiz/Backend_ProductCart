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
    AllowOrigins: "https://satsetin.github.io", // Izinkan origin frontend
    AllowMethods: "GET,POST,PUT,DELETE",        // Izinkan metode HTTP
    AllowHeaders: "Content-Type,Authorization",  // Izinkan header tertentu
    AllowCredentials: true,  
})

// IPPort bisa disesuaikan dengan environment atau hardcoded
var IPPort = ":8080"
