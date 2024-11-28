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
	AllowOrigins: "https://chekoutgobiz.github.io, http://localhost:5502", // Daftar origins yang dipisahkan koma
	AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",                           // Daftar metode yang dipisahkan koma
	AllowHeaders: "Content-Type,Authorization",                            // Daftar headers yang dipisahkan koma
})

// IPPort bisa disesuaikan dengan environment atau hardcoded
var IPPort = ":8080"
