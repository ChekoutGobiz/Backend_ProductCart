package config

import (
	"github.com/ChekoutGobiz/BackendChekout/helper"

	"github.com/gofiber/fiber/v2"
)

// Global variable to hold the database client instance
var IPPort, Net = helper.GetAddress()

var Iteung = fiber.Config{
	Prefork:       true,
	CaseSensitive: true,
	StrictRouting: true,
	ServerHeader:  "GoBiz",
	AppName:       "Gibizyuhu",
	Network:       Net,
}
