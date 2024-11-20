package main

import (
	"log"

	"github.com/ChekoutGobiz/BackendChekout/config"
	"github.com/ChekoutGobiz/BackendChekout/url"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Membuat instance aplikasi Fiber dengan konfigurasi yang telah disediakan
	app := fiber.New(config.GoBiz)

	// Menambahkan middleware CORS dengan pengaturan yang ada di config
	app.Use(config.Cors)

	// Setup semua routes
	url.SetupRoutes(app)

	// Memulai aplikasi pada port yang telah disetting
	log.Fatal(app.Listen(config.IPPort))
}
