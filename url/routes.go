package url

import (
	controllers "github.com/ChekoutGobiz/BackendChekout/controller"
	"github.com/ChekoutGobiz/BackendChekout/middleware"
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

	// Product routes - Protected by JWT middleware
	api.Use(middleware.VerifyJWT) // Hanya proteksi rute produk
	api.Post("/products", controllers.CreateProduct)
	api.Get("/products", controllers.GetProducts)

	// Cart routes - Protected by JWT middleware
	cart := api.Group("/cart")
	cart.Use(middleware.VerifyJWT) // Hanya proteksi rute cart
	cart.Post("/", controllers.AddItemToCart)
	cart.Get("/", controllers.GetCart)
	cart.Put("/item", controllers.UpdateCartItem)
	cart.Delete("/item", controllers.RemoveCartItem)
	cart.Delete("/item/:product_id", controllers.RemoveItemFromCart)
}
