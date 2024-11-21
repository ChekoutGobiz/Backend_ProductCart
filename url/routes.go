package url

import (
	controllers "github.com/ChekoutGobiz/BackendChekout/controller"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes mendefinisikan semua rute aplikasi
func SetupRoutes(app *fiber.App) {
	// Grouping untuk lebih rapi
	api := app.Group("/api")

	// Authentication routes
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)

	// Region routes
	api.Post("/regions", controllers.CreateRegion)
	api.Get("/regions", controllers.GetRegions)

	// Product routes
	api.Post("/products", controllers.CreateProduct)
	api.Get("/products", controllers.GetProducts)

	// Cart routes
	api.Post("/cart", controllers.AddItemToCart)
	api.Get("/cart", controllers.GetCart)
	api.Put("/cart/item", controllers.UpdateCartItem)
	api.Delete("/cart/item", controllers.RemoveCartItem)
	api.Delete("/cart/item/:product_id", controllers.RemoveItemFromCart)
}
