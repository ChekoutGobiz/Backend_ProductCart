package main

import (
	"log"

	"gocroot/config"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"gocroot/url"

	"github.com/gofiber/fiber/v2"
)

func main() {
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	url.Web(site)
	log.Fatal(site.Listen(config.IPPort))
}
